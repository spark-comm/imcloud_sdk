package server_api

import (
	authPb "github.com/OpenIMSDK/protocol/auth"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spark-comm/imcloud_sdk/internal/util"
	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	"github.com/spark-comm/imcloud_sdk/pkg/sdkerrs"
	"github.com/spark-comm/imcloud_sdk/pkg/server_api/convert"
	usermodel "github.com/spark-comm/spark-api/api/common/model/user/v2"
	userPb "github.com/spark-comm/spark-api/api/im_cloud/user/v2"
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

// SearchUser 搜索用户
func SearchUser(ctx context.Context, req *userPb.SearchProfileReq) (*usermodel.UserProfile, error) {
	resp := &userPb.SearchProfileReply{}
	err := util.CallPostApi[*userPb.SearchProfileReq, *userPb.SearchProfileReply](
		ctx, constant.SearchUserInfoRouter,
		req,
		resp,
	)
	if err != nil {
		return nil, sdkerrs.Warp(err, "SearchUser failed")
	}
	return resp.Data, nil
}

// FindFullProfileByUserId 获取完整用户信息
func FindFullProfileByUserId(ctx context.Context, userIDs ...string) ([]*usermodel.UserProfile, error) {
	resp := &userPb.FindFullProfileByUserIdReply{}
	err := util.CallPostApi[*userPb.FindFullProfileByUserIdReq, *userPb.FindFullProfileByUserIdReply](
		ctx, constant.FindFullProfileByUserIdRouter,
		&userPb.FindFullProfileByUserIdReq{UserIds: userIDs},
		resp)
	if err != nil {
		return nil, sdkerrs.Warp(err, "FindFullProfileByUserId failed")
	}
	return resp.List, nil
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

// GetUsersOption sets the option for the specified user.
func GetUsersOption(ctx context.Context, loginUserId, option string) (int32, error) {
	res, err := util.ProtoApiPost[userPb.GetOptionValReq, userPb.GetOptionValReply](
		ctx,
		constant.GetUserOperation,
		&userPb.GetOptionValReq{
			UserId: loginUserId,
			Option: userPb.UserOption(userPb.UserOption_value[option]),
		},
	)
	if err != nil {
		return 0, err
	}
	return res.Value, nil
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
