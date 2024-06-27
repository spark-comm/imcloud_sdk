package server_api

import (
	"github.com/golang/protobuf/ptypes/empty"
	friendmodel "github.com/miliao_apis/api/common/model/friend/v2"
	netmodel "github.com/miliao_apis/api/common/net/v2"
	v2 "github.com/miliao_apis/api/im_cloud/friend/v2"
	"github.com/openimsdk/openim-sdk-core/internal/util"
	"github.com/openimsdk/openim-sdk-core/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/pkg/server_api/convert"
	"golang.org/x/net/context"
)

// ProcessFriendApplication 处理好友申请
func ProcessFriendApplication(ctx context.Context, req *v2.ProcessFriendApplicationReq) error {
	if _, err := util.ProtoApiPost[v2.ProcessFriendApplicationReq, empty.Empty](
		ctx,
		constant.AddFriendResponse,
		req,
	); err != nil {
		return err
	}
	return nil
}

// BothFriendRequest 好友请求
func BothFriendRequest(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	res := &v2.GetFriendRequestByApplicantReply{}
	err := util.CallPostApi[*v2.GetFriendRequestByApplicantReq, *v2.GetFriendRequestByApplicantReply](
		ctx, constant.GetFriendRequestByApplicantRouter,
		&v2.GetFriendRequestByApplicantReq{FromUserID: fromUserID, ToUserID: toUserID},
		res,
	)
	if err != nil {
		return nil, err
	}
	if res.FriendRequest == nil {
		return nil, sdkerrs.ErrUserIDNotFound.Wrap("sync friend request failed")
	}
	return convert.ServerFriendRequestToLocalFriendRequest(res.FriendRequest), nil
}

// GetSendFriendApplication 自己发送的好友请求
func GetSendFriendApplication(ctx context.Context, loginUserId string) ([]*model_struct.LocalFriendRequest, error) {
	req := &netmodel.GetByUserListSdk{UserID: loginUserId, Pagination: &netmodel.RequestPagination{}}
	fn := func(resp *v2.GetSendFriendsApplyReply) []*friendmodel.FriendRequest {
		return resp.List
	}
	resp := &v2.GetSendFriendsApplyReply{}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	if requests == nil {
		return nil, sdkerrs.ErrUserIDNotFound.Wrap("sync friend request failed")
	}
	return util.Batch(convert.ServerFriendRequestToLocalFriendRequest, requests), nil
}

// GetReceiveFriendApplication 获取收到的好友申请列表
func GetReceiveFriendApplication(ctx context.Context, loginUserId string) ([]*model_struct.LocalFriendRequest, error) {
	req := &netmodel.GetByUserListSdk{UserID: loginUserId, Pagination: &netmodel.RequestPagination{}}
	fn := func(resp *v2.GetReceiveFriendsApplyReply) []*friendmodel.FriendRequest {
		return resp.List
	}
	resp := &v2.GetReceiveFriendsApplyReply{}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendReceiveApplicationListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	if requests == nil {
		return nil, sdkerrs.ErrUserIDNotFound.Wrap("sync friend request failed")
	}
	return util.Batch(convert.ServerFriendRequestToLocalFriendRequest, requests), nil
}
