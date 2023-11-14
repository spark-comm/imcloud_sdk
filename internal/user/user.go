// Copyright © 2023 OpenIM SDK. All rights reserved.
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

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/syncer"

	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"

	imUserPb "github.com/imCloud/api/user/v1"
	"github.com/imCloud/im/pkg/common/log"
	authPb "github.com/imCloud/im/pkg/proto/auth"
	"github.com/imCloud/im/pkg/proto/sdkws"
)

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID    string
	listener       open_im_sdk_callback.OnUserListener
	loginTime      int64
	userSyncer     *syncer.Syncer[*model_struct.LocalUser, string]
	conversationCh chan common.Cmd2Value
}

// LoginTime gets the login time of the user.
func (u *User) LoginTime() int64 {
	return u.loginTime
}

// SetLoginTime sets the login time of the user.
func (u *User) SetLoginTime(loginTime int64) {
	u.loginTime = loginTime
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	user := &User{DataBase: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	user.initSyncer()
	return user
}

func (u *User) initSyncer() {
	u.userSyncer = syncer.New(
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return u.InsertLoginUser(ctx, value)
		},
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return fmt.Errorf("not support delete user %s", value.UserID)
		},
		func(ctx context.Context, serverUser, localUser *model_struct.LocalUser) error {
			return u.DataBase.UpdateLoginUser(context.Background(), serverUser)
		},
		func(user *model_struct.LocalUser) string {
			return user.UserID
		},
		nil,
		func(ctx context.Context, state int, server, local *model_struct.LocalUser) error {
			if u.listener == nil {
				return nil
			}
			switch state {
			case syncer.Update:
				u.listener.OnSelfInfoUpdated(utils.StructToJsonString(server))
				if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{UserID: server.UserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationCh)
				}
			}
			return nil
		},
	)
}

// DoNotification handles incoming notifications for the user.
func (u *User) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "user notification", "msg", *msg)
	if u.listener == nil {
		// log.Error(operationID, "listener == nil")
		return
	}
	//小于用户的登录时间忽略通知
	if msg.SendTime < u.loginTime {
		log.ZWarn(ctx, "ignore notification ", nil, "msg", *msg)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.UserInfoUpdatedNotification:
			u.userInfoUpdatedNotification(ctx, msg)
		default:
			// log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

// userInfoUpdatedNotification handles notifications about updated user information.
func (u *User) userInfoUpdatedNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userInfoUpdatedNotification", "msg", *msg)
	tips := sdkws.UserInfoUpdatedTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		log.ZError(ctx, "comm.UnmarshalTips failed", err, "msg", msg.Content)
		return
	}

	if tips.UserID == u.loginUserID {
		u.SyncLoginUserInfo(ctx)
	} else {
		log.ZDebug(ctx, "detail.UserID != u.loginUserID, do nothing", "detail.UserID", tips.UserID, "u.loginUserID", u.loginUserID)
	}
}

// GetUsersInfoFromSvr retrieves user information from the server.
func (u *User) GetUsersInfoFromSvr(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	//resp, err := util.CallApi[imUserPb.FindProfileByUserReply](ctx, constant.GetUsersInfoRouter,
	//	imUserPb.FindProfileByUserReq{UserIds: userIDs})
	resp := &imUserPb.FindProfileByUserReply{}
	err := util.CallPostApi[*imUserPb.FindProfileByUserReq, *imUserPb.FindProfileByUserReply](
		ctx, constant.GetUsersInfoRouter,
		&imUserPb.FindProfileByUserReq{UserIds: userIDs},
		resp,
	)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	//信息转换
	conversion, err := util.BatchConversion(ServerUserToLocalUser, resp.List)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	return conversion, nil
}

// GetUsersInfoFromSvrNoCallback retrieves user information from the server.
func (u *User) GetSingleUserFromSvr(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	users, err := u.GetUsersInfoFromSvr(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, sdkerrs.ErrUserIDNotFound.Wrap(fmt.Sprintf("getSelfUserInfo failed, userID: %s not exist", userID))
}

// getSelfUserInfo retrieves the user's information.
func (u *User) getSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	userInfo, errLocal := u.GetLoginUser(ctx, u.loginUserID)
	if errLocal != nil {
		log.ZError(ctx, fmt.Sprintf("登录的用户id:%s", u.loginUserID), nil)
		srvUserInfo, errServer := u.GetServerUserInfo(ctx, []string{u.loginUserID})
		if errServer != nil {
			return nil, errServer
		}
		log.ZError(ctx, fmt.Sprintf("服务端返回的数据:%s", utils.StructToJsonString(srvUserInfo)), nil)
		if len(srvUserInfo) == 0 {
			return nil, sdkerrs.ErrUserIDNotFound
		}
		ui, err := ServerUserToLocalUser(srvUserInfo[0])
		if err != nil {
			log.ZDebug(ctx, "get self user info error", err)
			return nil, err
		}
		userInfo = ui
		_ = u.InsertLoginUser(ctx, userInfo)
	}
	//填充默认options
	if userInfo.Options == "" {
		option, err := u.getDefUserOption()
		if err != nil {
			return nil, err
		}
		userInfo.Options = option
	}
	return userInfo, nil
}

// searchUser search user info.
func (u *User) searchUser(ctx context.Context, searchValue string, searchType int) (*model_struct.LocalUser, error) {
	//res, err := util.CallApi[imUserPb.SearchProfileReply](ctx, constant.SearchUserInfoRouter,
	//	&imUserPb.SearchProfileReq{SearchValue: searchValue, Type: int32(searchType)})
	res := &imUserPb.SearchProfileReply{}
	err := util.CallPostApi[*imUserPb.SearchProfileReq, *imUserPb.SearchProfileReply](
		ctx, constant.SearchUserInfoRouter,
		&imUserPb.SearchProfileReq{SearchValue: searchValue, Type: int32(searchType)},
		res,
	)
	if err != nil {
		return nil, err
	}
	user, err := ServerUserToLocalUser(res.Profile)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ParseTokenFromSvr parses a token from the server.
func (u *User) ParseTokenFromSvr(ctx context.Context) (int64, error) {
	//resp, err := util.CallApi[authPb.ParseTokenResp](ctx, constant.ParseTokenRouter, authPb.ParseTokenReq{})
	resp := &authPb.ParseTokenResp{}
	err := util.CallPostApi[*authPb.ParseTokenReq, *authPb.ParseTokenResp](
		ctx, constant.ParseTokenRouter,
		&authPb.ParseTokenReq{},
		resp,
	)
	return resp.ExpireTimeSeconds, err
}

// GetServerUserInfo retrieves user information from the server.
func (u *User) GetServerUserInfo(ctx context.Context, userIDs []string) ([]*imUserPb.ProfileReply, error) {
	//resp, err := util.CallApi[imUserPb.FindProfileByUserReply](ctx, constant.GetUsersInfoRouter,
	//	&userPb.GetDesignateUsersReq{UserIDs: userIDs})
	resp := &imUserPb.FindProfileByUserReply{}
	err := util.CallPostApi[*imUserPb.FindProfileByUserReq, *imUserPb.FindProfileByUserReply](
		ctx, constant.GetUsersInfoRouter,
		&imUserPb.FindProfileByUserReq{UserIds: userIDs},
		resp,
	)
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

// updateSelfUserInfo updates the user's information.
func (u *User) updateSelfUserInfo(ctx context.Context, userInfo *imUserPb.UpdateProfileReq) error {
	userInfo.UserId = u.loginUserID
	//if err := util.ApiPost(ctx, constant.UpdateSelfUserInfoRouter, userInfo, nil); err != nil {
	//	return err
	//}
	if _, err := util.ProtoApiPost[imUserPb.UpdateProfileReq, empty.Empty](
		ctx,
		constant.UpdateSelfUserInfoRouter,
		userInfo,
	); err != nil {
		return err
	}
	_ = u.SyncLoginUserInfo(ctx)
	return nil
}

func (u *User) getUserLoginStatus(ctx context.Context, userIDs string) (*imUserPb.GetUserLoginStatusReps, error) {
	//resp := &imUserPb.GetUserLoginStatusReps{}
	//err := util.ApiPost(ctx, constant.GetUserLoginStatusRouter, &imUserPb.GetUserLoginStatusReq{
	//	UserID: userIDs,
	//}, resp)
	resp, err := util.ProtoApiPost[imUserPb.GetUserLoginStatusReq, imUserPb.GetUserLoginStatusReps](
		ctx,
		constant.GetUserLoginStatusRouter,
		&imUserPb.GetUserLoginStatusReq{
			UserID: userIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *User) setUsersOption(ctx context.Context, option string, value int32) error {
	//err := util.ApiPost(ctx, constant.SetUsersOption, &server_api_params.SetOptionReqReq{
	//	UserID: u.loginUserID,
	//	Option: option,
	//	Value:  value,
	//}, nil)
	//
	_, err := util.ProtoApiPost[imUserPb.SetOptionReqReq, empty.Empty](
		ctx,
		constant.SetUsersOption,
		&imUserPb.SetOptionReqReq{
			UserID: u.loginUserID,
			Option: option,
			Value:  value,
		},
	)
	if err != nil {
		return err
	}
	_ = u.SyncLoginUserInfo(ctx)
	return nil
}

func (u *User) screenUserProfile(ctx context.Context, keyWord string) ([]*imUserPb.ScreenUserInfo, error) {
	resp, err := util.ProtoApiPost[imUserPb.ScreenUserInfoReq, imUserPb.ScreenUserInfoResp](
		ctx,
		constant.ScreenUserProfile,
		&imUserPb.ScreenUserInfoReq{
			KeyWord: keyWord,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

const (
	WalletOperation = "is_open_wallet"
)

func (u *User) syncUserOperation(ctx context.Context) error {
	//获取远程operation数据
	//resp := imUserPb.GetOperationResp{}
	//err := util.ApiPost(ctx, constant.GetUserOperation, &imUserPb.GetOperationReq{
	//	OperationKeyWord: []string{
	//		WalletOperation, //只获取钱包数据
	//	},
	//}, &resp)
	//resp, err := util.ProtoApiPost[imUserPb.GetOperationReq, imUserPb.GetOperationResp](
	//	ctx,
	//	constant.GetUserOperation,
	//	&imUserPb.GetOperationReq{},
	//)
	//if err != nil {
	//	return err
	//}
	//respMap := resp.UserOperationMap
	//if _, ok := respMap[WalletOperation]; !ok {
	//	return nil
	//}
	////获取本地数据比对
	//localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	//if err != nil {
	//	return err
	//}
	//marshal, err := json.Marshal(respMap)
	//if err != nil {
	//	return err
	//}
	//return u.DataBase.UpdateLoginUserByMap(ctx, localUser, map[string]interface{}{
	//	"optionst": string(marshal),
	//})
	return u.SyncLoginUserInfo(ctx)
	////获取远程operation数据
	//resp := imUserPb.GetOperationResp{}
	//err := util.ApiPost(ctx, constant.GetUserOperation, &imUserPb.GetOperationReq{}, &resp)
	//if err != nil {
	//	return err
	//}
	//respMap := resp.UserOperationMap
	//if _, ok := respMap[WalletOperation]; !ok {
	//	return nil
	//}
	////获取本地数据比对
	//localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	//if err != nil {
	//	return err
	//}
	//userOperation := make(map[string]int32)
	//err = json.Unmarshal([]byte(localUser.Options), &userOperation)
	//if err != nil {
	//	userOperation[WalletOperation] = respMap[WalletOperation]
	//	marshal, err := json.Marshal(userOperation)
	//	if err != nil {
	//		return err
	//	}
	//	return u.DataBase.UpdateLoginUserByMap(ctx, localUser, map[string]interface{}{
	//		"optionst": string(marshal),
	//	})
	//}
	//if val, ok := userOperation[WalletOperation]; ok {
	//	if val == respMap[WalletOperation] {
	//		return nil
	//	}
	//}
	//userOperation[WalletOperation] = respMap[WalletOperation]
	//marshal, err := json.Marshal(respMap)
	//if err != nil {
	//	return err
	//}
	//return u.DataBase.UpdateLoginUserByMap(ctx, localUser, map[string]interface{}{
	//	"optionst": string(marshal),
	//})
}

// getDefUserOption 获取默认option
func (u *User) getDefUserOption() (string, error) {
	options := map[string]int32{
		"is_real":               0,
		"is_open_moments":       1,
		"group_add":             1,
		"qr_code_add":           1,
		"card_add":              1,
		"code_add":              1,
		"phone_add":             1,
		"show_last_login":       0,
		"multiple_device_login": 0,
		"global_recv_msg_opt":   0,
		"app_manger_level":      0,
		"is_open_wallet":        0,
		"is_admin":              0,
		"not_login_status":      0,
		"is_customer_service":   0,
		"is_tenant":             0,
		"tenant_id":             0,
	}
	marshal, err := json.Marshal(options)
	return string(marshal), err
}
