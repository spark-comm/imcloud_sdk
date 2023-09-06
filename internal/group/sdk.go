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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/db/pg"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"sort"
	"strings"
	"time"

	"github.com/imCloud/im/pkg/common/log"

	groupv1 "github.com/imCloud/api/group/v1"
	"github.com/imCloud/im/pkg/utils"
)

// CreateGroup 创建群
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
	if err := g.syncGroupAndMember(ctx, resp.GroupID); err != nil {
		return nil, err
	}
	return resp, nil
}

// JoinGroup 加入群
func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32) error {
	if err := util.ApiPost(ctx, constant.JoinGroupRouter,
		&groupv1.JoinGroupReq{
			GroupID:  groupID,
			Remark:   reqMsg,
			SourceID: joinSource,
			UserID:   g.loginUserID,
		},
		nil); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplications(ctx, groupID); err != nil {
		return err
	}
	return nil
}

// QuitGroup 退出群聊
func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.QuitGroupRouter, groupv1.QuitAllGroupsReq{
		UserID: g.loginUserID,
	}, nil); err != nil {
		return err
	}
	//删除群信息
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	return nil
}

// DismissGroup 解散群
func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := util.ApiPost(ctx, constant.DismissGroupRouter, &groupv1.IsGroupMemberReq{
		GroupID: groupID,
		UserID:  g.loginUserID,
	}, nil); err != nil {
		return err
	}
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	return nil
}

// ChangeGroupMute 群禁言状态改变
func (g *Group) ChangeGroupMute(ctx context.Context, groupID string, isMute bool) (err error) {
	if isMute {
		err = util.ApiPost(ctx, constant.MuteGroupRouter, &groupv1.IsGroupMemberReq{
			GroupID: groupID,
			UserID:  g.loginUserID,
		}, nil)
	} else {
		err = util.ApiPost(ctx, constant.CancelMuteGroupRouter, &groupv1.CancelMuteGroupMemberReq{
			GroupID: groupID,
			UserID:  g.loginUserID,
		}, nil)
	}
	if err != nil {
		return err
	}
	//更新群状态
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	return nil
}

// ChangeGroupMemberMute 成员禁言
func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds int) (err error) {
	if mutedSeconds == 0 {
		err = util.ApiPost(ctx, constant.CancelMuteGroupMemberRouter, &groupv1.CancelMuteGroupMemberReq{
			GroupID: groupID,
			UserID:  userID,
			PUserID: g.loginUserID,
		}, nil)
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
	if err := g.syncGroupMembers(ctx, groupID, userID); err != nil {
		return err
	}
	return nil
}

// SetGroupMemberRoleLevel 设置群成员等级
func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfoReq{
		GroupID:   groupID,
		UserID:    userID,
		RoleLevel: rune(roleLevel),
		PUserID:   g.loginUserID,
	})
}

// SetGroupMemberNickname 设置群昵称
func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string) error {
	return g.SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfoReq{
		GroupID:  groupID,
		UserID:   userID,
		Nickname: groupMemberNickname,
		PUserID:  g.loginUserID,
	})
}

// SetBackgroundUrl 设置聊天背景图片
func (g *Group) SetBackgroundUrl(ctx context.Context, groupID, backgroundUrl string) error {
	err := g.SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfoReq{
		GroupID:       groupID,
		UserID:        g.loginUserID,
		BackgroundUrl: backgroundUrl,
		PUserID:       g.loginUserID,
	})
	if err != nil {
		return err
	}
	//设置会话的聊天背景
	common.TriggerCmdUpdateConversationBackgroundURL(ctx, g.getConversationIDBySessionType(groupID, constant.SuperGroup), backgroundUrl, g.conversationCh)
	return nil
}

// SetGroupMemberInfo 设置群信息
func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *groupv1.SetGroupMemberInfoReq) error {
	groupMemberInfo.PUserID = g.loginUserID
	if err := util.ApiPost(ctx, constant.SetGroupMemberInfoRouter, &groupMemberInfo, nil); err != nil {
		return err
	}
	return g.syncGroupMembers(ctx, groupMemberInfo.GroupID, groupMemberInfo.PUserID)
}

// GetJoinedGroupList 获取加入的群列表
func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

// GetSpecifiedGroupsInfo 根据id获取群信息
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

// SetGroupVerification 获取群是否需要验证
func (g *Group) SetGroupVerification(ctx context.Context, groupID string, verification int32) error {
	return g.SetGroupInfo(ctx, &groupv1.EditGroupProfileRequest{
		GroupID:          groupID,
		NeedVerification: int64(verification),
	})
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

// SetGroupInfo 更新群信息
func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *groupv1.EditGroupProfileRequest) error {
	groupInfo.UserID = g.loginUserID
	if err := util.ApiPost(ctx, constant.SetGroupInfoRouter, &groupInfo, nil); err != nil {
		return err
	}
	return g.SyncGroups(ctx, groupInfo.GroupID)
}

// GetGroupMemberList 获取群成员列表
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
	// 检查是否同步过
	// i, err := g.db.GetGroupMemberCount(ctx, groupID)
	// if i == 0 || err != nil {
	// 	if err := g.SyncAllGroupMember(ctx, groupID); err != nil {
	// 		return nil, err
	// 	}
	// }
	// return g.db.GetGroupMemberListSplit(
	// 	ctx,
	// 	groupID,
	// 	filter,
	// 	int(pg.BuildOffsetByPage(int(offset), int(count))),
	// 	int(count))
	//检查是否同步过
	// 检查是否同步过
	// i, err := g.db.GetGroupMemberCount(ctx, groupID)
	// if i == 0 || err != nil {
	// 	//从远端读取
	// 	members, err := g.GetServerFirstPageGroupMembers(ctx, groupID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	//通知同步群成员
	// 	common.TriggerCmdSyncGroupMembers(ctx, groupID, g.groupCh)
	// 	return util.Batch(ServerGroupMemberToLocalGroupMember, members), nil
	// } else {

	// }
}

// GetGroupMemberOwnerAndAdmin 获取群主和管理员
func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

// GetGroupMemberListByJoinTimeFilter 查询群成员根据时间过滤
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

// GetSpecifiedGroupMembersInfo 获取群成员信息
func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	return g.GetGroupMemberFromLocal2Svr(ctx, groupID, userIDList)
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.KickGroupMemberRouter, &groupv1.KickGroupMemberReq{
		GroupID:          groupID,
		KickedUserIdList: userIDList,
		UserID:           g.loginUserID,
		HandledMsg:       reason,
	},
		nil); err != nil {
		return err
	}
	return g.deleteGroupMembers(ctx, groupID, userIDList...)
}

// TransferGroupOwner 转让群主
func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	if err := util.ApiPost(ctx, constant.TransferGroupRouter, &groupv1.TransferGroupReq{
		GroupID:        groupID,
		NewOwnerUserID: newOwnerUserID,
		UserID:         g.loginUserID,
	},
		nil); err != nil {
		return err
	}
	if err := g.syncGroupMembers(ctx, groupID, newOwnerUserID, g.loginUserID); err != nil {
		return err
	}
	return nil
}

// InviteUserToGroup 邀请用户进群
func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := util.ApiPost(ctx, constant.InviteUserToGroupRouter, &groupv1.InviteUserToGroupReq{
		GroupID:           groupID,
		InvitedUserIdList: userIDList,
		Reason:            reason,
		UserID:            g.loginUserID,
	},
		nil); err != nil {
		return nil
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.syncGroupMembers(ctx, groupID, userIDList...); err != nil {
		return err
	}
	return nil
}

func (g *Group) GetGroupApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

// GetPageGroupApplicationListAsRecipient 分页获取加群请求数据
func (g *Group) GetPageGroupApplicationListAsRecipient(ctx context.Context, groupId string, no, size int64) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetPageGroupApplicationListAsRecipient(ctx, groupId, &pg.Page{
		NO:   no,
		Size: size,
	})
}

// GetGroupApplicationListAsApplicant 获取发出的群请求
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
	//申请状态同步
	return g.syncUserReqGroupInfo(ctx, req.FromUserID, req.GroupID)
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

func (g *Group) KickGroupMemberList(ctx context.Context, searchParam *sdk_params_callback.GetKickGroupListReq) (sdk_params_callback.SearchKickGroupListInfoRes, error) {
	if searchParam.PageNum == 0 {
		searchParam.PageNum = 1
	}
	if searchParam.PageSize == 0 {
		searchParam.PageSize = 20
	}
	result := sdk_params_callback.SearchKickGroupListInfoRes{}
	kickMemberList, total, err := g.db.SearchKickMemberList(ctx, sdk_params_callback.GetKickGroupListReq{
		GroupID:  searchParam.GroupID,
		IsManger: searchParam.IsManger,
		Name:     searchParam.Name,
		PageSize: searchParam.PageSize,
		PageNum:  searchParam.PageNum,
		UserID:   g.loginUserID,
	})
	if err != nil {
		return result, err
	}
	result.Total = total
	result.KickGroupList = kickMemberList
	return result, nil
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

func (g *Group) GetUserOwnerJoinRequestNum(ctx context.Context) (groupv1.GetOwnerJoinRequestNumReps, error) {
	resp := groupv1.GetOwnerJoinRequestNumReps{}
	if err := util.ApiPost(ctx, constant.GetUserOwnerJoinRequestNumRouter,
		&groupv1.GetOwnerJoinRequestNumReq{
			UserID: g.loginUserID,
		},
		&resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (g *Group) GetAppointGroupRequestInfo(ctx context.Context, groupID string, offset, count int) ([]model_struct.LocalGroupRequest, error) {
	if offset == 0 {
		offset = 1
	}
	if count == 0 {
		count = 20
	}
	return g.db.GetOwnerOrAdminGroupReqInfo(
		ctx,
		groupID,
		(offset-1)*count,
		count,
	)
}

// syncDelGroup 同步删除群
func (g *Group) syncDelGroup(ctx context.Context, groupID string) error {
	localData, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync group", "data from local", localData)
	return g.groupSyncer.Delete(ctx, []*model_struct.LocalGroup{localData}, nil)
}

// SearchGroupInfo 搜索群
func (g *Group) SearchGroupInfo(ctx context.Context, keyWord string, pageSize, pageNum int64) (groupv1.SearchGroupInfoResp, error) {
	resp := groupv1.SearchGroupInfoResp{}
	err := util.ApiPost(ctx, constant.SearchGroupInfoRouter, &groupv1.SearchGroupInfoReq{
		KeyWord:  keyWord,
		PageSize: pageSize,
		PageNum:  pageNum,
	}, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DelGroupConversation 删除群会话
func (g *Group) DelGroupConversation(ctx context.Context, groupID string) {
	//删除会话
	conversationID := utils.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)
	err := common.TriggerCmdDeleteConversationAndMessage(
		ctx,
		groupID,
		conversationID,
		constant.SuperGroupChatType,
		g.conversationCh)
	if err != nil {
		log.ZDebug(ctx, "QuitGroup  after delete conversation err", err)
	}
}

// getConversationIDBySessionType 获取会话类型
func (g *Group) getConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case constant.SingleChatType:
		l := []string{g.loginUserID, sourceID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		return "g_" + sourceID // group chat
	case constant.SuperGroupChatType:
		return "sg_" + sourceID // super group chat
	case constant.NotificationChatType:
		return "sn_" + sourceID // server notification chat
	}
	return ""
}
