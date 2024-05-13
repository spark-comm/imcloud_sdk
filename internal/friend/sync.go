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
	"fmt"
	commonPb "github.com/imCloud/api/common"
	friendPb "github.com/imCloud/api/friend/v1"
	"github.com/imCloud/im/pkg/common/log"
	"gorm.io/gorm"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"time"
)

// SyncSelfFriendApplication 自己发送的好友请求
func (f *Friend) SyncSelfFriendApplication(ctx context.Context) error {
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.GetPaginationFriendsApplyFromResp) []*friendPb.FriendRequests {
		return resp.FriendRequests
	}
	resp := &friendPb.GetPaginationFriendsApplyFromResp{}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, resp, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	err = f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
	if err != nil {
		return err
	}
	return nil
}

// SyncFriendApplication 同步自己收到的好友请求
func (f *Friend) SyncFriendApplication(ctx context.Context) error {
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.GetPaginationFriendsApplyReceiveResp) []*friendPb.FriendRequests {
		return resp.FriendRequests
	}
	resp := &friendPb.GetPaginationFriendsApplyReceiveResp{}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendReceiveApplicationListRouter, req, resp, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestRecvSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

// SyncUntreatedFriendReceiveFriendApplication 同步未处理的好友请求
func (f *Friend) SyncUntreatedFriendReceiveFriendApplication(ctx context.Context) error {
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.GetUntreatedFriendsApplyReceiveReply) []*friendPb.FriendRequests {
		return resp.FriendRequests
	}
	resp := &friendPb.GetUntreatedFriendsApplyReceiveReply{}
	requests, err := util.GetPageAll(ctx, constant.GetUntreatedFriendsApplyReceive, req, resp, fn)
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
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.ListFriendReply) []*friendPb.FriendInfo {
		return resp.List
	}
	resp := &friendPb.ListFriendReply{}
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, resp, fn)
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

// SyncFirstFriendList 同步第一页好友列表
func (f *Friend) SyncFirstFriendList(ctx context.Context) error {
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{
		ShowNumber: 10,
	}}
	fn := func(resp *friendPb.ListFriendReply) []*friendPb.FriendInfo {
		return resp.List
	}
	resp := &friendPb.ListFriendReply{}
	friends, err := util.GetFirstPage(ctx, constant.GetFriendListRouter, req, resp, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	if err = f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil); err != nil {
		log.ZDebug(ctx, "sync first page friend error", err)
	}
	//加入延迟队列做同步
	f.syncFriendQueue.Push(1, time.Second*20)
	return err
}

// SyncBlackList 同步黑名单信息
func (f *Friend) SyncBlackList(ctx context.Context) error {
	req := &friendPb.GetPaginationFriendsInfo{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.BlackListResponse) []*friendPb.BlackList { return resp.Data }
	resp := &friendPb.BlackListResponse{}
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, resp, fn)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := f.db.GetBlackListAllDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return f.blockSyncer.Sync(ctx, util.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}

// syncFriendApplicationById 根据id同步好友请求
func (f *Friend) syncFriendApplicationById(ctx context.Context, fromUserID, toUserID string) error {
	//req := &friendPb.GetFriendRequestByApplicantReq{FromUserID: fromUserID, ToUserID: toUserID}
	//res, err := util.CallApi[friendPb.GetFriendRequestByApplicantReps](ctx, constant.GetFriendRequestByApplicantRouter, req)
	res := &friendPb.GetFriendRequestByApplicantReps{}
	err := util.CallPostApi[*friendPb.GetFriendRequestByApplicantReq, *friendPb.GetFriendRequestByApplicantReps](
		ctx, constant.GetFriendRequestByApplicantRouter,
		&friendPb.GetFriendRequestByApplicantReq{FromUserID: fromUserID, ToUserID: toUserID},
		res,
	)
	if err != nil {
		return err
	}
	if res.FriendRequest == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nil")
		return nil
	}
	localData, err := f.db.GetFriendApplicationByBothID(ctx, fromUserID, toUserID)
	localList := make([]*model_struct.LocalFriendRequest, 0)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else {
		localList = append(localList, localData)
	}
	if err = f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest,
		[]*friendPb.FriendRequests{res.FriendRequest}), localList, nil); err != nil {
		return err
	}
	//log.ZInfo(ctx, "开始根据id同步好友请求")
	//if err = f.requestRecvSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest,
	//	[]*friendPb.FriendRequests{res.FriendRequest}), localList, nil); err != nil {
	//	log.ZInfo(ctx, fmt.Sprintf("开始根据id同步好友请求失败,err:%+v", err))
	//	return err
	//}
	return nil
}

// syncFriendById 根据id同步好友
func (f *Friend) syncFriendById(ctx context.Context, friendId string) error {
	//req := &friendPb.ListFriendByIdsReq{UserID: fromUserID, FriendIds: []string{friendId}}
	//res, err := util.CallApi[friendPb.ListFriendByIdsReply](ctx, constant.GetFriendByAppIdsRouter, req)
	svr, err1 := f.GetFriendByIdsSvr(ctx, []string{friendId})
	if err1 != nil {
		return err1
	}
	localData, err := f.db.GetFriendInfoList(ctx, []string{friendId}, false)
	//localList := make([]*model_struct.LocalFriend, 0)
	if err != nil {
		log.ZError(ctx, "syncFriendById->GetFriendInfoList", err)
		return err
	}
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, svr), localData, nil, true)
}

// syncFriendByInfo 同步指定好友信息
func (f *Friend) syncFriendByInfo(ctx context.Context, data []*model_struct.LocalFriend) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	friendIds := make([]string, len(data))
	for i, v := range data {
		friendIds[i] = v.FriendUserID
	}
	localData, err := f.db.GetFriendInfoList(ctx, friendIds, false)
	localList := make([]*model_struct.LocalFriend, 0)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else {
		localList = append(localList, localData...)
	}
	return f.friendSyncer.Sync(ctx, data, localData, nil, true)
}

// SyncDelFriend 同步删除好友列表
func (f *Friend) SyncDelFriend(ctx context.Context, friendId string) error {
	localData, err := f.db.GetFriendInfoNotPeersList(ctx, []string{friendId})
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from local", localData)
	return f.friendSyncer.Delete(ctx, localData, nil)
}

// SyncQuantityFriendList 全量同步好友，部分字段
func (f *Friend) SyncQuantityFriendList(ctx context.Context) error {
	//获取远端数据
	respList, err := f.GetFriendBaseInfoSvr(ctx)
	if err != nil {
		return err
	}
	fmt.Println("获取到的好友信息", len(respList))
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	friends := make([]*friendPb.FriendInfo, 0)
	for _, info := range respList {
		friends = append(friends, &friendPb.FriendInfo{
			OwnerUserID:   f.loginUserID,
			FriendUserID:  info.FriendID,
			Nickname:      info.NickName,
			FaceURL:       info.FaceURL,
			Code:          info.Code,
			Remark:        info.Remark,
			BackgroundUrl: info.BackgroundURL,
			IsComplete:    IsNotComplete, //未同步完成标记
			UpdateAt:      info.UpdateAt,
		})
	}
	//同步本地没有的数据
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	if err = f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil); err != nil {
		log.ZDebug(ctx, "sync first page friend error", err)
	}
	return err
}

// SyncFriend 同步好友信息
func (f *Friend) SyncFriend(ctx context.Context) error {
	udata, err := f.db.GetFriendUpdateTime(ctx)
	if err != nil || len(udata) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(udata) == 0 {
			return f.SyncQuantityFriendList(ctx)
		}
		return err
	} else {
		return f.SyncFriendByTime(ctx, udata)
	}
}

func (f *Friend) SyncFriendByTime(ctx context.Context, udata map[string]int64) error {
	addList, updateList, delIds, err := f.SyncFriendInfoByTime(ctx, udata)
	if err != nil {
		log.ZError(ctx, "sync first page friend error", err)
		return err
	}
	//定义一条改变的数据
	friend := &model_struct.LocalFriend{}
	//处理新增
	if len(addList) > 0 {
		for i, v := range addList {
			localFriend := ServerBaseFriendToLocalFriend(v)
			localFriend.OwnerUserID = f.loginUserID
			err = f.db.InsertFriend(ctx, localFriend)
			if err != nil {
				log.ZError(ctx, "insert friend error", err)
			}
			if len(addList) == i+1 {
				friend = localFriend
			}
		}
	}
	//处理修复爱
	if len(updateList) > 0 {
		for i, v := range updateList {
			localFriend := ServerBaseFriendToLocalFriend(v)
			localFriend.OwnerUserID = f.loginUserID
			err = f.db.UpdateFriend(ctx, localFriend)
			if err != nil {
				log.ZError(ctx, "insert friend error", err)
			}
			if len(addList) == i+1 {
				friend = localFriend
			}
		}
	}
	//处理删除
	if len(delIds) > 0 {
		friend.FriendUserID = delIds[0]
		err = f.db.DeleteFriendDB(ctx, delIds...)
		if err != nil {
			log.ZError(ctx, "insert friend error", err)
		}
		for _, v := range delIds {
			f.DelFriendConversation(ctx, v)
		}
	}
	if friend.FriendUserID != "" {
		//最后触发用户信息变动
		f.friendListener.OnFriendInfoChanged(*friend)
	}
	return nil
}
