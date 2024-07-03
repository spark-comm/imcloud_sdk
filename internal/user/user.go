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
	"fmt"

	"github.com/OpenIMSDK/protocol/sdkws"
	userPb "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/log"
	"github.com/spark-comm/imcloud_sdk/internal/cache"
	"github.com/spark-comm/imcloud_sdk/pkg/db/db_interface"
	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	"github.com/spark-comm/imcloud_sdk/pkg/sdkerrs"
	"github.com/spark-comm/imcloud_sdk/pkg/server_api"
	"github.com/spark-comm/imcloud_sdk/pkg/syncer"
	usermodel "github.com/spark-comm/spark-api/api/common/model/user/v2"
	imUserPb "github.com/spark-comm/spark-api/api/im_cloud/user/v2"

	"github.com/spark-comm/imcloud_sdk/open_im_sdk_callback"
	"github.com/spark-comm/imcloud_sdk/pkg/common"
	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/utils"
)

type BasicInfo struct {
	Nickname string
	FaceURL  string
}

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID       string
	listener          func() open_im_sdk_callback.OnUserListener
	userSyncer        *syncer.Syncer[*model_struct.LocalUser, string]
	commandSyncer     *syncer.Syncer[*model_struct.LocalUserCommand, string]
	conversationCh    chan common.Cmd2Value
	UserBasicCache    *cache.Cache[string, *BasicInfo]
	OnlineStatusCache *cache.Cache[string, *userPb.OnlineStatus]
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener func() open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	user := &User{DataBase: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	user.initSyncer()
	user.UserBasicCache = cache.NewCache[string, *BasicInfo]()
	user.OnlineStatusCache = cache.NewCache[string, *userPb.OnlineStatus]()
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
			switch state {
			case syncer.Update:
				u.listener().OnSelfInfoUpdated(utils.StructToJsonString(server))
				if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{SessionType: constant.SingleChatType, UserID: server.UserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationCh)
				}
			}
			return nil
		},
	)
	u.commandSyncer = syncer.New(
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to insert a command
			return u.DataBase.ProcessUserCommandAdd(ctx, command)
		},
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to delete a command
			return u.DataBase.ProcessUserCommandDelete(ctx, command)
		},
		func(ctx context.Context, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			// Logic to update a command
			if serverCommand == nil || localCommand == nil {
				return fmt.Errorf("nil command reference")
			}
			return u.DataBase.ProcessUserCommandUpdate(ctx, serverCommand)
		},
		func(command *model_struct.LocalUserCommand) string {
			// Return a unique identifier for the command
			if command == nil {
				return ""
			}
			return command.Uuid
		},
		func(a *model_struct.LocalUserCommand, b *model_struct.LocalUserCommand) bool {
			// Compare two commands to check if they are equal
			if a == nil || b == nil {
				return false
			}
			return a.Uuid == b.Uuid && a.Type == b.Type && a.Value == b.Value
		},
		func(ctx context.Context, state int, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			if u.listener == nil {
				return nil
			}
			switch state {
			case syncer.Delete:
				u.listener().OnUserCommandDelete(utils.StructToJsonString(serverCommand))
			case syncer.Update:
				u.listener().OnUserCommandUpdate(utils.StructToJsonString(serverCommand))
			case syncer.Insert:
				u.listener().OnUserCommandAdd(utils.StructToJsonString(serverCommand))
			}
			return nil
		},
	)
}

// DoNotification handles incoming notifications for the user.
func (u *User) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "user notification", "msg", *msg)
	go func() {
		switch msg.ContentType {
		case constant.UserInfoUpdatedNotification:
			u.userInfoUpdatedNotification(ctx, msg)
		case constant.UserStatusChangeNotification:
			u.userStatusChangeNotification(ctx, msg)
		case constant.UserCommandAddNotification:
			u.userCommandAddNotification(ctx, msg)
		case constant.UserCommandDeleteNotification:
			u.userCommandDeleteNotification(ctx, msg)
		case constant.UserCommandUpdateNotification:
			u.userCommandUpdateNotification(ctx, msg)
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

// userStatusChangeNotification get subscriber status change callback
func (u *User) userStatusChangeNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userStatusChangeNotification", "msg", *msg)
	tips := sdkws.UserStatusChangeTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		log.ZError(ctx, "comm.UnmarshalTips failed", err, "msg", msg.Content)
		return
	}
	if tips.FromUserID == u.loginUserID {
		log.ZDebug(ctx, "self terminal login", "tips", tips)
		return
	}
	u.SyncUserStatus(ctx, tips.FromUserID, tips.Status, tips.PlatformID)
}

// userCommandAddNotification handle notification when user add favorite
func (u *User) userCommandAddNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", *msg)
	tip := sdkws.UserCommandAddTips{}
	if tip.ToUserID == u.loginUserID {
		u.SyncAllCommand(ctx)
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}

// userCommandDeleteNotification handle notification when user delete favorite
func (u *User) userCommandDeleteNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", *msg)
	tip := sdkws.UserCommandDeleteTips{}
	if tip.ToUserID == u.loginUserID {
		u.SyncAllCommand(ctx)
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}

// userCommandUpdateNotification handle notification when user update favorite
func (u *User) userCommandUpdateNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", *msg)
	tip := sdkws.UserCommandUpdateTips{}
	if tip.ToUserID == u.loginUserID {
		u.SyncAllCommand(ctx)
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}

// GetUsersInfoFromSvr retrieves user information from the server.
func (u *User) GetUsersInfoFromSvr(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	resp, err := server_api.GetServerUserInfo(ctx, userIDs)
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	return resp, nil
}

// GetSingleUserFromSvr retrieves user information from the server.
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
		srvUserInfo, errServer := u.GetServerUserInfo(ctx, []string{u.loginUserID})
		if errServer != nil {
			return nil, errServer
		}
		if len(srvUserInfo) == 0 {
			return nil, sdkerrs.ErrUserIDNotFound
		}
		userInfo = srvUserInfo[0]
		_ = u.InsertLoginUser(ctx, userInfo)
	}
	return userInfo, nil
}

// updateSelfUserInfo updates the user's information.
func (u *User) updateSelfUserInfo(ctx context.Context, userInfo *imUserPb.UpdateProfileReq) error {
	userInfo.UserId = u.loginUserID
	if err := server_api.UpdateSelfUserInfo(ctx, userInfo); err != nil {
		return err
	}
	_ = u.SyncLoginUserInfo(ctx)
	return nil
}

// updateSelfUserInfoEx updates the user's information with Ex field.
func (u *User) updateSelfUserInfoEx(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	//userInfo.UserID = u.loginUserID
	//if err := util.ApiPost(ctx, constant.UpdateSelfUserInfoExRouter, userPb.UpdateUserInfoExReq{UserInfo: userInfo}, nil); err != nil {
	//	return err
	//}
	//_ = u.SyncLoginUserInfo(ctx)
	return nil
}

// CRUD user command
func (u *User) ProcessUserCommandAdd(ctx context.Context, userCommand *userPb.ProcessUserCommandAddReq) error {
	//if err := util.ApiPost(ctx, constant.ProcessUserCommandAdd, userPb.ProcessUserCommandAddReq{UserID: u.loginUserID, Type: userCommand.Type, Uuid: userCommand.Uuid, Value: userCommand.Value}, nil); err != nil {
	//	return err
	//}
	//return u.SyncAllCommand(ctx)
	return nil
}

// ProcessUserCommandDelete delete user's choice
func (u *User) ProcessUserCommandDelete(ctx context.Context, userCommand *userPb.ProcessUserCommandDeleteReq) error {
	//if err := util.ApiPost(ctx, constant.ProcessUserCommandDelete, userPb.ProcessUserCommandDeleteReq{UserID: u.loginUserID,
	//	Type: userCommand.Type, Uuid: userCommand.Uuid}, nil); err != nil {
	//	return err
	//}
	//return u.SyncAllCommand(ctx)
	return nil
}

// ProcessUserCommandUpdate update user's choice
func (u *User) ProcessUserCommandUpdate(ctx context.Context, userCommand *userPb.ProcessUserCommandUpdateReq) error {
	//if err := util.ApiPost(ctx, constant.ProcessUserCommandUpdate, userPb.ProcessUserCommandUpdateReq{UserID: u.loginUserID,
	//	Type: userCommand.Type, Uuid: userCommand.Uuid, Value: userCommand.Value}, nil); err != nil {
	//	return err
	//}
	//return u.SyncAllCommand(ctx)
	return nil
}

// ProcessUserCommandGet get user's choice
func (u *User) ProcessUserCommandGetAll(ctx context.Context) ([]*userPb.CommandInfoResp, error) {
	localCommands, err := u.DataBase.ProcessUserCommandGetAll(ctx)
	if err != nil {
		return nil, err // Handle the error appropriately
	}

	var result []*userPb.CommandInfoResp
	for _, localCommand := range localCommands {
		result = append(result, &userPb.CommandInfoResp{
			Type:       localCommand.Type,
			CreateTime: localCommand.CreateTime,
			Uuid:       localCommand.Uuid,
			Value:      localCommand.Value,
		})
	}

	return result, nil
}

// ParseTokenFromSvr parses a token from the server.
func (u *User) ParseTokenFromSvr(ctx context.Context) (int64, error) {
	return server_api.ParseTokenFromSvr(ctx)
}

// GetServerUserInfo retrieves user information from the server.
func (u *User) GetServerUserInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	return server_api.GetServerUserInfo(ctx, userIDs)
}

// subscribeUsersStatus Presence status of subscribed users.
func (u *User) subscribeUsersStatus(ctx context.Context, userIDs []string) ([]*userPb.OnlineStatus, error) {
	//resp, err := util.CallApi[userPb.SubscribeOrCancelUsersStatusResp](ctx, constant.SubscribeUsersStatusRouter, &userPb.SubscribeOrCancelUsersStatusReq{
	//	UserID:  u.loginUserID,
	//	UserIDs: userIDs,
	//	Genre:   PbConstant.SubscriberUser,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return resp.StatusList, nil
	return nil, nil
}

// unsubscribeUsersStatus Unsubscribe a user's presence.
func (u *User) unsubscribeUsersStatus(ctx context.Context, userIDs []string) error {
	//_, err := util.CallApi[userPb.SubscribeOrCancelUsersStatusResp](ctx, constant.SubscribeUsersStatusRouter, &userPb.SubscribeOrCancelUsersStatusReq{
	//	UserID:  u.loginUserID,
	//	UserIDs: userIDs,
	//	Genre:   PbConstant.Unsubscribe,
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

// getSubscribeUsersStatus Get the online status of subscribers.
func (u *User) getSubscribeUsersStatus(ctx context.Context) ([]*userPb.OnlineStatus, error) {
	//resp, err := util.CallApi[userPb.GetSubscribeUsersStatusResp](ctx, constant.GetSubscribeUsersStatusRouter, &userPb.GetSubscribeUsersStatusReq{
	//	UserID: u.loginUserID,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return resp.StatusList, nil
	return nil, nil
}

// getUserStatus Get the online status of users.
func (u *User) getUserStatus(ctx context.Context, userID string) (*usermodel.OnlineStatus, error) {
	resp, err := server_api.GetUserLoginStatus(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &usermodel.OnlineStatus{
		Status: resp.Data.Status,
		UserId: resp.Data.UserId,
	}, nil
}
