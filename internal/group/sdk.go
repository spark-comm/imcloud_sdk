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
	"github.com/imCloud/im/pkg/common/log"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"time"

	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/utils"
)

//// deprecated use CreateGroup
//funcation (g *Group) CreateGroup(ctx context.Context, groupBaseInfo sdk_params_callback.CreateGroupBaseInfoParam, memberList sdk_params_callback.CreateGroupMemberRoleParam) (*sdkws.GroupInfo, error) {
//	req := &group.CreateGroupReq{
//		GroupInfo: &sdkws.GroupInfo{
//			GroupName:    groupBaseInfo.GroupName,
//			Notification: groupBaseInfo.Notification,
//			Introduction: groupBaseInfo.Introduction,
//			FaceURL:      groupBaseInfo.FaceURL,
//			Ex:           groupBaseInfo.Ex,
//			GroupType:    groupBaseInfo.GroupType,
//		},
//	}
//	if groupBaseInfo.NeedVerification != nil {
//		req.GroupInfo.NeedVerification = *groupBaseInfo.NeedVerification
//	}
//	for _, info := range memberList {
//		switch info.RoleLevel {
//		case constant.GroupOrdinaryUsers:
//			req.InitMembers = append(req.InitMembers, info.UserID)
//		case constant.GroupOwner:
//			req.OwnerUserID = info.UserID
//		case constant.GroupAdmin:
//			req.AdminUserIDs = append(req.AdminUserIDs, info.UserID)
//		default:
//			return nil, sdkerrs.ErrArgs.Wrap(fmt.Sprintf("CreateGroup: invalid role level %d", info.RoleLevel))
//		}
//	}
//	return g.CreateGroup(ctx, req)
//}

func (g *Group) CreateGroup(ctx context.Context, req *groupv1.CrateGroupReq) (*groupv1.GroupInfo, error) {
	if req.CreatorUserID == "" {
		req.CreatorUserID = g.loginUserID
	}
	if req.GroupType != constant.WorkingGroup {
		return nil, sdkerrs.ErrGroupType
	}
	resp, err := util.CallApi[groupv1.GroupInfo](ctx, constant.CreateGroupRouter, req)
	if err != nil {
		return nil, err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return nil, err
	}
	if err := g.SyncGroupMember(ctx, resp.GroupID); err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32) error {
	if err := util.ApiPost(ctx, constant.JoinGroupRouter,
		&groupv1.JoinGroupReq{
			GroupID:  groupID,
			Remark:   reqMsg,
			SourceID: joinSource,
			UserID:   g.loginUserID,
		},
		//&group.JoinGroupReq{
		//	GroupID:       groupID,
		//	ReqMessage:    reqMsg,
		//	JoinSource:    joinSource,
		//	InviterUserID: g.loginUserID,
		//},
		nil); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplication(ctx); err != nil {
		return err
	}
	//if err := g.SyncJoinedGroup(ctx); err != nil {
	//	return err
	//}
	//if err := g.SyncGroupMember(ctx, groupID); err != nil {
	//	return err
	//}
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.QuitGroupRouter, groupv1.CancelMuteGroupMemberReq{
		GroupID: groupID,
		UserID:  g.loginUserID,
	},
		//&group.QuitGroupReq{GroupID: groupID},
		nil); err != nil {
		return err
	}
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	//if err := g.SyncGroupMember(ctx, groupID); err != nil {
	//	return err
	//}
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.DismissGroupRouter, &groupv1.IsGroupMemberReq{
		GroupID: groupID,
		UserID:  g.loginUserID,
	}, nil); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMute(ctx context.Context, groupID string, isMute bool) (err error) {
	if isMute {
		err = util.ApiPost(ctx, constant.MuteGroupRouter, &groupv1.IsGroupMemberReq{
			GroupID: groupID,
			UserID:  g.loginUserID,
		},
			//group.MuteGroupReq{GroupID: groupID},
			nil)
	} else {
		err = util.ApiPost(ctx, constant.CancelMuteGroupRouter, &groupv1.CancelMuteGroupMemberReq{
			GroupID: groupID,
			UserID:  g.loginUserID,
		},
			//group.CancelMuteGroupReq{GroupID: groupID},
			nil)
	}
	if err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds int) (err error) {
	if mutedSeconds == 0 {
		err = util.ApiPost(ctx, constant.CancelMuteGroupMemberRouter, &groupv1.CancelMuteGroupMemberReq{
			GroupID: groupID,
			UserID:  userID,
			PUserID: g.loginUserID,
		},
			//group.CancelMuteGroupMemberReq{GroupID: groupID, UserID: userID},
			nil)
	} else {
		err = util.ApiPost(ctx, constant.MuteGroupMemberRouter, &groupv1.MuteGroupMemberReq{
			GroupID:      groupID,
			UserID:       userID,
			MutedSeconds: int64(mutedSeconds),
			PUserID:      g.loginUserID,
		},
			//&group.MuteGroupMemberReq{GroupID: groupID, UserID: userID, MutedSeconds: uint32(mutedSeconds)},
			nil)
	}
	if err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfo{
		GroupID:   groupID,
		UserID:    userID,
		RoleLevel: rune(roleLevel),
		PUserID:   g.loginUserID,
	})
}

func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string) error {
	return g.SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfo{
		GroupID:  groupID,
		UserID:   userID,
		Nickname: groupMemberNickname,
		PUserID:  g.loginUserID,
	})
	//&group.SetGroupMemberInfo{GroupID: groupID, UserID: userID, Nickname: wrapperspb.String(groupMemberNickname)})
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *groupv1.SetGroupMemberInfo) error {
	groupMemberInfo.PUserID = g.loginUserID
	if err := util.ApiPost(ctx, constant.SetGroupMemberInfoRouter, &groupv1.SetGroupMemberInfoReq{
		Members: []*groupv1.SetGroupMemberInfo{groupMemberInfo},
	},
		//&group.SetGroupMemberInfoReq{Members: []*group.SetGroupMemberInfo{groupMemberInfo}},
		nil); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupMemberInfo.GroupID)
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	superGroupList, err := g.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return nil, err
	}
	groupIDMap := utils.SliceSet(groupIDs)
	//获取所有群数据（普通群和超级群）
	groups := append(groupList, superGroupList...)
	res := make([]*model_struct.LocalGroup, 0, len(groupIDs))
	for i, v := range groups {
		if _, ok := groupIDMap[v.GroupID]; ok {
			delete(groupIDMap, v.GroupID)
			res = append(res, groups[i])
		}
	}
	if len(groupIDMap) > 0 {
		groups, err := util.CallApi[groupv1.GetGroupInfoResponse](
			ctx,
			constant.GetGroupsInfoRouter,
			&groupv1.GetGroupInfoReq{GroupID: utils.Keys(groupIDMap)})
		if err != nil {
			log.ZError(ctx, "Call GetGroupsInfoRouter", err)
		}
		if groups != nil && len(groups.Data) > 0 {
			for i := range groups.Data {
				groups.Data[i].MemberCount = 0
			}
			//转换为本地的群组数据格式
			res = append(res, util.Batch(ServerGroupToLocalGroup, groups.Data)...)
		}
	}
	return res, nil
}

// SearchGroups 本地数据过滤
func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, sdkerrs.ErrArgs.Wrap("keyword is null or search field all false")
	}
	groups, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(
		ctx,
		param.KeywordList[0],
		param.IsSearchGroupID,
		param.IsSearchGroupName) // todo	param.KeywordList[0]
	if err != nil {
		return nil, err
	}
	return groups, nil
}

//funcation (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdk_params_callback.SetGroupInfoParam, groupID string) error {
//	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{
//		GroupID:          groupID,
//		GroupName:        groupInfo.GroupName,
//		Notification:     groupInfo.Notification,
//		Introduction:     groupInfo.Introduction,
//		FaceURL:          groupInfo.FaceURL,
//		Ex:               groupInfo.Ex,
//		NeedVerification: wrapperspb.Int32Ptr(groupInfo.NeedVerification),
//	})
//}

func (g *Group) SetGroupVerification(ctx context.Context, groupID string, verification int32) error {
	return g.SetGroupInfo(ctx, &groupv1.EditGroupProfileRequest{
		GroupID:          groupID,
		NeedVerification: int64(verification),
	})

	//&sdkws.GroupInfoForSet{
	//GroupID: groupID,
	//NeedVerification: wrapperspb.Int32(verification)})
}

func (g *Group) SetGroupLookMemberInfo(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx,
		&groupv1.EditGroupProfileRequest{
			GroupID:        groupID,
			LookMemberInfo: rule,
			//&sdkws.GroupInfoForSet{GroupID: groupID, LookMemberInfo: wrapperspb.Int32(rule)
		})
}

func (g *Group) SetGroupApplyMemberFriend(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupInfo(ctx,
		&groupv1.EditGroupProfileRequest{
			GroupID:           groupID,
			ApplyMemberFriend: rule,
			//&sdkws.GroupInfoForSet{GroupID: groupID, ApplyMemberFriend: wrapperspb.Int32(rule)
		})
}

func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *groupv1.EditGroupProfileRequest) error {
	groupInfo.UserID = g.loginUserID
	if err := util.ApiPost(ctx, constant.SetGroupInfoRouter, &groupInfo, nil); err != nil {
		return err
	}
	return g.SyncJoinedGroup(ctx)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	if offset == 0 {
		offset = 1
	}
	if count == 0 {
		count = 20
	}
	return g.db.GetGroupMemberListSplit(
		ctx,
		groupID,
		filter,
		int((offset-1)*count),
		int(count))
}

func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}
	if offset == 0 {
		offset = 1
	}
	if count == 0 {
		count = 20
	}
	return g.db.GetGroupMemberListSplitByJoinTimeFilter(
		ctx,
		groupID,
		int((offset-1)*count),
		int(count),
		joinTimeBegin,
		joinTimeEnd,
		userIDs)
}

func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.KickGroupMemberRouter, &groupv1.KickGroupMemberReq{
		GroupID:          groupID,
		KickedUserIdList: userIDList,
		UserID:           g.loginUserID,
		HandledMsg:       reason,
	},
		//&group.KickGroupMemberReq{GroupID: groupID, KickedUserIDs: userIDList, Reason: reason},
		nil); err != nil {
		return err
	}
	return g.SyncGroupMember(ctx, groupID)
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	if err := util.ApiPost(ctx, constant.TransferGroupRouter, &groupv1.TransferGroupReq{
		GroupID:        groupID,
		NewOwnerUserID: newOwnerUserID,
		UserID:         g.loginUserID,
	},
		//&group.TransferGroupOwnerReq{GroupID: groupID, OldOwnerUserID: g.loginUserID, NewOwnerUserID: newOwnerUserID},
		nil); err != nil {
		return err
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &groupv1.InviteUserToGroupReq{
		GroupID:           groupID,
		InvitedUserIdList: userIDList,
		Reason:            reason,
		UserID:            g.loginUserID,
	},
		//group.InviteUserToGroupReq{GroupID: groupID, Reason: reason, InvitedUserIDs: userIDList},
		nil); err != nil {
		return nil
	}
	if err := g.SyncJoinedGroup(ctx); err != nil {
		return err
	}
	if err := g.SyncGroupMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) GetGroupApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

func (g *Group) GetGroupApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	return g.db.GetSendGroupApplication(ctx)
}

func (g *Group) AcceptGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx,
		&groupv1.ApplicationResponseReq{
			GroupID:      groupID,
			FromUserID:   fromUserID,
			HandledMsg:   handleMsg,
			HandleResult: constant.GroupResponseAgree,
			UserID:       g.loginUserID,
		})
	//&group.GroupApplicationResponseReq{GroupID: groupID,
	//FromUserID: fromUserID, HandledMsg: handleMsg,
	//HandleResult: constant.GroupResponseAgree})
}

func (g *Group) RefuseGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx,
		&groupv1.ApplicationResponseReq{
			GroupID:      groupID,
			FromUserID:   fromUserID,
			HandledMsg:   handleMsg,
			HandleResult: constant.GroupResponseRefuse,
			UserID:       g.loginUserID,
		})
	//&group.GroupApplicationResponseReq{GroupID: groupID,
	//	FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseRefuse})
}

func (g *Group) HandlerGroupApplication(ctx context.Context, req *groupv1.ApplicationResponseReq) error {
	if err := util.ApiPost(ctx, constant.AcceptGroupApplicationRouter, req, nil); err != nil {
		return err
	}
	// SyncAdminGroupApplication todo
	return nil
}

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam *sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	if searchParam.Offset == 0 {
		searchParam.Offset = 1
	}
	if searchParam.Count == 0 {
		searchParam.Count = 20
	}
	return g.db.SearchGroupMembersDB(
		ctx,
		searchParam.KeywordList[0],
		searchParam.GroupID,
		searchParam.IsSearchMemberNickname,
		searchParam.IsSearchUserID,
		(searchParam.Offset-1)*searchParam.Count,
		searchParam.Count)
}

func (g *Group) IsJoinGroup(ctx context.Context, groupID string) (bool, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return false, err
	}
	for _, localGroup := range groupList {
		if localGroup.GroupID == groupID {
			return true, nil
		}
	}
	superGroupList, err := g.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return false, err
	}
	for _, localGroup := range superGroupList {
		if localGroup.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}

func (g *Group) KickGroupMemberList(ctx context.Context, searchParam *sdk_params_callback.GetKickGroupListReq) ([]*sdk_params_callback.KickGroupList, error) {
	if searchParam.PageNum == 0 {
		searchParam.PageNum = 1
	}
	if searchParam.PageSize == 0 {
		searchParam.PageSize = 20
	}
	return g.db.SearchKickMemberList(ctx, sdk_params_callback.GetKickGroupListReq{
		GroupID:  searchParam.GroupID,
		IsManger: searchParam.IsManger,
		Name:     searchParam.Name,
		PageSize: searchParam.PageSize,
		PageNum:  searchParam.PageNum,
		UserID:   g.loginUserID,
	})
}

func (g *Group) GetNotInGroupFriendInfoList(ctx context.Context, searchParam *sdk_params_callback.SearchNotInGroupUserReq) (sdk_params_callback.SearchNotInGroupUserInfoRes, error) {
	if searchParam.PageNum == 0 {
		searchParam.PageNum = 1
	}
	if searchParam.PageSize == 0 {
		searchParam.PageSize = 20
	}
	result := sdk_params_callback.SearchNotInGroupUserInfoRes{}
	groupMember, err := g.db.GetGroupMemberListByGroupID(ctx, searchParam.GroupID)
	if err != nil {
		return result, err
	}
	groupFriends := []string{}
	for _, member := range groupMember {
		groupFriends = append(groupFriends, member.UserID)
	}
	info, total, err := g.db.GetNotInListFriendInfo(
		ctx,
		searchParam.Name,
		g.loginUserID,
		groupFriends,
		searchParam.PageSize,
		searchParam.PageNum,
	)
	result.Total = total
	result.Friends = info
	return result, nil
}
