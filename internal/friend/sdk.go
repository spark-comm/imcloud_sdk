// Copyright 2021 OpenIM Corporation
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

	"github.com/OpenIMSDK/tools/log"
	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	sdk "github.com/spark-comm/imcloud_sdk/pkg/sdk_params_callback"
	"github.com/spark-comm/imcloud_sdk/pkg/sdkerrs"
	"github.com/spark-comm/imcloud_sdk/pkg/server_api"
	"github.com/spark-comm/imcloud_sdk/pkg/server_api_params"
	friendPb "github.com/spark-comm/spark-api/api/im_cloud/friend/v2"
)

func (f *Friend) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList, false)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "localFriendList", localFriendList)
	blackList, err := f.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "blackList", blackList)
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range blackList {
		m[black.BlackUserID] = blackList[i]
	}
	res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return res, nil
}

func (f *Friend) AddFriend(ctx context.Context, req *friendPb.AddFriendReq) error {
	if req.FromUserID == "" {
		req.FromUserID = f.loginUserID
	}
	if err := server_api.AddFriend(ctx, req); err != nil {
		return err
	}
	return f.SyncAllFriendApplication(ctx)
}

func (f *Friend) GetFriendApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetRecvFriendApplication(ctx)
}

func (f *Friend) GetFriendApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetSendFriendApplication(ctx)
}

func (f *Friend) AcceptFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friendPb.ProcessFriendApplicationReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, Flag: constant.FriendResponseAgree, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RefuseFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friendPb.ProcessFriendApplicationReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, Flag: constant.FriendResponseRefuse, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RespondFriendApply(ctx context.Context, req *friendPb.ProcessFriendApplicationReq) error {
	if req.ToUserID == "" {
		req.ToUserID = f.loginUserID
	}
	if err := server_api.ProcessFriendApplication(ctx, req); err != nil {
		return err
	}
	if req.Flag == constant.FriendResponseAgree {
		_ = f.SyncFriends(ctx, []string{req.FromUserID})
	}
	_ = f.SyncAllFriendApplication(ctx)
	return nil
}

// GetUnProcessFriendRequestNum 获取未处理的好友请求
func (f *Friend) GetUnProcessFriendRequestNum(ctx context.Context) (int64, error) {
	return f.db.GetUnProcessFriendRequestNum(ctx, f.loginUserID)
}

func (f *Friend) CheckFriend(ctx context.Context, friendUserIDList []string) ([]*server_api_params.UserIDResult, error) {
	friendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList, true)
	if err != nil {
		return nil, err
	}
	blackList, err := f.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	res := make([]*server_api_params.UserIDResult, 0, len(friendUserIDList))
	for _, v := range friendUserIDList {
		var r server_api_params.UserIDResult
		isBlack := false
		isFriend := false
		for _, b := range blackList {
			if v == b.BlackUserID {
				isBlack = true
				break
			}
		}
		for _, f := range friendList {
			if v == f.FriendUserID {
				isFriend = true
				break
			}
		}
		r.UserID = v
		if isFriend {
			r.Result = 1
		} else if isBlack {
			r.Result = 2
		} else {
			r.Result = 0
		}
		res = append(res, &r)
	}
	return res, nil
}

func (f *Friend) DeleteFriend(ctx context.Context, friendUserID string) error {
	if err := server_api.DeleteFriend(ctx, f.loginUserID, friendUserID); err != nil {
		return err
	}
	return f.deleteFriend(ctx, friendUserID)
}

func (f *Friend) GetFriendList(ctx context.Context) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlackUserID] = localBlackList[i]
	}
	res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return res, nil
}

func (f *Friend) GetFriendListPage(ctx context.Context, offset, count int32) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetPageFriendList(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlackUserID] = localBlackList[i]
	}
	res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return res, nil
}
func (f *Friend) GetFriendsByPage(ctx context.Context, page, size int) (*sdk.FriendPage, error) {
	localFriendList, total, err := f.db.GetFriendsByPage(ctx, page, size)
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlackUserID] = localBlackList[i]
	}
	res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return &sdk.FriendPage{
		Total: total,
		List:  res,
	}, nil
}

// SearchFriendsList 搜索好友
func (f *Friend) SearchFriendsList(ctx context.Context, keyword string, notPeersFriend bool, page, size int) (*sdk.OnlyFriendPage, error) {
	data, total, err := f.db.SearchFriends(ctx, keyword, notPeersFriend, page, size)
	if err != nil {
		return nil, err
	}
	return &sdk.OnlyFriendPage{
		Total: total,
		List:  data,
	}, nil
}

// GetFriendsNotInGroup 不在指定群当中的好友
func (f *Friend) GetFriendsNotInGroup(ctx context.Context, groupID, keyword string, page, size int) (*sdk.OnlyFriendPage, error) {
	data, total, err := f.db.GetFriendsNotInGroup(ctx, groupID, keyword, page, size)
	if err != nil {
		return nil, err
	}
	return &sdk.OnlyFriendPage{
		Total: total,
		List:  data,
	}, nil
}
func (f *Friend) SearchFriends(ctx context.Context, param *sdk.SearchFriendsParam) ([]*sdk.SearchFriendItem, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchNickname && !param.IsSearchUserID && !param.IsSearchRemark) {
		return nil, sdkerrs.ErrArgs.Wrap("keyword is null or search field all false")
	}
	localFriendList, err := f.db.SearchFriendList(ctx, param.KeywordList[0], param.IsSearchUserID, param.IsSearchNickname, param.IsSearchRemark)
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{})
	for _, black := range localBlackList {
		m[black.BlackUserID] = struct{}{}
	}
	res := make([]*sdk.SearchFriendItem, 0, len(localFriendList))
	for i, localFriend := range localFriendList {
		var relationship int
		if _, ok := m[localFriend.FriendUserID]; ok {
			relationship = constant.BlackRelationship
		} else {
			relationship = constant.FriendRelationship
		}
		res = append(res, &sdk.SearchFriendItem{
			LocalFriend:  *localFriendList[i],
			Relationship: relationship,
		})
	}
	return res, nil
}

func (f *Friend) SetFriendRemark(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	return f.SetFriendInfo(ctx, &friendPb.SetFriendInfoReq{
		ToUserID: userIDRemark.ToUserID,
		Remark:   userIDRemark.Remark,
	})
}

func (f *Friend) PinFriends(ctx context.Context, friends *sdk.SetFriendPinParams) error {
	//if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friends.ToUserIDs, IsPinned: friends.IsPinned}, nil); err != nil {
	//	return err
	//}
	//return f.SyncFriends(ctx, friends.ToUserIDs)
	return nil
}

func (f *Friend) SetFriendInfo(ctx context.Context, req *friendPb.SetFriendInfoReq) error {
	req.FromUserID = f.loginUserID
	req.UserID = f.loginUserID
	if err := server_api.SetFriendInfo(ctx, req); err != nil {
		return err
	}
	return f.SyncFriends(ctx, []string{req.ToUserID})
}

func (f *Friend) AddBlack(ctx context.Context, blackUserID string, ex string) error {
	if err := server_api.AddBlack(ctx, f.loginUserID, blackUserID); err != nil {
		return err
	}
	return f.SyncAllBlackList(ctx)
}

func (f *Friend) RemoveBlack(ctx context.Context, blackUserID string) error {
	if err := server_api.RemoveBlack(ctx, f.loginUserID, blackUserID); err != nil {
		return err
	}
	return f.SyncAllBlackList(ctx)
}

func (f *Friend) GetBlackList(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	return f.db.GetBlackListDB(ctx)
}

func (f *Friend) SetFriendsEx(ctx context.Context, friendIDs []string, ex string) error {
	//if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs, Ex: &wrapperspb.StringValue{
	//	Value: ex,
	//}}, nil); err != nil {
	//	return err
	//}
	//// Check if the specified ID is a friend
	//friendResults, err := f.CheckFriend(ctx, friendIDs)
	//if err != nil {
	//	return errs.Wrap(err, "Error checking friend status")
	//}
	//
	//// Determine if friendID is indeed a friend
	//// Iterate over each friendID
	//for _, friendID := range friendIDs {
	//	isFriend := false
	//
	//	// Check if this friendID is in the friendResults
	//	for _, result := range friendResults {
	//		if result.UserID == friendID && result.Result == 1 { // Assuming result 1 means they are friends
	//			isFriend = true
	//			break
	//		}
	//	}
	//
	//	// If this friendID is not a friend, return an error
	//	if !isFriend {
	//		return errs.ErrRecordNotFound.Wrap("Not friend")
	//	}
	//}
	//
	//// If the code reaches here, all friendIDs are confirmed as friends
	//// Update friend information if they are friends
	//
	//updateErr := f.db.UpdateColumnsFriend(ctx, friendIDs, map[string]interface{}{"Ex": ex})
	//if updateErr != nil {
	//	return errs.Wrap(updateErr, "Error updating friend information")
	//}
	//return nil
	return nil
}
