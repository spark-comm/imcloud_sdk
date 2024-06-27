package server_api

import (
	authPb "github.com/OpenIMSDK/protocol/auth"
	"github.com/golang/protobuf/ptypes/empty"
	userPb "github.com/miliao_apis/api/im_cloud/user/v2"
	"github.com/openimsdk/openim-sdk-core/internal/util"
	"github.com/openimsdk/openim-sdk-core/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/pkg/server_api/convert"
	"golang.org/x/net/context"
)

// GetSelfUserInfoFromSvr 从服务端获取个人信息
func GetSelfUserInfoFromSvr(ctx context.Context, loginUserId string) (*model_struct.LocalUser, error) {
	resp := &userPb.GetProfileReply{}
	err := util.CallPostApi[*userPb.GetProfileReq, *userPb.GetProfileReply](
		ctx, constant.GetSelfUserInfoRouter,
		&userPb.GetProfileReq{UserId: loginUserId},
		resp,
	)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	//信息转换
	conversion, err := convert.ServerUserToLocalUser(resp.Profile)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	return conversion, nil
}

// GetServerUserInfo retrieves user information from the server.
func GetServerUserInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	resp := &userPb.FindProfileByUserReply{}
	err := util.CallPostApi[*userPb.FindProfileByUserReq, *userPb.FindProfileByUserReply](
		ctx, constant.GetUsersInfoRouter,
		&userPb.FindProfileByUserReq{UserIds: userIDs},
		resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.List == nil || len(resp.List) == 0 {
		return nil, sdkerrs.ErrUserIDNotFound
	}
	//信息转换
	var users []*model_struct.LocalUser
	for _, v := range resp.List {
		u, err := convert.ServerUserToLocalUser(v)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserLoginStatus retrieves user login status from the server.
func GetUserLoginStatus(ctx context.Context, userIDs string) (*userPb.GetUserLoginStatusReply, error) {
	resp, err := util.ProtoApiPost[userPb.GetUserLoginStatusReq, userPb.GetUserLoginStatusReply](
		ctx,
		constant.GetUserLoginStatusRouter,
		&userPb.GetUserLoginStatusReq{
			UserID: userIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetUsersOption sets the option for the specified user.
func SetUsersOption(ctx context.Context, loginUserId, option string, value int32) error {
	_, err := util.ProtoApiPost[userPb.SetOptionReq, empty.Empty](
		ctx,
		constant.SetUsersOption,
		&userPb.SetOptionReq{
			UserId: loginUserId,
			Option: userPb.UserOption(userPb.UserOption_value[option]),
			Value:  value,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSelfUserInfo updates the user's information.
func UpdateSelfUserInfo(ctx context.Context, userInfo *userPb.UpdateProfileReq) error {
	if _, err := util.ProtoApiPost[userPb.UpdateProfileReq, empty.Empty](
		ctx,
		constant.UpdateSelfUserInfoRouter,
		userInfo,
	); err != nil {
		return err
	}
	return nil
}

// ParseTokenFromSvr parses a token from the server.
func ParseTokenFromSvr(ctx context.Context) (int64, error) {
	resp := &authPb.ParseTokenResp{}
	err := util.CallPostApi[*authPb.ParseTokenReq, *authPb.ParseTokenResp](
		ctx, constant.ParseTokenRouter,
		&authPb.ParseTokenReq{},
		resp,
	)
	return resp.ExpireTimeSeconds, err
}
