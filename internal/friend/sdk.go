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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/db/pg"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/server_api_params"
	"sort"
	"strings"

	friendPb "github.com/imCloud/api/friend/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/friend"
)

func (f *Friend) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList)
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

// AddFriend 添加好友
func (f *Friend) AddFriend(ctx context.Context, addRequest *friendPb.AddFriendRequest) error {
	if addRequest.FromUserID == "" {
		addRequest.FromUserID = f.loginUserID
	}
	if err := util.ApiPost(ctx, constant.AddFriendRouter, addRequest, nil); err != nil {
		return err
	}
	return f.SyncFriendApplication(ctx)
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
	if err := util.ApiPost(ctx, constant.AddFriendResponse, req, nil); err != nil {
		return err
	}
	if req.HandleResult == constant.FriendResponseAgree {
		_ = f.SyncFriendList(ctx)
	}
	_ = f.SyncFriendApplication(ctx)
	return nil
	//return f.SyncFriendApplication(ctx)
}

func (f *Friend) CheckFriend(ctx context.Context, friendUserIDList []string) ([]*server_api_params.UserIDResult, error) {
	friendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList)
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
		if isFriend && !isBlack {
			r.Result = 1
		} else {
			r.Result = 0
		}
		res = append(res, &r)
	}
	return res, nil
}

func (f *Friend) DeleteFriend(ctx context.Context, friendUserID string) error {
	if err := util.ApiPost(ctx, constant.DeleteFriendRouter, &friend.DeleteFriendReq{OwnerUserID: f.loginUserID, FriendUserID: friendUserID}, nil); err != nil {
		return err
	}
	//获取会话id
	conversationID := f.getConversationIDBySessionType(friendUserID, constant.SingleChatType)
	//删除好友后删除对应的会话消息
	err := common.TriggerCmdDeleteConversationAndMessage(friendUserID, conversationID, constant.SingleChatType, f.conversationCh)
	if err != nil {
		log.ZDebug(ctx, "delete friend after delete conversation and message")
	}
	return f.syncDelFriend(ctx, friendUserID)
}

func (f *Friend) GetFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	localFriendList, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return nil, err
	}
	//localBlackList, err := f.db.GetBlackListDB(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//m := make(map[string]*model_struct.LocalBlack)
	//for i, black := range localBlackList {
	//	m[black.BlockUserID] = localBlackList[i]
	//}
	//res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	//for _, localFriend := range localFriendList {
	//	res = append(res, &server_api_params.FullUserInfo{
	//		PublicInfo: nil,
	//		FriendInfo: localFriend,
	//		BlackInfo:  m[localFriend.FriendUserID],
	//	})
	//}
	return localFriendList, nil
}

func (f *Friend) GetFriendListPage(ctx context.Context, no, size int64) ([]*model_struct.LocalFriend, error) {
	localFriendList, err := f.db.GetFriendList(ctx, &pg.Page{NO: no, Size: size})
	if err != nil {
		return nil, err
	}
	//localBlackList, err := f.db.GetBlackListDB(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//m := make(map[string]*model_struct.LocalBlack)
	//for i, black := range localBlackList {
	//	m[black.BlockUserID] = localBlackList[i]
	//}
	//res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	//for _, localFriend := range localFriendList {
	//	res = append(res, &server_api_params.FullUserInfo{
	//		PublicInfo: nil,
	//		FriendInfo: localFriend,
	//		BlackInfo:  m[localFriend.FriendUserID],
	//	})
	//}
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

func (f *Friend) SetFriendRemark(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	if err := util.ApiPost(ctx, constant.SetFriendRemark, &friend.SetFriendRemarkReq{OwnerUserID: f.loginUserID, FriendUserID: userIDRemark.ToUserID, Remark: userIDRemark.Remark}, nil); err != nil {
		return err
	}
	return f.SyncFriendList(ctx)
}

func (f *Friend) AddBlack(ctx context.Context, blackUserID string) error {
	if err := util.ApiPost(ctx, constant.AddBlackRouter, &friend.AddBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
		return err
	}
	return f.SyncBlackList(ctx)
}

func (f *Friend) RemoveBlack(ctx context.Context, blackUserID string) error {
	if err := util.ApiPost(ctx, constant.RemoveBlackRouter, &friend.RemoveBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
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

// getConversationIDBySessionType 获取会话类型
func (f *Friend) getConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case constant.SingleChatType:
		l := []string{f.loginUserID, sourceID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		return "g_" + sourceID // group chat
	case constant.SuperGroupChatType:
		return "sg_" + sourceID // super group chat
	case constant.NotificationChatType:
		return "sn_" + sourceID // server notification chat
	}
	return ""
}
