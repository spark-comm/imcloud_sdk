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

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*groupv1.GroupInfo, error) {
	fn := func(resp *groupv1.UserJoinGroupInfoList) []*groupv1.GroupInfo { return resp.Groups }
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &group.GetGroupApplicationListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, fn)
}

func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, fn)
}

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
