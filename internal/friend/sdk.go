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
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	friendPb "github.com/imCloud/api/friend/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/friend"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/db/pg"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (f *Friend) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList, true)
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
	res := make([]*server_api_params.FullUserInfo, 0)
	//判断是否有信息未同步的好友
	completed, notCompleteIds := f.CheckCompleted(localFriendList)
	localFriends := make([]*model_struct.LocalFriend, 0)
	if len(notCompleteIds) > 0 || len(localFriendList) == 0 {
		friendInfos, err := f.GetFriendByIdsSvr(ctx, friendUserIDList)
		if err != nil {
			return nil, err
		}
		if len(friendInfos) > 0 {
			serverFriends := util.Batch(ServerFriendToLocalFriend, friendInfos)
			if len(notCompleteIds) == len(localFriendList) {
				localFriends = serverFriends
			} else {
				localFriends = append(completed, serverFriends...)
			}
			go f.syncFriendByInfo(ctx, serverFriends)
		}
	} else {
		localFriends = localFriendList
	}
	for _, localFriend := range localFriends {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return res, nil
}

// AddFriend 添加好友
func (f *Friend) AddFriend(ctx context.Context, addRequest *friendPb.AddFriendRequest) error {
	if addRequest.FromUserID == "" {
		addRequest.FromUserID = f.loginUserID
	}
	//if err := util.ApiPost(ctx, constant.AddFriendRouter, addRequest, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[friendPb.AddFriendRequest, empty.Empty](
		ctx,
		constant.AddFriendRouter,
		addRequest,
	); err != nil {
		return err
	}
	if err := f.SyncFriendApplication(ctx); err != nil {
		return err
	}
	return nil
}

func (f *Friend) GetFriendApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetRecvFriendApplication(ctx)
}

// GetPageFriendApplicationListAsRecipient 分页获取我收到的数据
func (f *Friend) GetPageFriendApplicationListAsRecipient(ctx context.Context, no, size int64) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetRecvFriendApplicationList(ctx, &pg.Page{
		NO:   no,
		Size: size,
	})
}
func (f *Friend) GetFriendApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetSendFriendApplication(ctx)
}
func (f *Friend) GetPageFriendApplicationListAsApplicant(ctx context.Context, no, size int64) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetSendFriendApplicationList(ctx, &pg.Page{
		NO:   no,
		Size: size,
	})
}
func (f *Friend) AcceptFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, HandleResult: constant.FriendResponseAgree, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RefuseFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, HandleResult: constant.FriendResponseRefuse, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RespondFriendApply(ctx context.Context, req *friend.RespondFriendApplyReq) error {
	if req.ToUserID == "" {
		req.ToUserID = f.loginUserID
	}
	if _, err := util.ProtoApiPost[friendPb.AddFriendResponseRequest, empty.Empty](
		ctx,
		constant.AddFriendResponse,
		&friendPb.AddFriendResponseRequest{
			FromUserID: req.FromUserID,
			ToUserID:   req.ToUserID,
			Flag:       int64(req.HandleResult),
			HandleMsg:  req.HandleMsg,
			UserID:     f.loginUserID,
		}); err != nil {
		return err
	}
	if req.HandleResult == constant.FriendResponseAgree {
		_ = f.SyncFriendList(ctx)
	}
	_ = f.SyncFriendApplication(ctx)
	return nil
}

func (f *Friend) CheckFriend(ctx context.Context, friendUserIDList []string) ([]*server_api_params.UserIDResult, error) {
	friendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList, true)
	log.ZInfo(ctx, fmt.Sprintf("获取本地数据的好友列表数据：%+v", friendList))
	if err != nil || len(friendList) != len(friendUserIDList) {
		svr, err := f.GetFriendByIdsSvr(ctx, friendUserIDList)
		log.ZInfo(ctx, fmt.Sprintf("获取远程数据的好友列表数据：%+v", svr))
		if err != nil {
			return nil, err
		}
		friendList = util.Batch(ServerFriendToLocalFriend, svr)
		log.ZInfo(ctx, fmt.Sprintf("本地和远程数据对比处理后的好友列表数据：%+v", friendList))
	}
	blackList, err := f.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	log.ZInfo(ctx, fmt.Sprintf("获取本地的黑名单数据信息为：%+v", blackList))
	res := make([]*server_api_params.UserIDResult, 0, len(friendUserIDList))
	for _, v := range friendUserIDList {
		var r server_api_params.UserIDResult
		isBlack := false
		isFriend := false
		for _, b := range blackList {
			if v == b.BlackUserID || v == b.OwnerUserID {
				isBlack = true
				break
			}
		}
		for _, f := range friendList {
			if f.NotPeersFriend != constant.NotPeersFriend && v == f.FriendUserID {
				isFriend = true
				break
			}
		}
		r.UserID = v
		if isFriend && !isBlack {
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
	//if err := util.ApiPost(ctx, constant.DeleteFriendRouter, &friend.DeleteFriendReq{OwnerUserID: f.loginUserID, FriendUserID: friendUserID}, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[friendPb.DeleteFriendRequest, empty.Empty](
		ctx,
		constant.DeleteFriendRouter,
		&friendPb.DeleteFriendRequest{FromUserID: f.loginUserID, ToUserID: friendUserID},
	); err != nil {
		return err
	}
	//获取会话id
	conversationID := f.getConversationIDBySessionType(friendUserID, constant.SingleChatType)
	//删除好友后删除对应的会话消息
	err := common.TriggerCmdDeleteConversationAndMessage(ctx, friendUserID, conversationID, constant.SingleChatType, f.conversationCh)
	//加密会话处理
	//获取会话id
	ecConversationID := f.getConversationIDBySessionType(friendUserID, constant.EncryptedChatType)
	//删除好友后删除对应的会话消息
	err = common.TriggerCmdDeleteConversationAndMessage(ctx, friendUserID, ecConversationID, constant.EncryptedChatType, f.conversationCh)
	if err != nil {
		log.ZDebug(ctx, "delete friend after delete conversation and message")
	}
	f.RemoveBlack(ctx, friendUserID)
	return f.SyncDelFriend(ctx, friendUserID)
}

func (f *Friend) GetFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	localFriendList, err := f.db.GetAllFriendList(ctx)
	if err != nil || len(localFriendList) == 0 {
		//从远程重新拉取
		err = f.SyncFriend(ctx)
		if err != nil {
			return nil, err
		}
		localFriendList, err = f.db.GetAllFriendList(ctx)
		if err != nil {
			return nil, err
		}
	}

	return localFriendList, nil
}

func (f *Friend) GetFriendListPage(ctx context.Context, no, size int64) ([]*model_struct.LocalFriend, error) {
	localFriendList, err := f.db.GetFriendList(ctx, &pg.Page{NO: no, Size: size})
	if err != nil || len(localFriendList) == 0 {
		//从远程重新拉取
		err = f.SyncFriend(ctx)
		if err != nil {
			return nil, err
		}
		localFriendList, err = f.db.GetFriendList(ctx, &pg.Page{NO: no, Size: size})
		if err != nil {
			return nil, err
		}
	}
	return localFriendList, nil
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

// SetFriendRemark 设置备注
func (f *Friend) SetFriendRemark(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	return f.SetFriendInfo(ctx, userIDRemark)
}

// SetBackgroundUrl 设置聊天背景图片
func (f *Friend) SetBackgroundUrl(ctx context.Context, friendId, backgroundUrl string) error {
	err := f.SetFriendInfo(ctx, &sdk.SetFriendRemarkParams{
		ToUserID:      friendId,
		BackgroundUrl: backgroundUrl,
	})
	if err != nil {
		return err
	}
	//设置会话的聊天背景
	common.TriggerCmdUpdateConversationBackgroundURL(ctx, f.getConversationIDBySessionType(friendId, constant.SingleChatType), backgroundUrl, f.conversationCh)
	return nil
}

// SetFriendInfo 设置好友信息
func (f *Friend) SetFriendInfo(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	//if err := util.ApiPost(ctx, constant.SetFriendInfoRouter, &friendPb.SetFriendInfoRequest{FromUserID: f.loginUserID, ToUserID: userIDRemark.ToUserID, Remark: userIDRemark.Remark, BackgroundUrl: userIDRemark.BackgroundUrl}, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[friendPb.SetFriendInfoRequest, empty.Empty](
		ctx,
		constant.SetFriendInfoRouter,
		&friendPb.SetFriendInfoRequest{
			FromUserID:    f.loginUserID,
			ToUserID:      userIDRemark.ToUserID,
			Remark:        userIDRemark.Remark,
			BackgroundUrl: userIDRemark.BackgroundUrl},
	); err != nil {
		return err
	}
	return f.syncFriendById(ctx, userIDRemark.ToUserID)
}
func (f *Friend) AddBlack(ctx context.Context, blackUserID string) error {
	//if err := util.ApiPost(ctx, constant.AddBlackRouter, &friend.AddBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[friend.AddBlackReq, empty.Empty](
		ctx,
		constant.AddBlackRouter,
		&friend.AddBlackReq{
			OwnerUserID: f.loginUserID,
			BlackUserID: blackUserID}); err != nil {
		return err
	}
	return f.SyncBlackList(ctx)
}

func (f *Friend) RemoveBlack(ctx context.Context, blackUserID string) error {
	//if err := util.ApiPost(ctx, constant.RemoveBlackRouter, &friend.RemoveBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[friendPb.RemoveBlackListRequest, empty.Empty](
		ctx,
		constant.RemoveBlackRouter,
		&friendPb.RemoveBlackListRequest{
			FromUserID: f.loginUserID,
			ToUserID:   blackUserID}); err != nil {
		return err
	}
	return f.SyncBlackList(ctx)
}

func (f *Friend) GetBlackList(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	return f.db.GetBlackListDB(ctx)
}

// GetPageBlackList 分页获取黑名单
func (f *Friend) GetPageBlackList(ctx context.Context, no, size int64) ([]*model_struct.LocalBlack, error) {
	return f.db.GetBlackList(ctx, &pg.Page{NO: no, Size: size})
}

// GetUnprocessedNum 获取待处理的好友请求
func (f *Friend) GetUnprocessedNum(ctx context.Context) (int64, error) {
	return f.db.GetUnprocessedNum(ctx)
}

// SetFriendDestroyMsgStatus 设置好友阅后即焚
// friendID   string 好友状态
// status     int  状态1:开启,0:关闭
func (f *Friend) SetFriendDestroyMsgStatus(ctx context.Context, friendID string, status int) error {
	if _, err := util.ProtoApiPost[friendPb.SetDestroyMsgStatusReq, empty.Empty](
		ctx,
		constant.SetDestroyMsgStatus,
		&friendPb.SetDestroyMsgStatusReq{
			OwnerUserId:  f.loginUserID,
			FriendUserId: friendID,
			Status:       uint32(status),
		},
	); err != nil {
		return err
	}
	return f.syncFriendById(ctx, friendID)
}

// getConversationIDBySessionType 获取会话类型
func (f *Friend) getConversationIDBySessionType(sourceID string, sessionType int) string {
	return utils.GetConversationIDBySessionType(sessionType, f.loginUserID, sourceID)
}

// CheckCompleted 检查数据是否完整
func (f *Friend) CheckCompleted(friends []*model_struct.LocalFriend) ([]*model_struct.LocalFriend, []string) {
	completed := make([]*model_struct.LocalFriend, 0)
	ids := make([]string, 0)
	if friends != nil && len(friends) > 0 {
		for _, v := range friends {
			if v.IsComplete == IsNotComplete {
				ids = append(ids, v.FriendUserID)
			} else {
				completed = append(completed, v)
			}
		}
	}
	return completed, ids
}
