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
	"fmt"
	"time"

	"github.com/spark-comm/imcloud_sdk/pkg/server_api"
	groupmodel "github.com/spark-comm/spark-api/api/common/model/group/v2"
	v2 "github.com/spark-comm/spark-api/api/im_cloud/group/v2"

	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	"github.com/spark-comm/imcloud_sdk/pkg/sdk_params_callback"
	"github.com/spark-comm/imcloud_sdk/pkg/sdkerrs"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/utils"
)

func (g *Group) CreateGroup(ctx context.Context, req *v2.CrateGroupReq) (*groupmodel.GroupInfo, error) {
	if req.GroupType != constant.WorkingGroup {
		return nil, sdkerrs.ErrGroupType
	}
	req.CreatorUserID = g.loginUserID
	groupInfo, err := server_api.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := g.SyncGroups(ctx, groupInfo.GroupID); err != nil {
		return nil, err
	}
	if err := g.SyncAllGroupMember(ctx, groupInfo.GroupID); err != nil {
		return nil, err
	}
	return groupInfo, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32) error {
	if err := server_api.JoinGroup(ctx, g.loginUserID, groupID, reqMsg, joinSource); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplications(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := server_api.QuitGroup(ctx, groupID, g.loginUserID); err != nil {
		return err
	}
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		return err
	}
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := server_api.DismissGroup(ctx, groupID, g.loginUserID); err != nil {
		return err
	}
	if err := g.deleteGroup(ctx, groupID); err != nil {
		return err
	}
	if err := g.db.DeleteGroupAllMembers(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMute(ctx context.Context, groupID string, isMute bool) (err error) {
	err = server_api.ChangeGroupMute(ctx, groupID, g.loginUserID, isMute)
	if err != nil {
		return err
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds int) (err error) {
	err = server_api.ChangeGroupMemberMute(ctx, groupID, userID, g.loginUserID, mutedSeconds)
	if err != nil {
		return err
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncGroupMembers(ctx, groupID, userID); err != nil {
		return err
	}
	return nil
}

func (g *Group) SetGroupMemberRoleLevel(ctx context.Context, groupID, userID string, roleLevel int) error {
	return g.SetGroupMemberInfo(ctx, &v2.SetGroupMemberInfoReq{GroupID: groupID, UserID: userID, RoleLevel: rune(roleLevel)})
}

func (g *Group) SetGroupMemberNickname(ctx context.Context, groupID, userID string, groupMemberNickname string) error {
	return g.SetGroupMemberInfo(ctx, &v2.SetGroupMemberInfoReq{GroupID: groupID, UserID: userID, Nickname: groupMemberNickname})
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *v2.SetGroupMemberInfoReq) error {
	groupMemberInfo.PUserID = g.loginUserID
	if err := server_api.SetGroupMemberInfo(ctx, groupMemberInfo); err != nil {
		return err
	}
	return g.SyncGroupMembers(ctx, groupMemberInfo.GroupID, groupMemberInfo.UserID)
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	groupList, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, err
	}
	groupIDMap := utils.SliceSet(groupIDs)
	res := make([]*model_struct.LocalGroup, 0, len(groupIDs))
	for i, v := range groupList {
		if _, ok := groupIDMap[v.GroupID]; ok {
			delete(groupIDMap, v.GroupID)
			res = append(res, groupList[i])
		}
	}
	if len(groupIDMap) > 0 {
		groups, err := server_api.GetSpecifiedGroupsInfo(ctx, groupIDs)
		if err != nil {
			log.ZError(ctx, "Call GetGroupsInfoRouter", err)
		}
		if groups != nil && len(groups) > 0 {
			for i := range groups {
				groups[i].MemberCount = 0
			}
			res = append(res, groups...)
		}
	}
	fmt.Printf("GetSpecifiedGroupsInfo val:%s", utils.StructToJsonString(res))
	return res, nil
}

func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, sdkerrs.ErrArgs.Wrap("keyword is null or search field all false")
	}
	groups, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(ctx, param.KeywordList[0], param.IsSearchGroupID, param.IsSearchGroupName) // todo	param.KeywordList[0]
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// funcation (g *Group) SetGroupInfo(ctx context.Context, groupInfo *sdk_params_callback.SetGroupInfoParam, groupID string) error {
//	return g.SetGroupInfo(ctx, &sdkws.GroupInfoForSet{
//		GroupID:          groupID,
//		GroupName:        groupInfo.GroupName,
//		Notification:     groupInfo.Notification,
//		Introduction:     groupInfo.Introduction,
//		FaceURL:          groupInfo.FaceURL,
//		Ex:               groupInfo.Ex,
//		NeedVerification: wrapperspb.Int32Ptr(groupInfo.NeedVerification),
//	})
// }

func (g *Group) SetGroupVerification(ctx context.Context, groupID string, verification int32) error {
	return g.SetGroupSwitchInfo(ctx, groupID, v2.GroupSwitchOption_needVerification.String(), verification)
}

func (g *Group) SetGroupLookMemberInfo(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupSwitchInfo(ctx, groupID, v2.GroupSwitchOption_lookMemberInfo.String(), rule)
}

func (g *Group) SetGroupApplyMemberFriend(ctx context.Context, groupID string, rule int32) error {
	return g.SetGroupSwitchInfo(ctx, groupID, v2.GroupSwitchOption_applyMemberFriend.String(), rule)
}

// SetGroupSwitchInfo 设置群开关
func (g *Group) SetGroupSwitchInfo(ctx context.Context, groupID string, field string, ups int32) error {
	if err := server_api.SetGroupSwitchInfo(ctx, groupID, g.loginUserID, field, ups); err != nil {
		return err
	}
	return g.SyncGroups(ctx, groupID)
}

func (g *Group) SetGroupInfo(ctx context.Context, req *v2.EditGroupProfileReq) error {
	req.UserID = g.loginUserID
	if err := server_api.SetGroupInfo(ctx, req); err != nil {
		return err
	}
	return g.SyncGroups(ctx, req.GroupID)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberListSplit(ctx, groupID, filter, int(offset), int(count))
}
func (g *Group) GetGroupMemberListPage(ctx context.Context, groupID string, filter, offset, count int32) (*sdk_params_callback.GroupMemberPage, error) {
	data, total, err := g.db.GetGroupMemberListPage(ctx, groupID, filter, int(offset), int(count))
	if err != nil {
		return nil, sdkerrs.GetDataError.Wrap(fmt.Sprintf("GetGroupMemberListPage err: %v", err))
	}
	return &sdk_params_callback.GroupMemberPage{
		List:  data,
		Total: total,
	}, nil
}
func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}
	return g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
}

func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	if err := server_api.KickGroupMember(ctx, groupID, g.loginUserID, reason, userIDList); err != nil {
		return err
	}
	return g.SyncGroupMembers(ctx, groupID, userIDList...)
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	oldOwner, err := g.db.GetGroupMemberOwner(ctx, groupID)
	if err != nil {
		return err
	}
	if err = server_api.TransferGroupOwner(ctx, groupID, g.loginUserID, newOwnerUserID); err != nil {
		return err
	}
	if err = g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err = g.SyncGroupMembers(ctx, groupID, newOwnerUserID, oldOwner.UserID); err != nil {
		return err
	}
	return nil
}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	if err := server_api.InviteUserToGroup(ctx, groupID, g.loginUserID, reason, userIDList); err != nil {
		return err
	}
	if err := g.SyncGroups(ctx, groupID); err != nil {
		return err
	}
	if err := g.SyncGroupMembers(ctx, groupID, userIDList...); err != nil {
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
	return g.HandlerGroupApplication(ctx, &v2.ApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseAgree})
}

func (g *Group) RefuseGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &v2.ApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseRefuse})
}

func (g *Group) HandlerGroupApplication(ctx context.Context, req *v2.ApplicationResponseReq) error {
	req.UserID = g.loginUserID
	if err := server_api.HandlerGroupApplication(ctx, req); err != nil {
		return err
	}
	// SyncAdminGroupApplication todo
	return nil
}

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam *sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	return g.db.SearchGroupMembersDB(ctx, searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
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
	return false, nil
}
func (g *Group) SearchGroupByCode(ctx context.Context, code string) (*model_struct.LocalGroup, error) {
	return server_api.SearchGroupByCode(ctx, g.loginUserID, code)
}
