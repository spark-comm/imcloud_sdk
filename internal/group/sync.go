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
	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/group"
	"github.com/imCloud/im/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdkerrs"
)

// SyncGroupMember 同步群成员
func (g *Group) SyncGroupMember(ctx context.Context, groupID string) error {
	//获取远程的成员列表
	members, err := g.GetServerGroupMembers(ctx, groupID)
	if err != nil {
		return err
	}
	//获取本地的成员列表
	localData, err := g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 9999999)
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "SyncGroupMember Info", "groupID", groupID, "members", len(members), "localData", len(localData))

	//util.Batch(ServerGroupMemberToLocalGroupMember, members) 远程数据序列化为本地结构
	err = g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, members), localData, nil)
	if err != nil {
		return err
	}
	//if len(members) != len(localData) {
	log.ZInfo(ctx, "SyncGroupMember Sync Group Member Count", "groupID", groupID, "members", len(members), "localData", len(localData))
	//获取远程群组数据（单条）
	gs, err := g.GetSpecifiedGroupsInfo(ctx, []string{groupID})
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "SyncGroupMember GetGroupsInfo", "groupID", groupID, "len", len(gs), "gs", gs)
	if len(gs) > 0 {
		v := gs[0]
		count := int32(len(members))
		if v.MemberCount != count {
			v.MemberCount = int32(len(members))
			if v.GroupType == constant.SuperGroupChatType {
				if err := g.db.UpdateSuperGroup(ctx, v); err != nil {
					//return err
					log.ZError(ctx, "SyncGroupMember UpdateSuperGroup", err, "groupID", groupID, "info", v)
				}
			} else {
				if err := g.db.UpdateGroup(ctx, v); err != nil {
					log.ZError(ctx, "SyncGroupMember UpdateGroup", err, "groupID", groupID, "info", v)
				}
			}
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			log.ZInfo(ctx, "SyncGroupMember OnGroupInfoChanged", "groupID", groupID, "data", string(data))
			g.listener.OnGroupInfoChanged(string(data))
		}
	}
	//}
	return nil
}

// SyncJoinedGroup 同步加入的群
func (g *Group) SyncJoinedGroup(ctx context.Context) error {
	_, err := g.syncJoinedGroup(ctx)
	if err != nil {
		return err
	}
	return err
}

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

// GetServerJoinGroup  获取服务端加入的群
func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*groupv1.GroupInfo, error) {
	fn := func(resp *groupv1.UserJoinGroupInfoList) []*groupv1.GroupInfo { return resp.Groups }
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

// GetServerAdminGroupApplicationList 获取服务端加群申请
func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &group.GetGroupApplicationListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, fn)
}

// GetServerSelfGroupApplication 获取服务端的自己的加群请求
func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, fn)
}

// GetServerGroupMembers 远程获取群成员
func (g *Group) GetServerGroupMembers(ctx context.Context, groupID string) ([]*groupv1.MembersInfo, error) {
	req := &group.GetGroupMemberListReq{GroupID: groupID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *groupv1.MemberListForSDKReps) []*groupv1.MembersInfo { return resp.Members }
	return util.GetPageAll(ctx, constant.GetGroupMemberListRouter, req, fn)
}

func (g *Group) syncGroupStatus(ctx context.Context, groupID string) error {
	svrGroup, err := g.getGroupsInfoFromSvr(ctx, []string{groupID})
	if err != nil {
		return err
	}
	if len(svrGroup) < 1 {
		return sdkerrs.ErrGroupIDNotFound.Wrap("server not this group")
	}
	return g.db.UpdateGroup(ctx, ServerGroupToLocalGroup(svrGroup[0]))
}

// syncJoinedGroupByID 根据id同步群信息
func (g *Group) syncJoinedGroupByID(ctx context.Context, id ...string) error {
	//根据id获取群（远程数据）
	groups, err := g.getGroupsInfoFromSvr(ctx, id)
	if err != nil {
		return err
	}
	//根据id获取本地群数据
	localData, err := g.db.GetGroupInfoByGroupIDs(ctx, id...)
	if err != nil {
		return err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return err
	}
	return nil
}

// syncGroupAndMember 同步群和群成员
func (g *Group) syncGroupAndMember(ctx context.Context, groupId string, memberId ...string) {
	//同步群数据
	if err := g.syncJoinedGroupByID(ctx, groupId); err != nil {
		log.ZDebug(ctx, "syncGroupAndMember->syncJoinedGroupByID err", err)
	}
	//同步群成员数据
}
func (g *Group) syncUserReqGroupInfo(ctx context.Context, fromUserID, groupID string) error {
	//获取用户加入单个群的申请信息
	req := groupv1.UserJoinGroupRequestReq{
		GroupID: groupID,
		UserID:  fromUserID,
	}
	reqInfos, err := util.CallApi[groupv1.UserJoinGroupRequestReps](ctx, constant.GetJoinGroupRequestDetailRouter, &req)
	if err != nil {
		return err
	}
	localGroupRequest := ServerGroupRequestToLocalGroupRequest(&groupv1.GroupRequestInfo{
		GroupID:       groupID,
		CreateTime:    reqInfos.CreateTime,
		GroupName:     reqInfos.GroupName,
		Notification:  reqInfos.Notification,
		Introduction:  reqInfos.Introduction,
		GroupFaceURL:  reqInfos.GroupFaceURL,
		Status:        reqInfos.GroupStatus,
		GroupType:     reqInfos.GroupType,
		GroupCode:     reqInfos.GroupCode,
		OwnerUserID:   reqInfos.OwnerUserID,
		CreatorUserID: reqInfos.CreatorUserID,
		MemberCount:   int32(reqInfos.MemberCount),
		UserID:        fromUserID,
		Nickname:      reqInfos.Nickname,
		UserFaceURL:   reqInfos.FaceURL,
		Gender:        reqInfos.Gender,
		Code:          reqInfos.Code,
		HandleResult:  reqInfos.HandleResult,
		ReqMsg:        reqInfos.ReqMsg,
		HandledMsg:    reqInfos.HandleMsg,
		ReqTime:       reqInfos.ReqTime,
		HandledTime:   reqInfos.HandleTime,
		HandleUserID:  reqInfos.HandleUserID,
		JoinSource:    reqInfos.JoinSource,
		InviterUserID: reqInfos.InviterUserID,
	})
	g.db.UpdateGroupRequest(ctx, localGroupRequest)
	return nil
}
