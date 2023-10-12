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
	"encoding/json"
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
	"testing"
	"time"

	groupv1 "github.com/imCloud/api/group/v1"
)

type GroupCallback struct {
}

func (g *GroupCallback) OnError(errCode int32, errMsg string) {
	log.Info("", "!!!!!!!OnError ")
}

func (g *GroupCallback) OnSuccess(data string) {
	log.Info("", "!!!!!!!OnSuccess ")
}
func Test_CreateGroupV2(t *testing.T) {
	req := &groupv1.CrateGroupReq{
		MemberList:       []string{"55476661776289792", "55474682626838528", "55474678449311744"},
		GroupName:        "月华馥",
		GroupType:        2,
		Notification:     "公告：这是一个超级大群",
		Introduction:     "介绍一下：大群",
		FaceURL:          "https://dfsjk/djfhsd/5d1f5562d/154452.jpg",
		NeedVerification: 1,
	}
	info, err := open_im_sdk.UserForSDK.Group().CreateGroup(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("group info: %s", info.String())
	sessionType := int32(3)
	conversation, err := open_im_sdk.UserForSDK.Conversation().GetOneConversation(
		ctx,
		sessionType,
		"118482918182912",
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sduky51515154==========%+v==============", conversation)
	dfsd, err := open_im_sdk.UserForSDK.Conversation().Get123dfsd(ctx, conversation.ConversationID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("conversationInfo-------------------%+v---------------", dfsd)
}

func Test_JoinGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().JoinGroup(ctx,
		"114711123202048",
		"进群收钱呀",
		1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JoinGroup success")
}

func Test_QuitGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().QuitGroup(ctx, "114711123202048")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("QuitGroup success")
}

func Test_DismissGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().DismissGroup(ctx, "139270752833536")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DismissGroup success")
}

func Test_ChangeGroupMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(ctx,
		"161276114243584", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelMuteGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(
		ctx,
		"850093334859776",
		false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_ChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(
		ctx,
		"120143539605504",
		"1463426512456",
		10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(ctx,
		"120143539605504",
		"1463426512456",
		0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CancelChangeGroupMemberMute success", ctx.Value("operationID"))
}

func Test_SetGroupMemberRoleLevel(t *testing.T) {
	// 1:普通成员 2:群主 3:管理员
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberRoleLevel(
		ctx,
		"120143539605504",
		"1463426512456",
		3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberRoleLevel success", ctx.Value("operationID"))
}

func Test_SetGroupMemberNickname(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberNickname(
		ctx,
		"120143539605504",
		"1463426512456",
		"头皮发麻-123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberNickname success", ctx.Value("operationID"))
}

func Test_SetGroupMemberInfo(t *testing.T) {
	// 1:普通成员 2:群主 3:管理员
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberInfo(ctx, &groupv1.SetGroupMemberInfoReq{
		GroupID:  "120143539605504",
		UserID:   "1463426512456",
		FaceURL:  "https://doc.rentsoft.cn/images/logo.png",
		Nickname: "熔火之心",
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

func Test_GetSpecifiedGroupsInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(
		ctx, []string{"136998765465600"})
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

func Test_GetPageGroupApplicationListAsRecipient(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetPageGroupApplicationListAsRecipient(ctx, "181523118559232", 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetSendGroupApplicationList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
	time.Sleep(time.Minute * 2)
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
	err := open_im_sdk.UserForSDK.Group().AcceptGroupApplication(ctx,
		"391667123228672",
		"55122326609465344",
		"test accept")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_RefuseGroupApplication(t *testing.T) {
	t.Log("operationID:", ctx.Value("operationID"))
	err := open_im_sdk.UserForSDK.Group().RefuseGroupApplication(
		ctx,
		"50913016287232",
		"1463426528082",
		"test refuse")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success")
}

func Test_HandlerGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().HandlerGroupApplication(ctx, &groupv1.ApplicationResponseReq{
		GroupID:      "78877338636288",
		FromUserID:   "1463426527031",
		HandledMsg:   "FDSFSFSF",
		HandleResult: constant.GroupResponseAgree,
		UserID:       "1463426574231",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_SearchGroupMembers(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().SearchGroupMembers(ctx, &sdk_params_callback.SearchGroupMembersParam{
		GroupID:                "57308482637824",
		KeywordList:            []string{""},
		IsSearchUserID:         false,
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
	err := open_im_sdk.UserForSDK.Group().KickGroupMember(
		ctx,
		"197403915325440",
		"我要踢人",
		[]string{"55224130529660928"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

func Test_TransferGroupOwner(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().TransferGroupOwner(
		ctx,
		"78877338636288",
		"1463426515762")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("TransferGroupOwner success", ctx.Value("operationID"))
}

func Test_InviteUserToGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().InviteUserToGroup(
		ctx,
		"106234346737664",
		"测试邀请人进群",
		[]string{"55122324243877888", "55122324344541184", "55122324461981696", "55122324579422208"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

func Test_SyncGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SyncGroups(ctx,
		"171491979169792")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 100000)
}

func Test_GetGroup(t *testing.T) {
	t.Log("--------------------------")
	infos, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(ctx,
		[]string{"105942007943168"})
	if err != nil {
		t.Fatal(err)
	}
	for i, info := range infos {
		t.Logf("%d: %#v", i, info)
	}
	//time.Sleep(time.Second * 100000)
}

func Test_IsJoinGroup(t *testing.T) {
	t.Log("--------------------------")
	join, err := open_im_sdk.UserForSDK.Group().IsJoinGroup(ctx, "1875806101")
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

	members, err := open_im_sdk.UserForSDK.Group().GetGroupMemberList(
		ctx,
		"153174048509952", 0, 1, 20)
	if err != nil {
		panic(err)
	}
	for i, member := range members {
		name := m[member.RoleLevel]
		t.Log("sdfs1f153165516165651456", i, member.UserID, member.Nickname, name)
	}

	t.Log("--------------------------")
}
func Test_GetGroupMemberList1(t *testing.T) {
	open_im_sdk.GetGroupMemberList(&GroupCallback{}, utils.OperationIDGenerator(), "503233625722880", 0, 0, 9999999)
	time.Sleep(time.Second * 4)
	t.Log("--------------------------")
}

func Test_CreateGroup(t *testing.T) {
	req := groupv1.CrateGroupReq{
		MemberList:       []string{"1463426311015", "1463426512456", "1463426515762", "1463426574029"},
		GroupName:        "123456789",
		GroupType:        2,
		Notification:     "公告：这是一个荣誉",
		Introduction:     "洗脑群",
		FaceURL:          "https://dfsjk/djfhsd/5d1f5562d/154452.jpg",
		NeedVerification: 1,
	}
	marshal, _ := json.Marshal(&req)
	open_im_sdk.CreateGroup(&GroupCallback{}, utils.OperationIDGenerator(), string(marshal))
	time.Sleep(time.Second * 4)
	t.Log("--------------------------")
}

func Test_SetGroupInfo(t *testing.T) {
	s := groupv1.EditGroupProfileRequest{
		GroupID:      "145196112",
		Notification: "-48559",
	}
	bytes, _ := json.Marshal(&s)
	open_im_sdk.SetGroupInfo(&GroupCallback{}, utils.OperationIDGenerator(), string(bytes))
	time.Sleep(time.Second * 4)
	t.Log("--------------------------")
}

func Test_SetGroupSwitchInfo(t *testing.T) {
	open_im_sdk.SetGroupSwitchInfo(
		&GroupCallback{},
		utils.OperationIDGenerator(), "1189039579336704", "need_verification", 0)
	time.Sleep(time.Second * 4)
	t.Log("--------------------------")
}
func Test_KickGroupUserList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().KickGroupMemberList(
		ctx, &sdk_params_callback.GetKickGroupListReq{
			GroupID:  "136998765465600",
			IsManger: false,
			Name:     "",
			PageSize: 10,
			PageNum:  1,
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetGroupsInfo: %d\n", len(info.KickGroupList))
	for _, localGroup := range info.KickGroupList {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetNotInGroupFriendInfoList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetNotInGroupFriendInfoList(
		ctx, &sdk_params_callback.SearchNotInGroupUserReq{
			GroupID:  "110908089044992",
			PageSize: 10,
			PageNum:  1,
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetGroupsInfo: %d\n", len(info.Friends))
	for _, localGroup := range info.Friends {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetUserOwnerJoinRequestNum(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetUserOwnerJoinRequestNum(
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, localGroup := range info.Data {
		t.Logf("%#v", localGroup)
	}
}

func Test_SearchGroupInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().SearchGroupInfo(
		ctx,
		"月华",
		20,
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, localGroup := range info.Data {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetAppointGroupRequestInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetAppointGroupRequestInfo(
		ctx,
		"45362416979968", 1, 20,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintf("%+v", info)
}

func Test_SetGroupChatBackground(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetBackgroundUrl(ctx, "105942007943168", "背景")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	t.Log("SetFriendChatBackground success")
}

func Test_GetUserMemberInfoInGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupChatBackground(ctx)
	if err != nil {
		return
	}
	return
}

func Test_SearchGroupByCode(t *testing.T) {
	groupInfo, err := open_im_sdk.UserForSDK.Group().SearchGroupByCode(ctx, "57p34d")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetFriendChatBackground success", groupInfo)
}
