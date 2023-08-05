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

// GetSpecifiedFriendsInfo 获取指定的好友信息
// par operationID  链路id
// par userIDList  用户ID集合
func GetSpecifiedFriendsInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, UserForSDK.Friend().GetSpecifiedFriendsInfo, userIDList)
}

// GetFriendList 获取好友信息
// par operationID  链路id
func GetFriendList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Friend().GetFriendList)
}

// GetPageFriendList 分页获取好友信息
// par operationID  链路id
// par no  页码
// par size   长度
func GetPageFriendList(callback open_im_sdk_callback.Base, operationID string, no int64, size int64) {
	call(callback, operationID, UserForSDK.Friend().GetFriendListPage, no, size)
}

// SearchFriends 搜索好友
// par operationID  链路id
func SearchFriends(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, UserForSDK.Friend().SearchFriends, searchParam)
}

// CheckFriend 校验是否好友
// par operationID  链路id
// par userIDList  用户ID
func CheckFriend(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, UserForSDK.Friend().CheckFriend, userIDList)
}

// AddFriend 添加好友
// par operationID  链路id
// par userIDReqMsg  用户ID
func AddFriend(callback open_im_sdk_callback.Base, operationID string, userIDReqMsg string) {
	call(callback, operationID, UserForSDK.Friend().AddFriend, userIDReqMsg)
}

// SetFriendRemark 设置好友备注
// par operationID  链路id
// par userIDReqMsg  用户ID和备注
func SetFriendRemark(callback open_im_sdk_callback.Base, operationID string, userIDRemark string) {
	call(callback, operationID, UserForSDK.Friend().SetFriendRemark, userIDRemark)
}

// DeleteFriend 删除好友
// par operationID  链路id
// par friendUserID  要删除的用户ID
func DeleteFriend(callback open_im_sdk_callback.Base, operationID string, friendUserID string) {
	call(callback, operationID, UserForSDK.Friend().DeleteFriend, friendUserID)
}

// GetFriendApplicationListAsRecipient 收到的好友请求
// par operationID  链路id
func GetFriendApplicationListAsRecipient(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Friend().GetFriendApplicationListAsRecipient)
}

// GetPageFriendApplicationListAsRecipient 分页获取收到的好友请求
// par operationID  链路id
// par no  页码
// par size   长度
func GetPageFriendApplicationListAsRecipient(callback open_im_sdk_callback.Base, operationID string, no int64, size int64) {
	call(callback, operationID, UserForSDK.Friend().GetPageFriendApplicationListAsRecipient, no, size)
}

// GetFriendApplicationListAsApplicant 发出的好友请求
// par operationID  链路id
func GetFriendApplicationListAsApplicant(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Friend().GetFriendApplicationListAsApplicant)
}

// GetPageFriendApplicationListAsApplicant 分页获取发出的好友请求
// par operationID  链路id
// par no  页码
// par size   长度
func GetPageFriendApplicationListAsApplicant(callback open_im_sdk_callback.Base, operationID string, no int64, size int64) {
	call(callback, operationID, UserForSDK.Friend().GetPageFriendApplicationListAsApplicant, no, size)
}

// AcceptFriendApplication 同意好友请求
// par operationID  链路id
// par userIDHandleMsg  处理好友请求信息
func AcceptFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, UserForSDK.Friend().AcceptFriendApplication, userIDHandleMsg)
}

// RefuseFriendApplication 拒绝好友请求
// par operationID  链路id
// par userIDHandleMsg  处理好友请求信息
func RefuseFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, UserForSDK.Friend().RefuseFriendApplication, userIDHandleMsg)
}

// AddBlack 加入和名单
// par operationID  链路id
// par blackUserID  加入黑名单的用户id
func AddBlack(callback open_im_sdk_callback.Base, operationID string, blackUserID string) {
	call(callback, operationID, UserForSDK.Friend().AddBlack, blackUserID)
}

// GetBlackList 黑名单列表
// par operationID  链路id
func GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Friend().GetBlackList)
}

// GetPageBlackList 分页获取黑明单列表
// par operationID  链路id
// par no  页码
// par size   长度
func GetPageBlackList(callback open_im_sdk_callback.Base, operationID string, no int64, size int64) {
	call(callback, operationID, UserForSDK.Friend().GetPageBlackList, no, size)
}

// RemoveBlack 将用户移出黑名单
// par operationID  链路id
// par removeUserID  移除黑明单的用户ID
func RemoveBlack(callback open_im_sdk_callback.Base, operationID string, removeUserID string) {
	call(callback, operationID, UserForSDK.Friend().RemoveBlack, removeUserID)
}

// GetUnprocessedNum 获取未处理的好友请求
// par operationID  链路id
// return count 未处理的角标数
func GetUnprocessedNum(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Friend().GetUnprocessedNum)
}
