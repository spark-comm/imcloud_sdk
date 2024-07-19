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
	"fmt"

	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	"github.com/spark-comm/imcloud_sdk/pkg/sdkerrs"
	"github.com/spark-comm/imcloud_sdk/pkg/server_api"

	"github.com/OpenIMSDK/tools/log"
)

func (f *Friend) SyncBothFriendRequest(ctx context.Context, fromUserID, toUserID string) error {
	friendRequests, err := server_api.BothFriendRequest(ctx, fromUserID, toUserID)
	localData, err := f.db.GetBothFriendReq(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}
	fmt.Println("localData", friendRequests)
	data := []*model_struct.LocalFriendRequest{friendRequests}
	if toUserID == f.loginUserID {
		return f.requestRecvSyncer.Sync(ctx, data, localData, nil)
	} else if fromUserID == f.loginUserID {
		return f.requestSendSyncer.Sync(ctx, data, localData, nil)
	}
	return nil
}

// send
func (f *Friend) SyncAllSelfFriendApplication(ctx context.Context) error {
	requests, err := server_api.GetSendFriendApplication(ctx, f.loginUserID)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestSendSyncer.Sync(ctx, requests, localData, nil)
}

// recv
func (f *Friend) SyncAllFriendApplication(ctx context.Context) error {
	requests, err := server_api.GetReceiveFriendApplication(ctx, f.loginUserID)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestRecvSyncer.Sync(ctx, requests, localData, nil)
}

func (f *Friend) SyncAllFriendList(ctx context.Context) error {
	friends, err := server_api.GetAllFriendList(ctx, f.loginUserID)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	return f.friendSyncer.Sync(ctx, friends, localData, nil)
}

func (f *Friend) deleteFriend(ctx context.Context, friendUserID string) error {
	friends, err := f.db.GetFriendInfoList(ctx, []string{friendUserID}, false)
	if err != nil {
		return err
	}
	if len(friends) == 0 {
		return sdkerrs.ErrUserIDNotFound.Wrap("friendUserID not found")
	}
	// todo 删除好友后，删除好友的聊天记录
	if err := f.db.DeleteFriendDB(ctx, friendUserID); err != nil {
		return err
	}
	f.friendListener.OnFriendDeleted(*friends[0])
	return nil
}

func (f *Friend) SyncFriends(ctx context.Context, friendIDs []string) error {
	friends, err := server_api.GetFriendByIds(ctx, f.loginUserID, friendIDs)
	if err != nil {
		return err
	}
	localData, err := f.db.GetFriendInfoList(ctx, friendIDs, false)
	if err != nil {
		return err
	}
	if friends == nil {
		friends = make([]*model_struct.LocalFriend, 0)
	}
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	return f.friendSyncer.Sync(ctx, friends, localData, nil)
}

func (f *Friend) SyncAllBlackList(ctx context.Context) error {
	serverData, err := server_api.GetAllBlackList(ctx, f.loginUserID)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return f.blockSyncer.Sync(ctx, serverData, localData, nil)
}
