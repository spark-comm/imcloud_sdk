package server_api

import (
	"github.com/OpenIMSDK/tools/log"
	"github.com/golang/protobuf/ptypes/empty"
	friendmodel "github.com/miliao_apis/api/common/model/friend/v2"
	netmodel "github.com/miliao_apis/api/common/net/v2"
	friendPb "github.com/miliao_apis/api/im_cloud/friend/v2"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api/convert"
	"golang.org/x/net/context"
)

// AddFriend 添加好友
func AddFriend(ctx context.Context, addRequest *friendPb.AddFriendReq) error {
	if _, err := util.ProtoApiPost[friendPb.AddFriendReq, empty.Empty](
		ctx,
		constant.AddFriendRouter,
		addRequest,
	); err != nil {
		return err
	}
	return nil
}

// DeleteFriend 删除好友
func DeleteFriend(ctx context.Context, loginUserID, friendUserID string) error {
	if _, err := util.ProtoApiPost[friendPb.DeleteFriendReq, empty.Empty](
		ctx,
		constant.DeleteFriendRouter,
		&friendPb.DeleteFriendReq{FromUserID: loginUserID, ToUserID: friendUserID},
	); err != nil {
		return err
	}
	return nil
}

// SetFriendInfo 设置好友信息
func SetFriendInfo(ctx context.Context, req *friendPb.SetFriendInfoReq) error {
	if _, err := util.ProtoApiPost[friendPb.SetFriendInfoReq, empty.Empty](
		ctx,
		constant.SetFriendInfoRouter,
		req,
	); err != nil {
		return err
	}
	return nil
}

// AddBlack 添加黑名单
func AddBlack(ctx context.Context, loginUserID, blackUserID string) error {
	if _, err := util.ProtoApiPost[friendPb.AddBlackReq, empty.Empty](
		ctx,
		constant.AddBlackRouter,
		&friendPb.AddBlackReq{
			OwnerUserID: loginUserID,
			BlackUserID: blackUserID,
		}); err != nil {
		return err
	}
	return nil
}

// RemoveBlack 移除黑名单
func RemoveBlack(ctx context.Context, loginUserID, blackUserID string) error {
	if _, err := util.ProtoApiPost[friendPb.RemoveBlackListReq, empty.Empty](
		ctx,
		constant.RemoveBlackRouter,
		&friendPb.RemoveBlackListReq{
			FromUserID: loginUserID,
			ToUserID:   blackUserID}); err != nil {
		return err
	}
	return nil
}

// GetAllFriendList 获取所有的好友
func GetAllFriendList(ctx context.Context, loginUserId string) ([]*model_struct.LocalFriend, error) {
	req := &netmodel.GetByUserListSdk{UserID: loginUserId, Pagination: &netmodel.RequestPagination{}}
	fn := func(resp *friendPb.ListFriendReply) []*friendmodel.FriendInfo {
		return resp.List
	}
	resp := &friendPb.ListFriendReply{}
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	if friends == nil {
		return nil, sdkerrs.ErrUserIDNotFound.Wrap("friend failed not found")
	}
	return util.Batch(convert.ServerFriendToLocalFriend, friends), nil
}

// GetFriendByIds 获取好友列表
func GetFriendByIds(ctx context.Context, loginUserId string, friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	res := &friendPb.ListFriendByIdsReply{}
	err := util.CallPostApi[*friendPb.ListFriendByIdsReq, *friendPb.ListFriendByIdsReply](
		ctx, constant.GetFriendByAppIdsRouter,
		&friendPb.ListFriendByIdsReq{UserID: loginUserId, FriendIds: friendUserIDList},
		res,
	)
	if err != nil {
		return nil, err
	}
	if res.List == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nill")
		return make([]*model_struct.LocalFriend, 0), nil
	}
	return util.Batch(convert.ServerFriendToLocalFriend, res.List), nil
}

// GetAllBlackList 获取黑名单信息
func GetAllBlackList(ctx context.Context, loginUserId string) ([]*model_struct.LocalBlack, error) {
	req := &netmodel.GetByUserListSdk{UserID: loginUserId, Pagination: &netmodel.RequestPagination{}}
	fn := func(resp *friendPb.BlackListReply) []*friendmodel.BlackInfo { return resp.List }
	resp := &friendPb.BlackListReply{}
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	if serverData == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nill")
		return make([]*model_struct.LocalBlack, 0), nil
	}
	log.ZDebug(ctx, "black from local", "data", resp)
	return util.Batch(convert.ServerBlackToLocalBlack, serverData), nil
}
