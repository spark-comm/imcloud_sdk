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
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"

	imUserPb "github.com/imCloud/api/user/v1"
)

// GetUsersInfo 获取用户信息直接查服务器
func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	return u.GetUsersInfoFromSvr(ctx, userIDs)
}

// GetSelfUserInfo 获取登录用户信息
func (u *User) GetSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.getSelfUserInfo(ctx)
}

// SetSelfInfo 修改自己的信息
func (u *User) SetSelfInfo(ctx context.Context, userInfo *imUserPb.UpdateProfileReq) error {
	return u.updateSelfUserInfo(ctx, userInfo)
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

// SearchUserInfo 1-手机 2-用户ID 3-扫码 5-邮箱
// todo 新添加方法
func (u *User) SearchUserInfo(ctx context.Context, searchValue string, searchType int) (*model_struct.LocalUser, error) {
	return u.searchUser(ctx, searchValue, searchType)
}

func (u *User) GetUserLoginStatus(ctx context.Context, userIDs string) (*imUserPb.GetUserLoginStatusReps, error) {
	return u.getUserLoginStatus(ctx, userIDs)
}

// SetUsersOption 设置用户配置项
// @par option string  配置项
// @par value  number  值
func (u *User) SetUsersOption(ctx context.Context, option string, value int32) error {
	return u.setUsersOption(ctx, option, value)
}

func (u *User) SyncUserOperation(ctx context.Context) error {
	return u.syncUserOperation(ctx)
}

func (u *User) ScreenUserProfile(ctx context.Context, keyWord string) ([]*imUserPb.ScreenUserInfo, error) {
	return u.screenUserProfile(ctx, keyWord)
}
