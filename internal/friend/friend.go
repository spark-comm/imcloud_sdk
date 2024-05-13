// Copyright 2021 OpenIM Corporation
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

package friend

import (
	"context"
	"fmt"
	"github.com/imCloud/api/common/enums"
	"github.com/imCloud/api/common/notice"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/sdkws"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/delayqueue"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
)

func NewFriend(ctx context.Context, loginUserID string, db db_interface.DataBase, user *user.User, conversationCh, groupCh chan common.Cmd2Value) *Friend {
	f := &Friend{loginUserID: loginUserID, db: db, user: user, conversationCh: conversationCh, groupCh: groupCh, syncFriendQueue: delayqueue.New[int]()}
	f.initSyncer()
	return f
}

type Friend struct {
	friendListener    open_im_sdk_callback.OnFriendshipListenerSdk
	loginUserID       string
	db                db_interface.DataBase
	user              *user.User
	friendSyncer      *syncer.Syncer[*model_struct.LocalFriend, [2]string]
	blockSyncer       *syncer.Syncer[*model_struct.LocalBlack, [2]string]
	requestRecvSyncer *syncer.Syncer[*model_struct.LocalFriendRequest, [2]string]
	requestSendSyncer *syncer.Syncer[*model_struct.LocalFriendRequest, [2]string]
	loginTime         int64
	conversationCh    chan common.Cmd2Value
	groupCh           chan common.Cmd2Value
	// 同步群组信息延迟队列
	syncFriendQueue    *delayqueue.DelayQueue[int]
	listenerForService open_im_sdk_callback.OnListenerForService
}

func (f *Friend) initSyncer() {
	//好友信息同步
	f.friendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.InsertFriend(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.DeleteFriendDB(ctx, value.FriendUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriend, local *model_struct.LocalFriend) error {
		return f.db.UpdateFriend(ctx, server)
	}, func(value *model_struct.LocalFriend) [2]string {
		return [...]string{value.OwnerUserID, value.FriendUserID}
	}, nil, f.NotificationFun)
	//黑名单同步
	f.blockSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.DeleteBlack(ctx, value.OwnerUserID, value.BlackUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return f.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlackUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalBlack) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnBlackAdded(*server)
		case syncer.Delete:
			f.friendListener.OnBlackDeleted(*local)
		}
		return nil
	})
	// 收到的申请同步
	f.requestRecvSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		err := f.db.InsertFriendRequest(ctx, value)
		if err != nil {
			log.ZInfo(ctx, fmt.Sprintf("收到的申请同步InsertFriendRequest失败：%+v,insert数据：%+v", err, value))
			return err
		}
		return nil
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalFriendRequest) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*server)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*local)
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*server)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*server)
			case constant.FriendResponseDefault:
				f.friendListener.OnFriendApplicationAdded(*server)
			}
		}
		return nil
	})
	//发起的申请同步
	f.requestSendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		err := f.db.InsertFriendRequest(ctx, value)
		if err != nil {
			log.ZInfo(ctx, fmt.Sprintf("发起的申请同步InsertFriendRequest失败：%+v,insert数据：%+v", err, value))
			return err
		}
		return nil
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalFriendRequest) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*server)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*local)
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*server)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*server)
			case constant.FriendResponseDefault:
				f.friendListener.OnFriendApplicationAdded(*server)
			}
		}
		return nil
	})
}

func (f *Friend) LoginTime() int64 {
	return f.loginTime
}

func (f *Friend) SetLoginTime(loginTime int64) {
	f.loginTime = loginTime
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func (f *Friend) SetListener(listener open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = open_im_sdk_callback.NewOnFriendshipListenerSdk(listener)
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}

func (f *Friend) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := f.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "doNotification error", err, "msg", msg)
		}
	}()
}

// syncApplication 同步好友申请
func (f *Friend) syncApplication(ctx context.Context, from *notice.FromToUserID) error {
	if from.Operation == enums.Operation_Delete {
		f.friendListener.OnFriendApplicationDeleted(model_struct.LocalFriendRequest{
			FromUserID: from.FromUserID,
			ToUserID:   from.ToUserID,
		})
	}
	if from.FromUserID == f.loginUserID {
		// 自己发起的好友请求
		return f.SyncSelfFriendApplication(ctx)
	} else if from.ToUserID == f.loginUserID {
		// 发给自己的请求,同步自己收到的好友请求
		return f.SyncFriendApplication(ctx)
	}
	return nil
}

// syncApplicationByNotification 根据通知同步好友请求
func (f *Friend) syncApplicationByNotification(ctx context.Context, from *sdkws.FromToUserID) error {
	err := f.syncFriendApplicationById(ctx, from.FromUserID, from.ToUserID)
	if err != nil {
		return fmt.Errorf("friend application notification error, fromUserID: %s, toUserID: %s", from.FromUserID, from.ToUserID)
	}
	log.ZInfo(ctx, "根据通知同步好友请求成功！")
	return nil
}

// syncFriendByNotification
func (f *Friend) syncFriendByNotification(ctx context.Context, friendId string) error {
	err := f.syncFriendById(ctx, friendId)
	if err != nil {
		return fmt.Errorf("friend  notification error, fromUserID: %s, toUserID: %s", f.loginUserID, friendId)
	}
	//生成对应的会话
	//_ = common.TriggerCmdAddFriendGenerateSession(ctx, common.SourceIDAndSessionType{SourceID: friendId, SessionType: constant.SingleChatType}, f.conversationCh)
	return nil
}

// DelLocalFriend 删除本地好友
func (f *Friend) DelLocalFriend(ctx context.Context, friendId string) error {
	return f.db.DeleteFriendDB(ctx, friendId)
}

// NotificationFun 好友通知
func (f *Friend) NotificationFun(ctx context.Context, state int, server, local *model_struct.LocalFriend) error {
	if f.friendListener == nil {
		return nil
	}
	switch state {
	case syncer.Insert:
		//尝试更新会话
		f.friendListener.OnFriendAdded(*server)
		if server.Remark != "" {
			server.Nickname = server.Remark
		}
		_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
			Args: common.SourceIDAndSessionType{SourceID: server.FriendUserID, SessionType: constant.SingleChatType, FaceURL: server.FaceURL, Nickname: server.Nickname}}, f.conversationCh)
	case syncer.Delete:
		log.ZDebug(ctx, "syncer OnFriendDeleted", "local", local)
		f.friendListener.OnFriendDeleted(*local)
	case syncer.Update:
		f.friendListener.OnFriendInfoChanged(*server)
		if local.Nickname != server.Nickname || local.FaceURL != server.FaceURL || local.Remark != server.Remark {
			if server.Remark != "" {
				server.Nickname = server.Remark
			}
			//更新会话
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
				Args: common.SourceIDAndSessionType{SourceID: server.FriendUserID, SessionType: constant.SingleChatType, FaceURL: server.FaceURL, Nickname: server.Nickname}}, f.conversationCh)
			//更新消息
			_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
				Args: common.UpdateMessageInfo{UserID: server.FriendUserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, f.conversationCh)
			//更新所在群的信息
			_ = common.TriggerCmdGroupMemberChange(ctx, common.UpdateGroupMemberInfo{UserId: server.FriendUserID, Nickname: server.Nickname, FaceUrl: server.FaceURL}, f.groupCh)
		}

	}
	return nil
}

// DelFriendConversation 删除群会话
func (f *Friend) DelFriendConversation(ctx context.Context, friend string) {
	//删除会话
	conversationID := utils.GetConversationIDBySessionType(constant.SingleChatType, friend)
	err := common.TriggerCmdDeleteConversationAndMessage(
		ctx,
		friend,
		conversationID,
		constant.SingleChatType,
		f.conversationCh)
	if err != nil {
		log.ZDebug(ctx, "QuitGroup  after delete conversation err", err)
	}
}
