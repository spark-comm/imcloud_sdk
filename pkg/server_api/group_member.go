package server_api

import (
	"github.com/brian-god/imcloud_sdk/internal/util"
	"github.com/brian-god/imcloud_sdk/pkg/constant"
	"github.com/brian-god/imcloud_sdk/pkg/db/model_struct"
	"github.com/brian-god/imcloud_sdk/pkg/server_api/convert"
	"github.com/golang/protobuf/ptypes/empty"
	groupmodel "github.com/spark-comm/spark-api/api/common/model/group/v2"
	netmodel "github.com/spark-comm/spark-api/api/common/net/v2"
	v2 "github.com/spark-comm/spark-api/api/im_cloud/group/v2"
	"golang.org/x/net/context"
)

// GetDesignatedGroupMember 获取指定群成员信息
func GetDesignatedGroupMember(ctx context.Context, groupID string, userIDs ...string) ([]*model_struct.LocalGroupMember, error) {
	resp, err := util.ProtoApiPost[v2.MemberByIdsReq, v2.MemberByIdsReply](
		ctx,
		constant.GetGroupMemberByIdsRouter,
		&v2.MemberByIdsReq{
			GroupID: groupID,
			UserIDs: userIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.List == nil {
		return make([]*model_struct.LocalGroupMember, 0), nil
	}
	return util.Batch(convert.ServerGroupMemberToLocalGroupMember, resp.List), nil
}

// GetServerGroupMembers 远程获取群成员
func GetServerGroupMembers(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	req := &groupmodel.GetByGroupListSdk{GroupID: groupID, Pagination: &netmodel.RequestPagination{ShowNumber: 100}}
	fn := func(resp *v2.MemberByIdsReply) []*groupmodel.MemberInfo { return resp.List }
	resp := &v2.MemberByIdsReply{}
	list, err := util.GetPageAll(ctx, constant.GetGroupMemberListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return make([]*model_struct.LocalGroupMember, 0), nil
	}
	return util.Batch(convert.ServerGroupMemberToLocalGroupMember, list), nil
}

// KickGroupMember 踢出群组成员
func KickGroupMember(ctx context.Context, groupID, loginUserId string, reason string, userIDList []string) error {
	if _, err := util.ProtoApiPost[v2.KickGroupMemberReq, empty.Empty](
		ctx,
		constant.KickGroupMemberRouter,
		&v2.KickGroupMemberReq{
			GroupID:          groupID,
			KickedUserIdList: userIDList,
			UserID:           loginUserId,
			HandledMsg:       reason,
		},
	); err != nil {
		return err
	}
	return nil
}

// SetGroupMemberInfo 设置群成员信息
func SetGroupMemberInfo(ctx context.Context, req *v2.SetGroupMemberInfoReq) error {
	_, err := util.ProtoApiPost[v2.SetGroupMemberInfoReq, empty.Empty](
		ctx, constant.SetGroupMemberInfoRouter, req,
	)
	return err
}

// ChangeGroupMemberMute 设置群成员消息接收开关
func ChangeGroupMemberMute(ctx context.Context, groupID, userID, loginUserId string, mutedSeconds int) (err error) {
	if mutedSeconds == 0 {
		req := &v2.CancelMuteGroupMemberReq{GroupID: groupID, PUserID: loginUserId, UserID: userID}
		_, err = util.ProtoApiPost[v2.CancelMuteGroupMemberReq, empty.Empty](
			ctx, constant.CancelMuteGroupMemberRouter, req,
		)
	} else {
		req := &v2.MuteGroupMemberReq{GroupID: groupID, PUserID: loginUserId, UserID: userID, MutedSeconds: int64(mutedSeconds)}
		_, err = util.ProtoApiPost[v2.MuteGroupMemberReq, empty.Empty](
			ctx, constant.MuteGroupMemberRouter, req,
		)
	}
	return
}
