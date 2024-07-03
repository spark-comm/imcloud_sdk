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

// CreateGroup 创建群
func CreateGroup(ctx context.Context, req *v2.CrateGroupReq) (*groupmodel.GroupInfo, error) {
	resp := &v2.CrateGroupReply{}
	err := util.CallPostApi[*v2.CrateGroupReq, *v2.CrateGroupReply](
		ctx, constant.CreateGroupRouter, req, resp,
	)
	return resp.GroupInfo, err
}

// JoinGroup 加入群
func JoinGroup(ctx context.Context, loginUseID, groupID, reqMsg string, joinSource int32) error {
	req := &v2.JoinGroupReq{GroupID: groupID, Remark: reqMsg, SourceID: joinSource, UserID: loginUseID}
	_, err := util.ProtoApiPost[v2.JoinGroupReq, empty.Empty](
		ctx, constant.JoinGroupRouter, req,
	)
	return err
}

// SetGroupInfo 更新群信息
func SetGroupInfo(ctx context.Context, groupInfo *v2.EditGroupProfileReq) error {
	if _, err := util.ProtoApiPost[v2.EditGroupProfileReq, empty.Empty](
		ctx,
		constant.SetGroupInfoRouter,
		groupInfo,
	); err != nil {
		return err
	}
	return nil
}

// SetGroupSwitchInfo 设置群开关
func SetGroupSwitchInfo(ctx context.Context, groupID, loginUseID string, field string, ups int32) error {
	if _, err := util.ProtoApiPost[v2.UpdateGroupSwitchReq, empty.Empty](
		ctx,
		constant.UpdateGroupSwitch,
		&v2.UpdateGroupSwitchReq{
			UserID:  loginUseID,
			GroupID: groupID,
			Field:   v2.GroupSwitchOption(v2.GroupSwitchOption_value[field]),
			Updates: ups,
		},
	); err != nil {
		return err
	}
	return nil
}

// QuitGroup 退出群
func QuitGroup(ctx context.Context, groupID, loginUseID string) error {
	req := &v2.QuitGroupReq{GroupID: groupID, UserID: loginUseID}
	_, err := util.ProtoApiPost[v2.QuitGroupReq, empty.Empty](
		ctx, constant.QuitGroupRouter, req,
	)
	return err
}

// DismissGroup 解散群
func DismissGroup(ctx context.Context, groupID, loginUseID string) error {
	req := &v2.DismissGroupNoticeReq{GroupID: groupID, UserID: loginUseID}
	_, err := util.ProtoApiPost[v2.DismissGroupNoticeReq, empty.Empty](
		ctx, constant.DismissGroupRouter, req,
	)
	return err
}

// ChangeGroupMute 设置群消息接收开关
func ChangeGroupMute(ctx context.Context, groupID, loginUseID string, isMute bool) (err error) {
	if isMute {
		req := &v2.MuteGroupReq{GroupID: groupID, PUserID: loginUseID}
		_, err = util.ProtoApiPost[v2.MuteGroupReq, empty.Empty](
			ctx, constant.MuteGroupRouter, req,
		)
	} else {
		req := &v2.CancelMuteGroupReq{GroupID: groupID, PUserID: loginUseID}
		_, err = util.ProtoApiPost[v2.CancelMuteGroupReq, empty.Empty](
			ctx, constant.CancelMuteGroupRouter, req,
		)
	}
	return
}

// GetSpecifiedGroupsInfo 获取群信息
func GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	resp, err := util.ProtoApiPost[v2.GetFullGroupInfoReq, v2.GetFullGroupInfoReply](
		ctx,
		constant.FindFullGroupInfoRouter,
		&v2.GetFullGroupInfoReq{
			GroupIds: groupIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.List == nil {
		return make([]*model_struct.LocalGroup, 0), nil
	}
	return util.Batch(convert.ServerBaseGroupToLocalGroup, resp.List), nil
}

// GetServerJoinGroup  获取服务端用户加入的群
func GetServerJoinGroup(ctx context.Context, loginUserId string) ([]*model_struct.LocalGroup, error) {
	fn := func(resp *v2.GetJoinedGroupListReply) []*groupmodel.BaseGroupInfo { return resp.List }
	req := &netmodel.GetByFormUserListSdk{FromUserID: loginUserId, Pagination: &netmodel.RequestPagination{}}
	resp := &v2.GetJoinedGroupListReply{}
	list, err := util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerBaseGroupToLocalGroup, list), nil
}

// GetGroupsInfo 从服务端获取群数据
func GetGroupsInfo(ctx context.Context, groupIDs ...string) ([]*model_struct.LocalGroup, error) {
	resp := &v2.GetGroupInfoReply{}
	err := util.CallPostApi[*v2.GetGroupInfoReq, *v2.GetGroupInfoReply](
		ctx, constant.GetGroupsInfoRouter, &v2.GetGroupInfoReq{GroupID: groupIDs}, resp,
	)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerGroupToLocalGroup, resp.List), nil
}

// TransferGroupOwner 转让群主
func TransferGroupOwner(ctx context.Context, groupID, loginUserId, newOwnerUserID string) error {
	if _, err := util.ProtoApiPost[v2.TransferGroupReq, empty.Empty](
		ctx,
		constant.TransferGroupRouter,
		&v2.TransferGroupReq{
			GroupID:        groupID,
			NewOwnerUserID: newOwnerUserID,
			UserID:         loginUserId,
		},
	); err != nil {
		return err
	}
	return nil
}

// InviteUserToGroup 邀请用户进群
func InviteUserToGroup(ctx context.Context, groupID, loginUserId, reason string, userIDList []string) error {
	if _, err := util.ProtoApiPost[v2.InviteUserToGroupReq, empty.Empty](
		ctx,
		constant.InviteUserToGroupRouter,
		&v2.InviteUserToGroupReq{
			GroupID:           groupID,
			InvitedUserIdList: userIDList,
			Reason:            reason,
			UserID:            loginUserId,
		},
	); err != nil {
		return nil
	}
	return nil
}

// SearchGroupByCode  搜索群组
func SearchGroupByCode(ctx context.Context, loginUserId, groupCode string) (*groupmodel.GroupInfo, error) {
	resp, err := util.ProtoApiPost[v2.GetGroupByCodeReq, v2.GetGroupByCodeReply](
		ctx, constant.SearchGroupByCodeRouter, &v2.GetGroupByCodeReq{UserID: loginUserId, Code: groupCode},
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
