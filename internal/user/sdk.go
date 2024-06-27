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

package user

import (
	"context"

	"github.com/brian-god/imcloud_sdk/pkg/constant"
	"github.com/brian-god/imcloud_sdk/pkg/db/model_struct"
	"github.com/brian-god/imcloud_sdk/pkg/server_api"
	usermodel "github.com/miliao_apis/api/common/model/user/v2"
	userPb "github.com/miliao_apis/api/im_cloud/user/v2"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	return u.GetUsersInfoFromSvr(ctx, userIDs)
}

func (u *User) GetSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.getSelfUserInfo(ctx)
}

// Deprecated: user SetSelfInfoEx instead
func (u *User) SetSelfInfo(ctx context.Context, userInfo *userPb.UpdateProfileReq) error {
	return u.updateSelfUserInfo(ctx, userInfo)
}
func (u *User) SetSelfInfoEx(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	return u.updateSelfUserInfoEx(ctx, userInfo)
}
func (u *User) SetGlobalRecvMessageOpt(ctx context.Context, opt int) error {
	if err := server_api.SetUsersOption(ctx, u.loginUserID, userPb.UserOption_globalRecvMsgOpt.String(), int32(opt)); err != nil {
		return err
	}
	u.SyncLoginUserInfo(ctx)
	return nil
}

func (u *User) UpdateMsgSenderInfo(ctx context.Context, nickname, faceURL string) (err error) {
	if nickname != "" {
		if err = u.DataBase.UpdateMsgSenderNickname(ctx, u.loginUserID, nickname, constant.SingleChatType); err != nil {
			return err
		}
	}
	if faceURL != "" {
		if err = u.DataBase.UpdateMsgSenderFaceURL(ctx, u.loginUserID, faceURL, constant.SingleChatType); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) SubscribeUsersStatus(ctx context.Context, userIDs []string) ([]*usermodel.OnlineStatus, error) {
	//userStatus, err := u.subscribeUsersStatus(ctx, userIDs)
	//if err != nil {
	//	return nil, err
	//}
	//u.OnlineStatusCache.DeleteAll()
	//u.OnlineStatusCache.StoreAll(func(value *userPb.OnlineStatus) string {
	//	return value.UserID
	//}, userStatus)
	//return userStatus, nil
	return nil, nil
}

func (u *User) UnsubscribeUsersStatus(ctx context.Context, userIDs []string) error {
	u.OnlineStatusCache.DeleteAll()
	return u.unsubscribeUsersStatus(ctx, userIDs)
}

func (u *User) GetSubscribeUsersStatus(ctx context.Context) ([]*usermodel.OnlineStatus, error) {
	//return u.getSubscribeUsersStatus(ctx)
	return nil, nil
}

func (u *User) GetUserStatus(ctx context.Context, userID string) (*usermodel.OnlineStatus, error) {
	return u.getUserStatus(ctx, userID)
}
