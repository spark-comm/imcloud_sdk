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

package friend

import (
	"context"
	"errors"
	friendPb "github.com/imCloud/api/friend/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/friend"
	"github.com/imCloud/im/pkg/proto/sdkws"
	"gorm.io/gorm"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
)

// SyncSelfFriendApplication 自己发送的好友请求
func (f *Friend) SyncSelfFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friendPb.GetPaginationFriendsApplyFromResp) []*friendPb.FriendRequests {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

// SyncFriendApplication 同步自己收到的好友请求
func (f *Friend) SyncFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friendPb.GetPaginationFriendsApplyReceiveResp) []*friendPb.FriendRequests {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendReceiveApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestRecvSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

// SyncFriendList 同步好友列表
func (f *Friend) SyncFriendList(ctx context.Context) error {
	req := &friend.GetPaginationFriendsReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friendPb.ListFriendReply) []*friendPb.FriendInfo {
		return resp.List
	}
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil)
}

// SyncBlackList 同步黑名单信息
func (f *Friend) SyncBlackList(ctx context.Context) error {
	req := &friend.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friendPb.BlackListResponse) []*friendPb.BlackList { return resp.Data }
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return f.blockSyncer.Sync(ctx, util.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}

// syncFriendApplicationById 根据id同步好友请求
func (f *Friend) syncFriendApplicationById(ctx context.Context, fromUserID, toUserID string) error {
	req := &friendPb.GetFriendRequestByApplicantReq{FromUserID: fromUserID, ToUserID: toUserID}
	res, err := util.CallApi[friendPb.GetFriendRequestByApplicantReps](ctx, constant.GetFriendRequestByApplicantRouter, req)
	if err != nil {
		return err
	}
	if res.FriendRequest == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nill")
		return nil
	}
	localData, err := f.db.GetFriendApplicationByBothID(ctx, fromUserID, toUserID)
	localList := make([]*model_struct.LocalFriendRequest, 0)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else {
		localList = append(localList, localData)
	}
	return f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, []*friendPb.FriendRequests{res.FriendRequest}), localList, nil)
}

// syncFriendById 根据id同步好友
func (f *Friend) syncFriendById(ctx context.Context, fromUserID, friendId string) error {
	req := &friendPb.ListFriendByIdsReq{UserID: fromUserID, FriendIds: []string{friendId}}
	res, err := util.CallApi[friendPb.ListFriendByIdsReply](ctx, constant.GetFriendByAppIdsRouter, req)
	if err != nil {
		return err
	}
	if res.FriendsInfo == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nill")
		return nil
	}
	localData, err := f.db.GetFriendInfoList(ctx, []string{friendId})
	localList := make([]*model_struct.LocalFriend, 0)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else {
		localList = append(localList, localData...)
	}
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, res.FriendsInfo), localData, nil, true)
}

// syncDelFriend 同步删除好友列表
func (f *Friend) syncDelFriend(ctx context.Context, friendId string) error {
	localData, err := f.db.GetFriendInfoList(ctx, []string{friendId})
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from local", localData)
	return f.friendSyncer.Delete(ctx, localData, nil)
}
