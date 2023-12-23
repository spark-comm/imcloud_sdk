package friend

import (
	"context"
	commonPb "github.com/imCloud/api/common"
	friendPb "github.com/imCloud/api/friend/v1"
	"github.com/imCloud/im/pkg/common/log"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
)

func (f *Friend) GetFriendByIdsSvr(ctx context.Context, friendUserIDList []string) ([]*friendPb.FriendInfo, error) {
	res := &friendPb.ListFriendByIdsReply{}
	err := util.CallPostApi[*friendPb.ListFriendByIdsReq, *friendPb.ListFriendByIdsReply](
		ctx, constant.GetFriendByAppIdsRouter,
		&friendPb.ListFriendByIdsReq{UserID: f.loginUserID, FriendIds: friendUserIDList},
		res,
	)
	if err != nil {
		return nil, err
	}
	friends := make([]*friendPb.FriendInfo, 0)
	if res.FriendsInfo == nil {
		log.ZDebug(ctx, "SyncFriendApplicationById res friend request nill")
		return friends, nil
	}
	return res.FriendsInfo, nil
}

// GetFriendBaseInfoSvr 登录同步数据
func (f *Friend) GetFriendBaseInfoSvr(ctx context.Context) ([]*friendPb.SyncFriendInfo, error) {
	req := &friendPb.GetSyncFriendReq{UserID: f.loginUserID, Pagination: &commonPb.RequestPagination{}}
	fn := func(resp *friendPb.GetSyncFriendResp) []*friendPb.SyncFriendInfo {
		return resp.List
	}
	resp := &friendPb.GetSyncFriendResp{}
	respList, err := util.GetPageAll(ctx, constant.GetSyncFriendList, req, resp, fn)
	if err != nil {
		return nil, err
	}
	return respList, nil
}

// SyncFriendInfoByTime 根据时间点同步信息
// par dataTime 时间问题
//
//	返回 新增好友，更新好友，删除好友
func (f *Friend) SyncFriendInfoByTime(ctx context.Context, dataTime map[string]int64) ([]*friendPb.SyncFriendInfo, []*friendPb.SyncFriendInfo, []string, error) {
	res := &friendPb.SyncFriendInfoByTimeReply{}
	err := util.CallPostApi[*friendPb.SyncFriendInfoByTimeReq, *friendPb.SyncFriendInfoByTimeReply](
		ctx, constant.SyncFriendInfoByTimeRouter,
		&friendPb.SyncFriendInfoByTimeReq{UserID: f.loginUserID, TimeData: dataTime},
		res,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return res.AddList, res.UpdateList, res.DelIds, nil
}
