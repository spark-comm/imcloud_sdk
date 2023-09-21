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
	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/group"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/delayqueue"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
)

func NewGroup(loginUserID string, db db_interface.DataBase,
	conversationCh, groupCh chan common.Cmd2Value) *Group {
	g := &Group{
		loginUserID:    loginUserID,
		db:             db,
		groupCh:        groupCh,
		syncGroup:      make(map[string]bool),
		conversationCh: conversationCh,
		syncGroupQueue: delayqueue.New[int](),
	}
	g.initSyncer()
	return g
}

// //utils.GetCurrentTimestampByMill()
type Group struct {
	listener                open_im_sdk_callback.OnGroupListener
	loginUserID             string
	db                      db_interface.DataBase
	groupSyncer             *syncer.Syncer[*model_struct.LocalGroup, string]
	groupMemberSyncer       *syncer.Syncer[*model_struct.LocalGroupMember, [2]string]
	groupRequestSyncer      *syncer.Syncer[*model_struct.LocalGroupRequest, [2]string]
	groupAdminRequestSyncer *syncer.Syncer[*model_struct.LocalAdminGroupRequest, [2]string]
	loginTime               int64
	joinedSuperGroupCh      chan common.Cmd2Value
	heartbeatCmdCh          chan common.Cmd2Value
	groupCh                 chan common.Cmd2Value
	conversationCh          chan common.Cmd2Value
	syncGroup               map[string]bool
	// 同步群组信息延迟队列
	syncGroupQueue     *delayqueue.DelayQueue[int]
	listenerForService open_im_sdk_callback.OnListenerForService
}

// Work 群工作
func (g *Group) Work(c2v common.Cmd2Value) {
	switch c2v.Cmd {
	case constant.CmdGroupMemberChange:
		g.handelGroupMemberInfo(c2v)
	case constant.CmdSyncGroup:
		//延迟同步群信息
		g.delaySyncJoinGroup(c2v.Ctx)
	case constant.CmdSyncGroupMembers:
		//同步群成员
		g.SyncAllGroupMember(c2v.Ctx, c2v.Value.(string))
	}
}

// GetCh 获取通道
func (g *Group) GetCh() chan common.Cmd2Value {
	return g.groupCh
}
func (g *Group) initSyncer() {
	g.groupSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroup) error {
		return g.db.InsertGroup(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroup) error {
		if err := g.db.DeleteGroupAllMembers(ctx, value.GroupID); err != nil {
			return err
		}
		return g.db.DeleteGroup(ctx, value.GroupID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroup) error {
		log.ZInfo(ctx, "groupSyncer trigger update funcation", "groupID", server.GroupID, "server", server, "local", local)
		return g.db.UpdateGroup(ctx, server)
	}, func(value *model_struct.LocalGroup) string {
		return value.GroupID
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroup) error {
		if g.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			g.listener.OnJoinedGroupAdded(utils.StructToJsonString(server))
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: server.GroupID,
				SessionType: constant.SuperGroupChatType, FaceURL: server.FaceURL, Nickname: server.GroupName}}, g.conversationCh)
		case syncer.Delete:
			g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(local))
		case syncer.Update:
			log.ZInfo(ctx, "groupSyncer trigger update", "groupID", server.GroupID, "data", server, "isDismissed", server.Status == constant.GroupStatusDismissed)
			if server.Status == constant.GroupStatusDismissed {
				if err := g.db.DeleteGroupAllMembers(ctx, server.GroupID); err != nil {
					log.ZError(ctx, "delete group all members failed", err)
				}
				g.listener.OnGroupDismissed(utils.StructToJsonString(server))
			} else {
				g.listener.OnGroupInfoChanged(utils.StructToJsonString(server))
				if server.GroupName != local.GroupName || local.FaceURL != server.FaceURL {
					_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: server.GroupID,
						SessionType: constant.SuperGroupChatType, FaceURL: server.FaceURL, Nickname: server.GroupName}}, g.conversationCh)
				}
			}
		}

		return nil
	})

	g.groupMemberSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.InsertGroupMember(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.DeleteGroupMember(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupMember) error {
		return g.db.UpdateGroupMember(ctx, server)
	}, func(value *model_struct.LocalGroupMember) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroupMember) error {
		if g.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			//log.ZError(ctx, fmt.Sprintf("触发自己信息同步%s", utils.StructToJsonString(server)), nil)
			//if server.UserID == g.loginUserID {
			//	//如果是自己则更新会话的背景
			//	sessionType := g.getConversationIDBySessionType(server.GroupID, constant.WorkingGroup)
			//	_ = common.TriggerCmdUpdateConversationBackgroundURL(ctx, sessionType, server.BackgroundURL, g.conversationCh)
			//}
			g.listener.OnGroupMemberAdded(utils.StructToJsonString(server))
		case syncer.Delete:
			g.listener.OnGroupMemberDeleted(utils.StructToJsonString(local))
		case syncer.Update:
			g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(server))
			if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL || server.GroupUserName != local.GroupUserName {
				nickname := server.Nickname
				if server.GroupUserName != "" {
					nickname = server.GroupUserName
				}
				// 更新本地消息
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: server.UserID, FaceURL: server.FaceURL,
					Nickname: nickname, GroupID: server.GroupID}}, g.conversationCh)
			}
		}
		return nil
	})

	g.groupRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.InsertGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.DeleteGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupRequest) error {
		return g.db.UpdateGroupRequest(ctx, server)
	}, func(value *model_struct.LocalGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

	g.groupAdminRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.InsertAdminGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.DeleteAdminGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalAdminGroupRequest) error {
		return g.db.UpdateAdminGroupRequest(ctx, server)
	}, func(value *model_struct.LocalAdminGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalAdminGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

}

func (g *Group) SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		return
	}
	g.listener = callback
}

func (g *Group) LoginTime() int64 {
	return g.loginTime
}

func (g *Group) SetLoginTime(loginTime int64) {
	g.loginTime = loginTime
}

func (g *Group) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	g.listenerForService = listener
}

func (g *Group) GetGroupOwnerIDAndAdminIDList(ctx context.Context, groupID string) (ownerID string, adminIDList []string, err error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	adminIDList, err = g.db.GetGroupAdminID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	return localGroup.OwnerUserID, adminIDList, nil
}

// GetGroupInfoFromLocal2Svr 从服务端获取群信息
func (g *Group) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err == nil && localGroup.GroupID != "" {
		return localGroup, nil
	}
	svrGroup, err := g.getGroupsInfoFromSvr(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}
	if len(svrGroup) == 0 {
		return nil, sdkerrs.ErrGroupIDNotFound.Wrap("server not this group")
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, svrGroup), []*model_struct.LocalGroup{localGroup}, nil); err != nil {
		log.ZDebug(ctx, "sync group info err:%v", err)
	}
	return ServerGroupToLocalGroup(svrGroup[0]), nil
}

// GetGroupInfoAndSelfGroupMemberInfoFromLocal2Svr 从本地和服务端获取群信息和当前用户在群中的信息
func (g *Group) GetGroupInfoAndSelfGroupMemberInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, *model_struct.LocalGroupMember, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	var localHaveGroup bool
	if err == nil && localGroup.GroupID != "" {
		localHaveGroup = true
	}
	localSelfGroupMember, err := g.db.GetGroupMemberInfoByGroupIDUserID(ctx, groupID, g.loginUserID)
	if err != nil {
		log.ZError(ctx, "GetGroupInfoAndSelfGroupMemberInfoFromLocal2Svr->GetGroupMemberInfoByGroupIDUserID failed", err)
	}
	if localHaveGroup {
		return localGroup, localSelfGroupMember, nil
	}
	return localGroup, localSelfGroupMember, nil
	//if localHaveGroup {
	//	return localGroup, localSelfGroupMember, nil
	//}
	//svrGroup, err := g.getGroupsInfoFromSvr(ctx, []string{groupID})
	//if err != nil {
	//	return nil, localSelfGroupMember, err
	//}
	//if len(svrGroup) == 0 {
	//	return nil, localSelfGroupMember, sdkerrs.ErrGroupIDNotFound.Wrap("server not this group")
	//}
	//if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, svrGroup), []*model_struct.LocalGroup{localGroup}, nil); err != nil {
	//	log.ZDebug(ctx, "sync group info err:%v", err)
	//}
	//return ServerGroupToLocalGroup(svrGroup[0]), localSelfGroupMember, nil
}

// getGroupsInfoFromSvr 从服务端获取群数据
func (g *Group) getGroupsInfoFromSvr(ctx context.Context, groupIDs []string) ([]*groupv1.GroupInfo, error) {
	//resp, err := util.CallApi[groupv1.GetGroupInfoResponse](ctx, constant.GetGroupsInfoRouter, &groupv1.GetGroupInfoReq{GroupID: groupIDs})
	resp := &groupv1.GetGroupInfoResponse{}
	err := util.CallPostApi[*groupv1.GetGroupInfoReq, *groupv1.GetGroupInfoResponse](
		ctx, constant.GetGroupsInfoRouter, &groupv1.GetGroupInfoReq{GroupID: groupIDs}, resp,
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (g *Group) getGroupAbstractInfoFromSvr(ctx context.Context, groupIDs []string) (*group.GetGroupAbstractInfoResp, error) {
	//return util.CallApi[group.GetGroupAbstractInfoResp](ctx, constant.GetGroupAbstractInfoRouter,
	//	&group.GetGroupAbstractInfoReq{GroupIDs: groupIDs})
	resp := &group.GetGroupAbstractInfoResp{}
	err := util.CallPostApi[*group.GetGroupAbstractInfoReq, *group.GetGroupAbstractInfoResp](
		ctx, constant.GetGroupAbstractInfoRouter,
		&group.GetGroupAbstractInfoReq{GroupIDs: groupIDs},
		resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *Group) GetJoinedDiffusionGroupIDListFromSvr(ctx context.Context) ([]string, error) {
	groups, err := g.GetServerJoinGroup(ctx)
	if err != nil {
		return nil, err
	}
	var groupIDs []string
	for _, g := range groups {
		if g.GroupType == constant.WorkingGroup {
			groupIDs = append(groupIDs, g.GroupID)
		}
	}
	return groupIDs, nil
}

// GetGroupMemberFromLocal2Svr 从本地和服务端获取群成员数据
func (g *Group) GetGroupMemberFromLocal2Svr(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	localGroupMembers, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
	if err == nil && len(localGroupMembers) > 0 {
		return localGroupMembers, nil
	}
	svrGroupMembers, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDList...)
	if err != nil {
		return nil, err
	}
	if len(svrGroupMembers) == 0 {
		return nil, sdkerrs.ErrGroupMemberNotFound.Wrap("group member not found")
	}
	if err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, svrGroupMembers), localGroupMembers, nil); err != nil {
		log.ZDebug(ctx, "sync group member info err:%v", err)
	}
	return util.Batch(ServerGroupMemberToLocalGroupMember, svrGroupMembers), nil
}

// handelGroupMemberInfo 处理群成员信息变更
func (g *Group) handelGroupMemberInfo(c2v common.Cmd2Value) {
	info := c2v.Value.(common.UpdateGroupMemberInfo)
	agrs := make(map[string]interface{})
	agrs["face_url"] = info.FaceUrl
	agrs["nickname"] = info.Nickname
	ctx := context.Background()
	if err := g.db.UpdateGroupMemberInfo(ctx, info.UserId, agrs); err != nil {
		log.ZError(ctx, "update group member info err", err)
	}
}
