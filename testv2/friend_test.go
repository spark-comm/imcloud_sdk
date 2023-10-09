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
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
	"testing"
	"time"

	friend "github.com/imCloud/api/friend/v1"
)

func Test_GetSpecifiedFriendsInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetSpecifiedFriendsInfo(ctx, []string{"69763666668425216"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetDesignatedFriendsInfo success", ctx.Value("operationID"))
	for _, userInfo := range info {
		t.Log(userInfo)
	}
}

func Test_AddFriend(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().AddFriend(ctx, &friend.AddFriendRequest{
		ToUserID:  "55122365784264704",
		ReqMsg:    "test add",
		RemarkMsg: "天加",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AddFriend success", ctx.Value("operationID"))
}

//funcation Test_GetRecvFriendApplicationList(t *testing.T) {
//	infos, err := open_im_sdk.UserForSDK.Friend().GetRecvFriendApplicationList(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, info := range infos {
//		t.Logf("%#v", info)
//	}
//}
//
//funcation Test_GetSendFriendApplicationList(t *testing.T) {
//	infos, err := open_im_sdk.UserForSDK.Friend().GetSendFriendApplicationList(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, info := range infos {
//		t.Logf("%#v", info)
//	}
//}

func Test_AcceptFriendApplication(t *testing.T) {
	req := &sdk_params_callback.ProcessFriendApplicationParams{ToUserID: "48676976868724736", HandleMsg: "test accept"}
	err := open_im_sdk.UserForSDK.Friend().AcceptFriendApplication(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptFriendApplication success", ctx.Value("operationID"))
	time.Sleep(time.Second * 30)
}

func Test_RefuseFriendApplication(t *testing.T) {
	req := &sdk_params_callback.ProcessFriendApplicationParams{ToUserID: "48676976868724736", HandleMsg: "test refuse"}
	err := open_im_sdk.UserForSDK.Friend().RefuseFriendApplication(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("RefuseFriendApplication success", ctx.Value("operationID"))
	time.Sleep(time.Second * 30)
}

func Test_CheckFriend(t *testing.T) {
	err2 := open_im_sdk.UserForSDK.Friend().DelLocalFriend(ctx, "55122332112392192")
	if err2 != nil {
		t.Fatal(err2)
	}
	res, err := open_im_sdk.UserForSDK.Friend().CheckFriend(ctx, []string{"55122332112392192"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CheckFriend success", ctx.Value("operationID"))
	for _, re := range res {
		t.Log(re)
	}
}

func Test_DeleteFriend(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().DeleteFriend(ctx, "863454357")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DeleteFriend success", ctx.Value("operationID"))
}

func Test_GetFriendList(t *testing.T) {
	infos, err := open_im_sdk.UserForSDK.Friend().GetFriendList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetFriendList success", ctx.Value("operationID"))
	for _, info := range infos {
		t.Logf("PublicInfo: %#v", info)
	}
}

// Test_GetPageFriendList 分页获取好友数据
func Test_GetPageFriendList(t *testing.T) {
	infos, err := open_im_sdk.UserForSDK.Friend().GetFriendListPage(ctx, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetFriendList success", ctx.Value("operationID"))
	for _, info := range infos {
		t.Logf("PublicInfo: %#v", info)
	}
}

func Test_SearchFriends(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().SearchFriends(ctx, &sdk_params_callback.SearchFriendsParam{KeywordList: []string{"我"}, IsSearchUserID: true, IsSearchRemark: true, IsSearchNickname: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SearchFriends success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}

func Test_SetFriendRemark(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SetFriendRemark(ctx, &sdk_params_callback.SetFriendRemarkParams{ToUserID: "863454357", Remark: "testRemark"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetFriendRemark success", ctx.Value("operationID"))
}

func Test_AddBlack(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().AddBlack(ctx, "863454357")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AddBlack success", ctx.Value("operationID"))
}

func Test_RemoveBlack(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().RemoveBlack(ctx, "48672487050842112")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("RemoveBlack success", ctx.Value("operationID"))
}

func Test_GetBlackList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetBlackList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetBlackList success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}

func Test_GetFriendListPage(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetFriendListPage(ctx, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetFriendListPage success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}

// Test_GetPageFriendApplicationListAsRecipient 收到申请
func Test_GetPageFriendApplicationListAsRecipient(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetPageFriendApplicationListAsRecipient(ctx, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetPageFriendApplicationListAsRecipient success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}

func Test_GetPageFriendApplicationListAsApplicant(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetPageFriendApplicationListAsApplicant(ctx, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetPageFriendApplicationListAsApplicant success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}
func Test_GetPageBlackList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Friend().GetPageBlackList(ctx, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetPageBlackList success", ctx.Value("operationID"))
	for _, item := range info {
		t.Log(*item)
	}
}

func Test_GetUnprocessedNum(t *testing.T) {
	count, err := open_im_sdk.UserForSDK.Friend().GetUnprocessedNum(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetPageBlackList success", ctx.Value("operationID"))

	t.Log("unprocessed num ", count)
}

func Test_SetFriendChatBackground(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SetBackgroundUrl(ctx, "55224333915656192", "背景1是群1")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	t.Log("SetFriendChatBackground success")
}

func Test_SetFriendDestroyMsgStatus(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SetFriendDestroyMsgStatus(ctx, "55122367646535680", 0)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	t.Log("SetFriendDestroyMsgStatus success")
}
func Test_GetFriendBaseInfoSvr(t *testing.T) {
	friends, err := open_im_sdk.UserForSDK.Friend().GetFriendBaseInfoSvr(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range friends {
		t.Log("sync data", fmt.Sprintf("%s", utils.StructToJsonString(v)))
	}
	t.Log("GetFriendBaseInfoSvr success")
}
