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
	imUserPb "github.com/imCloud/api/user/v1"
	"open_im_sdk/open_im_sdk"
	"testing"
	"time"
)

func Test_GetSelfUserInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(userInfo)
}

func Test_GetUsersInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.Full().GetUsersInfo(ctx, []string{"45778745637736448"})
	if err != nil {
		t.Error(err)
	}
	if userInfo[0].BlackInfo != nil {
		t.Log(userInfo[0].BlackInfo)
	}
	if userInfo[0].FriendInfo != nil {
		t.Log(userInfo[0].FriendInfo)
	}
	if userInfo[0].PublicInfo != nil {
		t.Log(userInfo[0].PublicInfo)
	}
}

func Test_SetSelfInfo(t *testing.T) {
	newNickName := "test1"
	newFaceURL := "http://test.com"
	err := open_im_sdk.UserForSDK.User().SetSelfInfo(ctx, &imUserPb.UpdateProfileReq{
		Nickname: newNickName,
		FaceUrl:  newFaceURL,
	})
	//newFaceURL := "http://test.com"

	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	if userInfo.UserID != UserID && userInfo.Nickname != newNickName && userInfo.FaceURL != newFaceURL {
		t.Error("user id not match")
	}
	t.Log(userInfo)
	time.Sleep(time.Second * 10)
}

func Test_UpdateMsgSenderInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.User().UpdateMsgSenderInfo(ctx, "test", "http://test.com")
	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(userInfo)
}

type SearchCallback struct {
}

func (m *SearchCallback) OnError(errCode int32, errMsg string) {
	fmt.Println("错误")
}
func (m *SearchCallback) OnSuccess(data string) {
	fmt.Println("成功返回", data)
}
func Test_SearchUserInfo(t *testing.T) {
	callback := SearchCallback{}
	open_im_sdk.SearchUser(&callback, "jhjguyuyguyguyguy", "au9699", 2)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(userInfo)
}
