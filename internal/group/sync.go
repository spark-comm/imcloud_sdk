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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"

	commonPb "github.com/imCloud/api/common"
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

// GetServerJoinGroup  获取服务端加入的群
func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*groupv1.GroupInfo, error) {
	fn := func(resp *groupv1.UserJoinGroupInfoList) []*groupv1.GroupInfo { return resp.Groups }
	req := &groupv1.GetJoinedGroupListR{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.UserJoinGroupInfoList{}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, resp, fn)
}

// GetServerFirstPageJoinGroup  从服务器获取第一页数据
func (g *Group) GetServerFirstPageJoinGroup(ctx context.Context) ([]*groupv1.GroupInfo, error) {
	fn := func(resp *groupv1.UserJoinGroupInfoList) []*groupv1.GroupInfo { return resp.Groups }
	req := &groupv1.GetJoinedGroupListR{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.UserJoinGroupInfoList{}
	return util.GetFirstPage(ctx, constant.GetJoinedGroupListRouter, req, resp, fn)
}

// GetServerAdminGroupApplicationList 获取服务端加群申请
func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &groupv1.GetJoinedGroupListR{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.GetRecvGroupApplicationListResp{}
	return util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, resp, fn)
}

// GetServerSelfGroupApplication 获取服务端的自己的加群请求
func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetRecvGroupApplicationListResp) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &groupv1.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.GetRecvGroupApplicationListResp{}
	return util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, resp, fn)
}

// GetServerAdminGroupUntreatedApplicationList 获取服务端未处理加群申请
func (g *Group) GetServerAdminGroupUntreatedApplicationList(ctx context.Context) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetUntreatedGroupApplicationListReply) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &groupv1.GetUntreatedRecvGroupApplicationList{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.GetUntreatedGroupApplicationListReply{}
	return util.GetPageAll(ctx, constant.GetUntreatedRecvGroupApplicationListRouter, req, resp, fn)
}

// GetServerGroupMembers 远程获取群成员
func (g *Group) GetServerGroupMembers(ctx context.Context, groupID string) ([]*groupv1.MembersInfo, error) {
	req := &groupv1.GroupMemberListReq{GroupID: groupID, Pagination: &commonPb.RequestPagination{ShowNumber: 100}}
	fn := func(resp *groupv1.MemberListForSDKReps) []*groupv1.MembersInfo { return resp.Members }
	resp := &groupv1.MemberListForSDKReps{}
	return util.GetPageAll(ctx, constant.GetGroupMemberListRouter, req, resp, fn)
}

// GetServerFirstPageGroupMembers  从服务器获取第一也群数据
func (g *Group) GetServerFirstPageGroupMembers(ctx context.Context, groupID string) ([]*groupv1.MembersInfo, error) {
	req := &groupv1.GroupMemberListReq{GroupID: groupID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *groupv1.MemberListForSDKReps) []*groupv1.MembersInfo { return resp.Members }
	resp := &groupv1.MemberListForSDKReps{}
	return util.GetFirstPage(ctx, constant.GetGroupMemberListRouter, req, resp, fn)
}

// SyncGroups 根据id同步群信息
func (g *Group) SyncGroups(ctx context.Context, id ...string) error {
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
	// 触发群改变通知
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupInfo))
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

// GetDesignatedGroupMembers 获取指定群成员信息
func (g *Group) GetDesignatedGroupMembers(ctx context.Context, groupID string, userIDs ...string) ([]*groupv1.MembersInfo, error) {
	//resp := &groupv1.MemberByIdsRes{}
	//if err := util.ApiPost(ctx, constant.GetGroupMemberByIdsRouter, &groupv1.MemberByIdsReq{
	//	GroupID: groupID,
	//	UserIDs: userIDs,
	//}, resp); err != nil {
	//	return nil, err
	//}
	resp, err := util.ProtoApiPost[groupv1.MemberByIdsReq, groupv1.MemberByIdsRes](
		ctx,
		constant.GetGroupMemberByIdsRouter,
		&groupv1.MemberByIdsReq{
			GroupID: groupID,
			UserIDs: userIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// deleteGroupMembers 删除指定群成员
func (g *Group) deleteGroupMembers(ctx context.Context, groupID string, memberId ...string) error {
	if err := g.db.DeleteGroupMembers(ctx, groupID, memberId...); err != nil {
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
	localData, err := g.db.GetSendGroupApplication(ctx)
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

func (g *Group) syncUserReqGroupInfo(ctx context.Context, fromUserID, groupID string) error {
	//获取用户加入单个群的申请信息
	//req := groupv1.UserJoinGroupRequestReq{
	//	GroupID: groupID,
	//	UserID:  fromUserID,
	//}
	//reqInfos, err := util.CallApi[groupv1.UserJoinGroupRequestReps](ctx, constant.GetJoinGroupRequestDetailRouter, &req)
	reqInfos := &groupv1.UserJoinGroupRequestReps{}
	err := util.CallPostApi[*groupv1.UserJoinGroupRequestReq, *groupv1.UserJoinGroupRequestReps](
		ctx, constant.GetJoinGroupRequestDetailRouter,
		&groupv1.UserJoinGroupRequestReq{
			GroupID: groupID,
			UserID:  fromUserID},
		reqInfos,
	)
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

// InitSyncGroupData 初始化同步群数据
func (g *Group) InitSyncGroupData(ctx context.Context) error {
	//从群中获取数据
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		log.ZError(ctx, "InitSyncGroupData err", err)
	}
	//获取远程数据
	fn := func(resp *groupv1.GetSyncGroupResp) []*groupv1.GetSyncGroupRespInfo { return resp.List }
	req := &groupv1.GetJoinedGroupListR{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.GetSyncGroupResp{}
	serveGroups, err := util.GetPageAll(ctx, constant.GetSyncGroupInfoList, req, resp, fn)
	var groupLists = make([]*groupv1.GroupInfo, 0)
	for _, list := range serveGroups {
		groupLists = append(groupLists, &groupv1.GroupInfo{
			GroupID:     list.GroupID,
			GroupName:   list.NickName,
			FaceURL:     list.FaceURL,
			MemberCount: list.GroupNumber,
			IsComplete:  IsNotComplete,
		})
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groupLists), groups, nil); err != nil {
		return err
	}
	//延迟30秒同步全量数据
	//g.syncGroupQueue.Push(1, time.Second*30)
	// 发出同步的通知
	common.TriggerCmdJoinGroup(ctx, g.groupCh)
	return nil
}

func (g *Group) GetUserMemberInfoInGroup(ctx context.Context) error {
	//初始化同步用户的所有群中成员信息
	info, _ := g.db.GetOwnerGroupMemberInfo(ctx, g.loginUserID)
	resp := &groupv1.GetUserInGroupMemberResp{}
	//获取远端数据
	err := util.CallPostApi[*groupv1.GetUserInGroupMemberReq, *groupv1.GetUserInGroupMemberResp](
		ctx,
		constant.GetUserMemberInfoInGroup,
		&groupv1.GetUserInGroupMemberReq{
			UserID: g.loginUserID,
		},
		resp,
	)
	if err != nil {
		return err
	}
	if len(resp.List) == 0 {
		return nil
	}
	//同步数据
	if err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, resp.List), info, nil); err != nil {
		return err
	}
	return nil
}
