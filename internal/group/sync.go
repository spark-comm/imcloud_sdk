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

package group

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"

	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/common/log"
)

// SyncAllGroupMember 同步所有群成员
func (g *Group) SyncAllGroupMember(ctx context.Context, groupID string) error {
	ctx1, cl := context.WithTimeout(ctx, time.Minute*5)
	defer cl()
	//获取远程的成员列表
	members, err := g.GetServerGroupMembers(ctx1, groupID)
	if err != nil {
		return err
	}
	//获取本地的成员列表
	localData, err := g.db.GetGroupMemberListSplit(ctx1, groupID, 0, 0, 9999999)
	if err != nil {
		return err
	}
	log.ZInfo(ctx1, "SyncGroupMember Info", "groupID", groupID, "members", len(members), "localData", len(localData))

	//util.Batch(ServerGroupMemberToLocalGroupMember, members) 远程数据序列化为本地结构
	err = g.groupMemberSyncer.Sync(ctx1, util.Batch(ServerGroupMemberToLocalGroupMember, members), localData, nil)
	if err != nil {
		return err
	}
	//if len(members) != len(localData) {
	log.ZInfo(ctx1, "SyncGroupMember Sync Group Member Count", "groupID", groupID, "members", len(members), "localData", len(localData))
	//获取远程群组数据（单条）
	gs, err := g.GetSpecifiedGroupsInfo(ctx1, []string{groupID})
	if err != nil {
		return err
	}
	log.ZInfo(ctx1, "SyncGroupMember GetGroupsInfo", "groupID", groupID, "len", len(gs), "gs", gs)
	if len(gs) > 0 {
		v := gs[0]
		count := int32(len(members))
		if v.MemberCount != count {
			v.MemberCount = int32(len(members))
			if v.GroupType == constant.SuperGroupChatType {
				if err := g.db.UpdateSuperGroup(ctx1, v); err != nil {
					//return err
					log.ZError(ctx1, "SyncGroupMember UpdateSuperGroup", err, "groupID", groupID, "info", v)
				}
			} else {
				if err := g.db.UpdateGroup(ctx1, v); err != nil {
					log.ZError(ctx1, "SyncGroupMember UpdateGroup", err, "groupID", groupID, "info", v)
				}
			}
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			log.ZInfo(ctx1, "SyncGroupMember OnGroupInfoChanged", "groupID", groupID, "data", string(data))
			g.listener.OnGroupInfoChanged(string(data))
		}
	}
	//标记群已经同步
	g.syncGroup[groupID] = true
	return nil
}

// SyncAllJoinedGroups 同步所有加入的群
func (g *Group) SyncAllJoinedGroups(ctx context.Context) error {
	_, err := g.syncJoinedGroup(ctx)
	if err != nil {
		return err
	}
	return err
}

// syncJoinedGroup 同步所有的群
func (g *Group) syncJoinedGroup(ctx context.Context) ([]*groupv1.GroupInfo, error) {
	//获取登录用户加入的群组列表（远程数据）
	groups, err := g.GetServerJoinGroup(ctx)
	if err != nil {
		return nil, err
	}
	//本地所有群组数据
	localData, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	//获取用户自己的加群申请信息（远程数据）
	list, err := g.GetServerSelfGroupApplication(ctx)
	if err != nil {
		return err
	}
	//获取本地加群请求列表
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	//更新/删除操作
	if err := g.groupRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil); err != nil {
		return err
	}
	// todo
	return nil
}

func (g *Group) SyncAdminGroupApplication(ctx context.Context) error {
	//(以管理员或群主身份)获取群的加群申请（远程数据）
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	//本地加群申请数据
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

// SyncAdminGroupUntreatedApplication 获取未处理的加群请求
func (g *Group) SyncAdminGroupUntreatedApplication(ctx context.Context) error {
	//(以管理员或群主身份)获取群的加群申请（远程数据）
	requests, err := g.GetServerAdminGroupUntreatedApplicationList(ctx)
	if err != nil {
		return err
	}
	//本地加群申请数据
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

// SyncGroups 根据id同步群信息
func (g *Group) SyncGroups(ctx context.Context, id ...string) error {
	//根据id获取群（远程数据）
	groups, err := g.getGroupsInfoFromSvr(ctx, id)
	if err != nil {
		return err
	}
	log.ZError(ctx, fmt.Sprintf("get local group data %s", utils.StructToJsonString(groups)), err)
	//根据id获取本地群数据
	localData, err := g.db.GetGroupInfoByGroupIDs(ctx, id...)
	log.ZError(ctx, fmt.Sprintf("get local group data %s", utils.StructToJsonString(localData)), err)
	if err != nil {
		localData = make([]*model_struct.LocalGroup, 0)
		log.ZError(ctx, "SyncGroups->GetGroupInfoByGroupIDs", err)
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return err
	}
	return nil
}

// deleteGroup 删除群
func (g *Group) deleteGroup(ctx context.Context, groupID string) error {
	groupInfo, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if err := g.db.DeleteGroup(ctx, groupID); err != nil {
		return err
	}
	// 删除群成员
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		log.ZDebug(ctx, "delete  all group member err ", err)
	}
	//删除会话
	g.DelGroupConversation(ctx, groupInfo.GroupID)
	// 删除群通知
	g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(groupInfo))
	return nil
}

// SyncAllJoinedGroupsAndMembers 同步所有的群和群成员
func (g *Group) SyncAllJoinedGroupsAndMembers(ctx context.Context) error {
	err := g.SyncAllJoinedGroups(ctx)
	if err != nil {
		return err
	}
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			if err := g.SyncAllGroupMember(ctx, groupID); err != nil {
				log.ZError(ctx, "SyncGroupMember failed", err)
			}
		}(group.GroupID)
	}
	wg.Wait()
	return nil
}

// deleteGroupMembers 删除指定群成员
func (g *Group) deleteGroupMembers(ctx context.Context, groupID string, memberId ...string) error {
	//if err := g.db.DeleteGroupMembers(ctx, groupID, memberId...); err != nil {
	//	return err
	//}
	//同步群数据
	if err := g.SyncGroups(ctx, groupID); err != nil {
		log.ZDebug(ctx, "syncGroupAndMembers->SyncGroups err", err)
	}
	//同步群成员数据
	if err := g.syncGroupMembers(ctx, groupID, memberId...); err != nil {
		log.ZDebug(ctx, "syncGroupAndMembers->syncGroupMembers err", err)
		return err
	}
	return nil
}

// syncGroupAndMember 同步群和群成员
func (g *Group) syncGroupAndMembers(ctx context.Context, groupId string, memberId ...string) error {
	//同步群数据
	if err := g.SyncGroups(ctx, groupId); err != nil {
		log.ZDebug(ctx, "syncGroupAndMembers->SyncGroups err", err)
	}
	//同步群成员数据
	if err := g.syncGroupMembers(ctx, groupId, memberId...); err != nil {
		log.ZDebug(ctx, "syncGroupAndMembers->syncGroupMembers err", err)
		return err
	}
	return nil
}

// syncGroupMembers 同步指定数据
func (g *Group) syncGroupMembers(ctx context.Context, groupID string, userIDs ...string) error {
	members, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs...)
	if err != nil {
		return err
	}
	localData, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDs)
	if err != nil {
		return err
	}
	if err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, members), localData, nil); err != nil {
		return err
	}
	return nil
}

// syncGroupAndMember 同步群和成员
func (g *Group) syncGroupAndMember(ctx context.Context, groupId string) error {
	if err := g.SyncGroups(ctx, groupId); err != nil {
		log.ZDebug(ctx, "syncGroupAndMember-> SyncGroups err", err)
	}
	if err := g.SyncAllGroupMember(ctx, groupId); err != nil {
		log.ZDebug(ctx, "syncGroupAndMember->SyncAllGroupMember err", err)
		return err
	}
	return nil
}

// SyncAllSelfGroupApplication 同步自己的加群申请
func (g *Group) SyncAllSelfGroupApplication(ctx context.Context) error {
	list, err := g.GetServerSelfGroupApplication(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetSendGroupApplication(ctx) //待处理
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil); err != nil {
		return err
	}
	// todo
	return nil
}

// SyncSelfGroupApplications  同步自己某个群的加群申请
func (g *Group) SyncSelfGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllSelfGroupApplication(ctx)
}

func (g *Group) SyncOneGroupApplications(ctx context.Context, groupID string) error {
	return g.SyncOneGroupApplicationsFunc(ctx, groupID)
}

// SyncAllAdminGroupApplication 同步所有群的加群申请
func (g *Group) SyncAllAdminGroupApplication(ctx context.Context) error {
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

// SyncAdminGroupApplications 根据群同步加群申请
func (g *Group) SyncAdminGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllAdminGroupApplication(ctx)
}

// InitSyncData 初始化同步数据
func (g *Group) InitSyncData(ctx context.Context) error {
	//从群中获取数据
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return err
	}
	serveGroups, err := g.GetServerFirstPageJoinGroup(ctx)
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, serveGroups), groups, nil); err != nil {
		return err
	}
	//延迟30秒同步全量数据
	g.syncGroupQueue.Push(1, time.Second*30)
	// 发出同步的通知
	common.TriggerCmdJoinGroup(ctx, g.groupCh)
	return nil
}

// delaySyncJoinGroup 延迟同步加入的群
func (g *Group) delaySyncJoinGroup(ctx context.Context) {
	ctx1, cl := context.WithTimeout(ctx, time.Minute*5)
	defer cl()
	//从延迟队列中取数据
	for emtry := range g.syncGroupQueue.Channel(ctx1, 1) {
		log.ZDebug(ctx1, "delay sync join group", emtry)
		err := g.SyncAllJoinedGroups(ctx)
		if err != nil {
			log.ZError(ctx, "SyncAllJoinedGroups failed", err)
		}
		// 同步所有群成员
		go g.SyncAllJoinedGroupMembers(ctx1)
	}
}

// SyncAllJoinedGroupMembers 同步所有加入群的群成员
func (g *Group) SyncAllJoinedGroupMembers(ctx context.Context) error {
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, group := range groups {
		// 获取用户在群中的信息
		groupMember, err := g.db.GetGroupMemberInfoByGroupIDUserID(ctx, group.GroupID, g.loginUserID)
		if err != nil {
			log.ZError(ctx, "sync all group member failed get group member ", err)
			continue
		}
		//普通成员不同步成员的所有信息
		if groupMember == nil || groupMember.RoleLevel == constant.GroupOrdinaryUsers {
			continue
		}
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			// 同步群成员
			if err := g.SyncGroupMembers(ctx, groupID); err != nil {
				log.ZError(ctx, "SyncGroupMember failed", err)
			}
			//if err := g.SyncAllGroupMember(ctx, groupID); err != nil {
			//	log.ZError(ctx, "SyncGroupMember failed", err)
			//}
		}(group.GroupID)
	}
	wg.Wait()
	return nil
}

// SyncQuantityGroupList 全量同步群
func (g *Group) SyncQuantityGroupList(ctx context.Context) error {
	//从群中获取数据
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		log.ZError(ctx, "InitSyncGroupData err", err)
	}
	serveGroups, err := g.GetSyncGroup(ctx)
	if err != nil {
		return err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerBaseGroupToLocalGroup, serveGroups), groups, nil); err != nil {
		return err
	}
	//延迟30秒同步全量数据
	//g.syncGroupQueue.Push(1, time.Second*30)
	// 发出同步的通知
	//common.TriggerCmdJoinGroup(ctx, g.groupCh)
	return nil
}

// SyncOneGroupApplicationsFunc 同步单个群的加群请求
func (g *Group) SyncOneGroupApplicationsFunc(ctx context.Context, groupID string) error {
	//远程单个群的申请
	svrGroupReq, err := g.GetServerAdminOneGroupUntreatedApplicationList(ctx, groupID)
	if err != nil {
		return err
	}
	//本地数据
	localData, err := g.db.GetOneSendGroupApplication(ctx, groupID)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalGroupRequest, svrGroupReq), localData, nil); err != nil {
		return err
	}
	return nil
}

// SyncUserMemberInfoInGroup 同步用户在所有群中的信息
func (g *Group) SyncUserMemberInfoInGroup(ctx context.Context) error {
	//初始化同步用户的所有群中成员信息
	info, _ := g.db.GetOwnerGroupMemberInfo(ctx, g.loginUserID)
	groups, err := g.GetUserMemberInfoInGroup(ctx)
	if err != nil {
		log.ZError(ctx, "sync user member info in group error", err)
		return nil
	}
	//同步数据
	if err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, groups), info, nil); err != nil {
		return err
	}
	return nil
}

// SyncGroup 同步群信息
func (g *Group) SyncGroup(ctx context.Context) error {
	udata, err := g.db.GetGroupUpdateTime(ctx)
	if err != nil || len(udata) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(udata) == 0 {
			return g.SyncQuantityGroupList(ctx)
		}
		return err
	} else {
		return g.SyncGroupByTime(ctx, udata)
	}
}

// SyncGroupMembers 同步群成员信息
func (g *Group) SyncGroupMembers(ctx context.Context, groupId string) error {
	udata, err := g.db.GetGroupMemberUpdateTime(ctx, groupId)
	if err != nil || len(udata) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(udata) == 0 {
			return g.SyncOneGroupMember(ctx, groupId)
		}
		return err
	} else {
		return g.SyncGroupMemberByTime(ctx, groupId, udata)
	}
}

// SyncOneGroupMember 同步单个群的成员
func (g *Group) SyncOneGroupMember(ctx context.Context, groupID string) error {
	//远程单个群的申请
	svrGroupReq, err := g.SyncGroupMemberInfo(ctx, groupID)
	if err != nil {
		return err
	}
	//本地数据
	localData, err := g.db.GetGroupMemberListByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerBaseGroupMemberToLocalGroupMember, svrGroupReq), localData, nil); err != nil {
		return err
	}
	return nil
}

// SyncGroupByTime 同步群组信息
func (g *Group) SyncGroupByTime(ctx context.Context, udata map[string]int64) error {
	addList, updateList, delIds, err := g.SyncUserJoinGroupInfoByTime(ctx, udata)
	if err != nil {
		log.ZError(ctx, "sync group error", err)
		return err
	}
	//处理新增
	if len(addList) > 0 {
		for _, v := range addList {
			localGroup := ServerBaseGroupToLocalGroup(v)
			err = g.db.InsertGroup(ctx, localGroup)
			if err != nil {
				log.ZError(ctx, "insert friend error", err)
			}
			g.listener.OnGroupMemberAdded(utils.StructToJsonString(localGroup))
		}
	}
	//处理修复爱
	if len(updateList) > 0 {
		for _, v := range updateList {
			localGroup := ServerBaseGroupToLocalGroup(v)
			err = g.db.UpdateGroup(ctx, localGroup)
			if err != nil {
				log.ZError(ctx, "insert friend error", err)
			}
			g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(localGroup))
		}
	}
	//处理删除
	if len(delIds) > 0 {
		for _, v := range delIds {
			err := g.deleteGroup(ctx, v)
			if err != nil {
				log.ZError(ctx, "insert friend error", err)
			}
		}
	}
	return nil
}

// SyncGroupMemberByTime 同步群组成员信息
func (g *Group) SyncGroupMemberByTime(ctx context.Context, groupId string, udata map[string]int64) error {
	addList, updateList, delIds, err := g.SyncGroupMemberInfoUpdateTime(ctx, groupId, udata)
	if err != nil {
		log.ZError(ctx, "sync group member error", err)
		return err
	}
	//处理新增
	if len(addList) > 0 {
		for _, v := range addList {
			localGroupMember := ServerBaseGroupMemberToLocalGroupMember(v)
			err = g.db.InsertGroupMember(ctx, localGroupMember)
			if err != nil {
				log.ZError(ctx, "insert group member error", err)
			}
		}
	}
	//处理修复爱
	if len(updateList) > 0 {
		for _, v := range updateList {
			localGroupMember := ServerBaseGroupMemberToLocalGroupMember(v)
			err = g.db.UpdateGroupMember(ctx, localGroupMember)
			if err != nil {
				log.ZError(ctx, "update group member error", err)
			}
		}
	}
	//处理删除
	if len(delIds) > 0 {
		err := g.deleteGroupMembers(ctx, groupId, delIds...)
		if err != nil {
			log.ZError(ctx, "delete group member error", err)
		}
	}
	return nil
}
