package user

import (
	"context"
	imUserPb "github.com/imCloud/api/user/v1"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdkerrs"
)

// GetSelfUserInfoFromSvr 从服务端获取个人信息
func (u *User) GetSelfUserInfoFromSvr(ctx context.Context) (*model_struct.LocalUser, error) {
	resp := &imUserPb.FindProfileByIdNoCacheReply{}
	err := util.CallPostApi[*imUserPb.FindProfileByIdNoCacheReq, *imUserPb.FindProfileByIdNoCacheReply](
		ctx, constant.GetSelfUserInfoRouter,
		&imUserPb.FindProfileByIdNoCacheReq{UserId: u.loginUserID},
		resp,
	)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	//信息转换
	conversion, err := ServerUserToLocalUser(resp.Profile)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	return conversion, nil
}
