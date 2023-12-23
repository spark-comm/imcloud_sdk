// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import (
	"context"
	"open_im_sdk/internal/cache"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"strings"
	"time"

	"github.com/imCloud/im/pkg/common/log"
)

// SyncConversations 同步会话
func (c *Conversation) SyncConversations(ctx context.Context) error {
	ccTime := time.Now()
	conversationsOnServer, err := c.getServerConversationList(ctx)
	if err != nil {
		log.ZError(ctx, "get server conversation list failed", err)
		return err
	}
	if len(conversationsOnServer) == 0 {
		return nil
	}
	log.ZDebug(ctx, "get server cost time", "cost time", time.Since(ccTime), "conversation on server", conversationsOnServer)
	conversationsOnLocal, err := c.db.GetAllConversations(ctx)
	if err != nil {
		log.ZError(ctx, "get server conversation list failed", err)
	}
	log.ZDebug(ctx, "get local cost time", "cost time", time.Since(ccTime), "conversation on local", conversationsOnLocal)
	conversationsOnServer, err = c.CompletingConversationInformation(ctx, conversationsOnServer)
	if err != nil {
		log.ZError(ctx, "get server conversation list failed", err)
		return err
	}
	if err = c.conversationSyncer.Sync(ctx, conversationsOnServer, conversationsOnLocal, func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
		if state == syncer.Update || state == syncer.Insert {
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: server.ConversationID, Action: constant.ConChange, Args: []string{server.ConversationID}}})
		}
		return nil
	}, true); err != nil {
		return err
	}
	conversationsOnLocal, err = c.db.GetAllConversations(ctx)
	if err != nil {
		return err
	}
	c.cache.UpdateConversations(conversationsOnLocal)
	return nil
}

func (c *Conversation) SyncConversationUnreadCount(ctx context.Context) error {
	var conversationChangedList []string
	allConversations := c.cache.GetAllHasUnreadMessageConversations()
	log.ZDebug(ctx, "get unread message length", "len", len(allConversations))
	for _, conversation := range allConversations {
		if deleteRows := c.db.DeleteConversationUnreadMessageList(ctx, conversation.ConversationID, conversation.UpdateUnreadCountTime); deleteRows > 0 {
			log.ZDebug(ctx, "DeleteConversationUnreadMessageList", conversation.ConversationID, conversation.UpdateUnreadCountTime, "delete rows:", deleteRows)
			if err := c.db.DecrConversationUnreadCount(ctx, conversation.ConversationID, deleteRows); err != nil {
				log.ZDebug(ctx, "DecrConversationUnreadCount", conversation.ConversationID, conversation.UpdateUnreadCountTime, "decr unread count err:", err.Error())
			} else {
				conversationChangedList = append(conversationChangedList, conversation.ConversationID)
			}
		}
	}
	if len(conversationChangedList) > 0 {
		if err := common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: conversationChangedList}, c.GetCh()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conversation) SyncConversationHashReadSeqs(ctx context.Context) error {
	log.ZDebug(ctx, "start SyncConversationHashReadSeqs")
	seqs, err := c.getServerHasReadAndMaxSeqs(ctx)
	if err != nil {
		return err
	}
	if len(seqs) == 0 {
		return nil
	}
	var conversations []*model_struct.LocalConversation
	var conversationIDs []string
	allConversations, err := c.db.GetAllConversationIDList(ctx)
	for conversationID, v := range seqs {
		c.maxSeqRecorder.Set(conversationID, v.MaxSeq)
		if len(allConversations) == 0 {
			continue
		}
		var unreadCount int32
		if v.MaxSeq-v.HasReadSeq < 0 {
			unreadCount = 0
		} else {
			unreadCount = int32(v.MaxSeq - v.HasReadSeq)
		}
		// 初次登录更新全会报错
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": unreadCount, "has_read_seq": v.HasReadSeq}); err != nil {
			log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
		}
		conversationIDs = append(conversationIDs, conversationID)
	}
	log.ZDebug(ctx, "update conversations", "conversations", conversations)
	// 会话注册
	if len(conversations) > 0 {
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: conversationIDs}, c.GetCh())
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.TotalUnreadMessageChanged, Args: conversationIDs}, c.GetCh())
	}
	return nil
}

// CompletingConversationInformation 补齐会话细信息
func (c *Conversation) CompletingConversationInformation(ctx context.Context, data []*model_struct.LocalConversation) ([]*model_struct.LocalConversation, error) {
	singleConversationRevIDs := make([]string, 0)
	groupConversationGroupIDs := make([]string, 0)
	mapBaseInfo := make(map[string]*cache.BaseInfo)
	localFriendIds := make([]string, 0)
	localGroupIds := make([]string, 0)
	//获取所有的好友
	list, err := c.db.GetAllFriendList(ctx)
	if err != nil {
		log.ZError(ctx, "CompletingConversationInformation get user friend err", err)
	}
	if len(list) > 0 {
		for _, v := range list {
			localFriendIds = append(localFriendIds, v.FriendUserID)
			showName := v.Nickname
			if v.Remark != "" {
				showName = v.Remark
			}
			mapBaseInfo[v.FriendUserID] = &cache.BaseInfo{
				FaceURL:       v.FaceURL,
				Nickname:      v.Nickname,
				BackgroundURL: showName,
			}
		}
	}
	//获取所有的群
	groups, err := c.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		log.ZError(ctx, "CompletingConversationInformation get user join group err", err)
	}
	if len(groups) > 0 {
		//获取用户在群中的信息
		memberList, err := c.db.GetUserInAllGroupMemberList(ctx, c.loginUserID)
		if err != nil {
			log.ZError(ctx, "CompletingConversationInformation get user in group member err", err)
		}
		mapGroupMember := make(map[string]model_struct.LocalGroupMember)
		if len(memberList) > 0 {
			for _, v := range memberList {
				mapGroupMember[v.GroupID] = v
			}
		}
		for _, v := range groups {
			localGroupIds = append(localGroupIds, v.GroupID)
			baseInfo := &cache.BaseInfo{
				FaceURL:  v.FaceURL,
				Nickname: v.GroupName,
			}
			if m, ok := mapGroupMember[v.GroupID]; ok {
				baseInfo.BackgroundURL = m.BackgroundURL
			}
		}
	}
	//取相关id
	for _, v := range data {
		if c.IsGroupConversation(v) {
			groupConversationGroupIDs = append(groupConversationGroupIDs, v.GroupID)
		} else {
			split := strings.Split(v.ConversationID, "_")
			if len(split) > 2 {
				singleConversationRevIDs = append(singleConversationRevIDs, split[2])
			} else if len(split) == 2 {
				singleConversationRevIDs = append(singleConversationRevIDs, split[1])
			}
		}
	}
	//取不在好友列表的单聊会话id
	if len(singleConversationRevIDs) > len(localFriendIds) {
		notInFriendIds := utils.DifferenceSubsetString(singleConversationRevIDs, localFriendIds)
		profile, err := c.user.FindFullProfile(ctx, notInFriendIds...)
		if err != nil {
			log.ZError(ctx, "CompletingConversationInformation get user in group member err", err)
		}
		if profile != nil && len(profile) > 0 {
			for _, v := range profile {
				mapBaseInfo[v.UserId] = &cache.BaseInfo{
					FaceURL:  v.FaceURL,
					Nickname: v.Nickname,
				}
			}
		}
	}
	//获取不在本地的群信息
	if len(groupConversationGroupIDs) > len(localGroupIds) {
		notInGroupIds := utils.DifferenceSubsetString(groupConversationGroupIDs, localGroupIds)
		groupInfos, err := c.group.FindFullGroupInfo(ctx, notInGroupIds...)
		if err != nil {
			log.ZError(ctx, "CompletingConversationInformation get user in group member err", err)
		}
		if groupInfos != nil && len(groupInfos) > 0 {
			for _, v := range groupInfos {
				mapBaseInfo[v.GroupID] = &cache.BaseInfo{
					FaceURL:  v.FaceURL,
					Nickname: v.NickName,
				}
			}
		}
	}
	res := make([]*model_struct.LocalConversation, len(data))
	for i, v := range data {
		key := v.UserID
		if c.IsGroupConversation(v) {
			key = v.GroupID
		}
		if info, ok := mapBaseInfo[key]; ok {
			v.ShowName = info.Nickname
			v.FaceURL = info.FaceURL
			v.BackgroundURL = info.BackgroundURL
			//更新缓存
			c.cache.Store(key, info)
		}
		res[i] = v
	}
	return res, nil
}
