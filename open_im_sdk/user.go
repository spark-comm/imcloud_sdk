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

package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

// GetUsersInfo 获取用户信息
func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, UserForSDK.Full().GetUsersInfo, userIDs)
}

// GetUsersInfo obtains the information about multiple users.
func GetUsersInfoFromSrv(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, UserForSDK.User().GetUsersInfo, userIDs)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, UserForSDK.User().SetSelfInfo, userInfo)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.User().GetSelfUserInfo)
}

// UpdateMsgSenderInfo updates the message sender's nickname and face URL.
func UpdateMsgSenderInfo(callback open_im_sdk_callback.Base, operationID string, nickname, faceURL string) {
	call(callback, operationID, UserForSDK.User().UpdateMsgSenderInfo, nickname, faceURL)
}

// SearchUser by search value and search type
// par operationID  链路id
// par searchValue  搜索的值
// par searchType   类型 1:手机号，2:id,3扫码
func SearchUser(callback open_im_sdk_callback.Base, operationID string, searchValue string, searchType int) {
	call(callback, operationID, UserForSDK.User().SearchUserInfo, searchValue, searchType)
}

// GetLoginUserStatus 获取用户状态
// par operationID  链路id
// par userID  用户id
func GetLoginUserStatus(callback open_im_sdk_callback.Base, operationID string, userID string) {
	call(callback, operationID, UserForSDK.User().GetUserLoginStatus, userID)
}

// SetUsersOption 设置用户配置项
// par operationID  链路id
// @par option string  配置项
// @par value  number  值
func SetUsersOption(callback open_im_sdk_callback.Base, operationID, option string, value int32) {
	call(callback, operationID, UserForSDK.User().SetUsersOption, option, value)
}

// SyncUsersWalletOption 同步用户钱包是否开通状态
func SyncUsersWalletOption(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.User().SyncUserOperation)
}
