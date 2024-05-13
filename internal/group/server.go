package group

import (
	"context"
	commonPb "github.com/imCloud/api/common"
	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/proto/group"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
)

// GetUserMemberInfoInGroup 同步用户在所有群中的信息
func (g *Group) GetUserMemberInfoInGroup(ctx context.Context) ([]*groupv1.MembersInfo, error) {
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
		return nil, err
	}
	return resp.List, nil
}

func (g *Group) GetServerPageGroupMembersInfo(ctx context.Context, groupID string, filter, offset, count int32) ([]*groupv1.MembersInfo, error) {
	fn := func(resp *groupv1.MemberListForSDKReps) []*groupv1.MembersInfo { return resp.Members }
	resp := &groupv1.MemberListForSDKReps{}
	return util.GetFirstPage(ctx, constant.GetGroupMemberListRouter, &groupv1.GroupMemberListReq{
		GroupID: groupID,
		Filter:  filter,
		Pagination: &commonPb.RequestPagination{
			PageNumber: offset,
			ShowNumber: count,
		}}, resp, fn)
}

func (g *Group) GetServerAllGroupMembersInfo(ctx context.Context, groupID string, count int32) ([]*groupv1.MembersInfo, error) {
	fn := func(resp *groupv1.MemberListForSDKReps) []*groupv1.MembersInfo { return resp.Members }
	resp := &groupv1.MemberListForSDKReps{}
	return util.GetPageAll(ctx, constant.GetGroupMemberListRouter, &groupv1.GroupMemberListReq{
		GroupID: groupID,
		Pagination: &commonPb.RequestPagination{
			ShowNumber: count,
		}}, resp, fn)
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

// GetSyncGroup 获取同步的用户数据
func (g *Group) GetSyncGroup(ctx context.Context) ([]*groupv1.BaseGroupInfo, error) {
	//获取远程数据
	fn := func(resp *groupv1.SyncUserJoinGroupInfoResp) []*groupv1.BaseGroupInfo { return resp.List }
	req := &groupv1.GetJoinedGroupListR{FromUserID: g.loginUserID, Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.SyncUserJoinGroupInfoResp{}
	serveGroups, err := util.GetPageAll(ctx, constant.GetSyncGroupInfoList, req, resp, fn)
	if err != nil {
		return nil, err
	}
	return serveGroups, nil
}

// SyncGroupMemberInfo 同步群成员
func (g *Group) SyncGroupMemberInfo(ctx context.Context, groupID string) ([]*groupv1.BaseGroupMemberInfo, error) {
	fn := func(resp *groupv1.SyncGroupMemberInfoResp) []*groupv1.BaseGroupMemberInfo {
		return resp.List
	}
	req := &groupv1.SyncGroupMemberInfoReq{
		GroupId: groupID,
		Pagination: &commonPb.RequestPagination{
			ShowNumber: 1000,
		}}
	resp := &groupv1.SyncGroupMemberInfoResp{}
	return util.GetPageAll(ctx, constant.SyncGroupMemberInfoRouter, req, resp, fn)
}

// GetServerAdminOneGroupUntreatedApplicationList 未处理的加群请求
func (g *Group) GetServerAdminOneGroupUntreatedApplicationList(ctx context.Context, groupID string) ([]*groupv1.GroupRequestInfo, error) {
	fn := func(resp *groupv1.GetUntreatedGroupApplicationListReply) []*groupv1.GroupRequestInfo {
		return resp.GroupRequests
	}
	req := &groupv1.GetUntreatedRecvGroupApplicationList{
		GroupID:    groupID,
		FromUserID: g.loginUserID,
		Pagination: &commonPb.RequestPagination{}}
	resp := &groupv1.GetUntreatedGroupApplicationListReply{}
	return util.GetPageAll(ctx, constant.GetUntreatedRecvGroupApplicationListRouter, req, resp, fn)
}

// SyncUserJoinGroupInfoByTime 根据时间点同步群信息
// par dataTime 时间问题
//
//	返回 新增群 修改群 删除群
func (g *Group) SyncUserJoinGroupInfoByTime(ctx context.Context, dataTime map[string]int64) ([]*groupv1.BaseGroupInfo, []*groupv1.BaseGroupInfo, []string, error) {
	res := &groupv1.SyncUserJoinGroupInfoByTimeReply{}
	err := util.CallPostApi[*groupv1.SyncUserJoinGroupInfoByTimeReq, *groupv1.SyncUserJoinGroupInfoByTimeReply](
		ctx, constant.SyncUserJoinGroupInfoByTimeRouter,
		&groupv1.SyncUserJoinGroupInfoByTimeReq{UserID: g.loginUserID, TimeData: dataTime},
		res,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return res.AddList, res.AddList, res.DelIds, nil
}

// SyncGroupMemberInfoUpdateTime 根据时间点同步群成员信息
// par dataTime 时间问题
//
//	返回 新增群成员 修改群成员 删除群成员
func (g *Group) SyncGroupMemberInfoUpdateTime(ctx context.Context, groupId string, dataTime map[string]int64) ([]*groupv1.BaseGroupMemberInfo, []*groupv1.BaseGroupMemberInfo, []string, error) {
	res := &groupv1.SyncGroupMemberInfoUpdateTimeReply{}
	err := util.CallPostApi[*groupv1.SyncGroupMemberInfoUpdateTimeReq, *groupv1.SyncGroupMemberInfoUpdateTimeReply](
		ctx, constant.SyncGroupMemberInfoUpdateTimeRouter,
		&groupv1.SyncGroupMemberInfoUpdateTimeReq{GroupId: groupId, TimeData: dataTime},
		res,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return res.AddList, res.AddList, res.DelIds, nil
}

// GetDesignatedGroupMembers 获取指定群成员信息
func (g *Group) GetDesignatedGroupMembers(ctx context.Context, groupID string, userIDs ...string) ([]*groupv1.MembersInfo, error) {
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

// FindFullGroupInfo 获取群完整信息忽略删除和解散
func (g *Group) FindFullGroupInfo(ctx context.Context, groupIDs ...string) ([]*groupv1.BaseGroupInfo, error) {
	resp, err := util.ProtoApiPost[groupv1.GetFullGroupInfoReq, groupv1.GetFullGroupInfoReply](
		ctx,
		constant.FindFullGroupInfoRouter,
		&groupv1.GetFullGroupInfoReq{
			GroupIds: groupIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
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

// syncUserReqGroupInfo 同步加群请求
func (g *Group) syncUserReqGroupInfo(ctx context.Context, fromUserID, groupID string) error {
	//获取用户加入单个群的申请信息
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
