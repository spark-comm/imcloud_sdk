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

package testv2

import (
	"fmt"
	"testing"

	"github.com/spark-comm/imcloud_sdk/open_im_sdk"
	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/sdk_params_callback"
	groupPb "github.com/spark-comm/spark-api/api/im_cloud/group/v2"
)

func Test_CreateGroupV2(t *testing.T) {
	req := &groupPb.CrateGroupReq{
		//MemberUserIDs: []string{"7299270930"},
		//AdminUserIDs:  []string{"1"},
		//OwnerUserID:   UserID,
		//GroupInfo: &sdkws.GroupInfo{
		//	GroupName: "test",
		//	GroupType: 2,
		//},
	}
	info, err := open_im_sdk.UserForSDK.Group().CreateGroup(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("group info: %s", info.String())
}

func Test_JoinGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().JoinGroup(ctx, "3889561099", "1234", 3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success")
}

func Test_QuitGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().QuitGroup(ctx, "xadxwr24")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("QuitGroup success")
}

func Test_DismissGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().DismissGroup(ctx, "1728503199")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DismissGroup success")
}

func Test_ChangeGroupMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(ctx, "3459296007", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelMuteGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(ctx, "3459296007", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_ChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(ctx, "3459296007", UserID, 10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(ctx, "3459296007", UserID, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CancelChangeGroupMemberMute success", ctx.Value("operationID"))
}

func Test_SetGroupMemberRoleLevel(t *testing.T) {
	// 1:普通成员 2:群主 3:管理员
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberRoleLevel(ctx, "3459296007", "45644221123", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberRoleLevel success", ctx.Value("operationID"))
}

func Test_SetGroupMemberNickname(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberNickname(ctx, "3459296007", "45644221123", "test1234")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberNickname success", ctx.Value("operationID"))
}

func Test_SetGroupMemberInfo(t *testing.T) {
	// 1:普通成员 2:群主 3:管理员
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberInfo(ctx, &groupPb.SetGroupMemberInfoReq{
		GroupID:  "3889561099",
		UserID:   UserID,
		FaceURL:  "https://doc.rentsoft.cn/images/logo.png",
		Nickname: "testupdatename",
		Ex:       "a",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberNickname success", ctx.Value("operationID"))
}

func Test_GetJoinedGroupList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetJoinedGroupList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetJoinedGroupList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_SearchJoinedGroupList(t *testing.T) {
	data, err := open_im_sdk.UserForSDK.Group().SearchJoinedGroupList(ctx, "小", 0, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, localGroup := range data.List {
		t.Logf("%#v", localGroup)
	}
	t.Log(data.Total)
}

func Test_GetSpecifiedGroupsInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(ctx, []string{"test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetGroupsInfo: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetGroupApplicationListAsRecipient(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsRecipient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetRecvGroupApplicationList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetGroupApplicationListAsApplicant(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsApplicant(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetSendGroupApplicationList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_AcceptGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().AcceptGroupApplication(ctx, "3459296007", "863454357", "test accept")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_RefuseGroupApplication(t *testing.T) {
	t.Log("operationID:", ctx.Value("operationID"))
	err := open_im_sdk.UserForSDK.Group().RefuseGroupApplication(ctx, "3459296007", "863454357", "test refuse")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success")
}

func Test_HandlerGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().HandlerGroupApplication(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_SearchGroupMembers(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().SearchGroupMembers(ctx, &sdk_params_callback.SearchGroupMembersParam{
		GroupID:                "237875325046784",
		KeywordList:            []string{"177"},
		IsSearchUserID:         true,
		IsSearchMemberNickname: true,
		Offset:                 0,
		Count:                  10,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SearchGroupMembers: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_KickGroupMember(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().KickGroupMember(ctx, "3459296007", "test", []string{"863454357"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

func Test_TransferGroupOwner(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().TransferGroupOwner(ctx, "1728503199", "5226390099")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("TransferGroupOwner success", ctx.Value("operationID"))
}

func Test_InviteUserToGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().InviteUserToGroup(ctx, "3459296007", "test", []string{"45644221123"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

//func Test_SyncGroup(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Group().SyncGroupMember(ctx, "3179997540")
//	if err != nil {
//		t.Fatal(err)
//	}
//	time.Sleep(time.Second * 100000)
//}

func Test_GetGroup(t *testing.T) {
	t.Log("--------------------------")
	infos, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(ctx, []string{"3179997540"})
	if err != nil {
		t.Fatal(err)
	}
	for i, info := range infos {
		t.Logf("%d: %#v", i, info)
	}
	// time.Sleep(time.Second * 100000)
}
func Test_GetGroupApplicantsList(t *testing.T) {
	t.Log("--------------------------")
	infos, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsRecipient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for i, info := range infos {
		t.Logf("%d: %#v", i, info)
	}
	// time.Sleep(time.Second * 100000)
}
func Test_IsJoinGroup(t *testing.T) {
	t.Log("--------------------------")
	join, err := open_im_sdk.UserForSDK.Group().IsJoinGroup(ctx, "237875325046784")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("join:", join)
}

func Test_GetGroupMemberList(t *testing.T) {
	t.Log("--------------------------")
	m := map[int32]string{
		constant.GroupOwner:         "群主",
		constant.GroupAdmin:         "管理员",
		constant.GroupOrdinaryUsers: "成员",
	}

	members, err := open_im_sdk.UserForSDK.Group().GetGroupMemberList(ctx, "3889561099", 0, 0, 9999999)
	if err != nil {
		panic(err)
	}
	for i, member := range members {
		name := m[member.RoleLevel]
		t.Log(i, member.UserID, member.Nickname, name)
	}

	t.Log("--------------------------")
}

func Test_SyncAllGroupMember(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SyncAllGroupMember(ctx, "3889561099")
	if err != nil {
		t.Fatal(err)
	}
}
func Test_SetGroupInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupInfo(ctx, &groupPb.EditGroupProfileReq{
		GroupID: "3889561099",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetGroupMemberByPage(t *testing.T) {
	data, err := open_im_sdk.UserForSDK.Group().GetGroupMemberByPage(ctx, "237875325046784", 0, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range data.List {
		t.Log(fmt.Sprintf("%#v", d))
	}
}

func Test_GetGroupMembers(t *testing.T) {
	data, err := open_im_sdk.UserForSDK.Group().GetGroupMembers(ctx, "237875325046784", "177", 0, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data.Total)
	for _, d := range data.List {
		t.Log(fmt.Sprintf("%#v", d))
	}
}

func Test_GetGroupMemberById(t *testing.T) {
	data, err := open_im_sdk.UserForSDK.Group().GetGroupMemberById(ctx, "237875325046784", UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%#v", data))
}

func Test_GetUnProcessGroupRequestNum(t *testing.T) {
	num, err := open_im_sdk.UserForSDK.Group().GetUnProcessGroupRequestNum(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(num)
}
