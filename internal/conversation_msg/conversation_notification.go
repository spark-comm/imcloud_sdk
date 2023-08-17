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

package conversation_msg

import (
	"context"
	"encoding/json"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/sdkws"
	utils2 "github.com/imCloud/im/pkg/utils"
)

// Work 会话工作
func (c *Conversation) Work(c2v common.Cmd2Value) {
	log.ZDebug(c2v.Ctx, "NotificationCmd start", "cmd", c2v.Cmd, "value", c2v.Value)
	defer log.ZDebug(c2v.Ctx, "NotificationCmd end", "cmd", c2v.Cmd, "value", c2v.Value)
	switch c2v.Cmd {
	case constant.CmdDeleteConversation:
		c.doDeleteConversation(c2v)
	case constant.CmdNewMsgCome:
		c.doMsgNew(c2v)
	case constant.CmdSuperGroupMsgCome:
		//c.doSuperGroupMsgNew(c2v)
	case constant.CmdUpdateConversation:
		c.doUpdateConversation(c2v)
	case constant.CmdUpdateMessage:
		//修改消息
		c.doUpdateMessage(c2v)
	case constant.CmSyncReactionExtensions:
		//c.doSyncReactionExtensions(c2v)
	case constant.CmdNotification:
		//处理通知
		c.doNotificationNew(c2v)
	case constant.CmdAcceptFriend:
		//好友请求被同意
		c.doAddFriend(c2v)
	case constant.CmdAddFriend:
		//新增会话
		c.doAddFriend(c2v)
	}
}

// doDeleteConversation 删除会话
func (c *Conversation) doDeleteConversation(c2v common.Cmd2Value) {
	node := c2v.Value.(common.DeleteConNode)
	ctx := c2v.Ctx
	// 将会话标记为删除
	// Mark messages related to this conversation for deletion
	err := c.db.UpdateMessageStatusBySourceID(context.Background(), node.SourceID, constant.MsgStatusHasDeleted, int32(node.SessionType))
	if err != nil {
		log.ZError(ctx, "setMessageStatusBySourceID", err)
		return
	}
	//重置清空会话
	//Reset the session information, empty session
	err = c.db.ResetConversation(ctx, node.ConversationID)
	if err != nil {
		log.ZError(ctx, "ResetConversation err:", err)
	}
	//清空服务端会话的消息
	c.clearConversationMsgFromSvr(ctx, node.ConversationID)
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
}

// doAddFriend 新增好友处理
func (c *Conversation) doAddFriend(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	node := c2v.Value.(common.SourceIDAndSessionType)
	cl, err := c.GetOneConversation(ctx, int32(node.SessionType), node.SourceID)
	if err != nil {
		log.ZDebug(ctx, "do add friend error", err)
		return
	}
	c.ConversationListener.OnNewConversation(utils.StructToJsonString([]*model_struct.LocalConversation{cl}))
}

// doUpdateConversation //更新会话
func (c *Conversation) doUpdateConversation(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		// log.Error("internal", "not set conversationListener")
		return
	}
	ctx := c2v.Ctx
	node := c2v.Value.(common.UpdateConNode)
	switch node.Action {
	case constant.AddConOrUpLatMsg:
		var list []*model_struct.LocalConversation
		lc := node.Args.(model_struct.LocalConversation)
		oc, err := c.db.GetConversation(ctx, lc.ConversationID)
		if err == nil {
			// log.Info("this is old conversation", *oc)
			if lc.LatestMsgSendTime >= oc.LatestMsgSendTime { //The session update of asynchronous messages is subject to the latest sending time
				err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"latest_msg_send_time": lc.LatestMsgSendTime, "latest_msg": lc.LatestMsg})
				if err != nil {
					// log.Error("internal", "updateConversationLatestMsgModel err: ", err)
				} else {
					oc.LatestMsgSendTime = lc.LatestMsgSendTime
					oc.LatestMsg = lc.LatestMsg
					list = append(list, oc)
					c.ConversationListener.OnConversationChanged(utils.StructToJsonString(list))
				}
			}
		} else {
			// log.Info("this is new conversation", lc)
			err4 := c.db.InsertConversation(ctx, &lc)
			if err4 != nil {
				// log.Error("internal", "insert new conversation err:", err4.Error())
			} else {
				list = append(list, &lc)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(list))
			}
		}

	case constant.UnreadCountSetZero:
		if err := c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"unread_count": 0}); err != nil {
			log.ZError(ctx, "updateConversationUnreadCountModel err", err, "conversationID", node.ConID)
		} else {
			totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
			if err == nil {
				c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
			} else {
				log.ZError(ctx, "getTotalUnreadMsgCountDB err", err)
			}

		}
	//case ConChange:
	//	err, list := u.getAllConversationListModel()
	//	if err != nil {
	//		sdkLog("getAllConversationListModel database err:", err.Error())
	//	} else {
	//		if list == nil {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString([]ConversationStruct{}))
	//		} else {
	//			u.ConversationListenerx.OnConversationChanged(structToJsonString(list))
	//
	//		}
	//	}
	case constant.IncrUnread:
		err := c.db.IncrConversationUnreadCount(ctx, node.ConID)
		if err != nil {
			// log.Error("internal", "incrConversationUnreadCount database err:", err.Error())
			return
		}
	case constant.TotalUnreadMessageChanged:
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			// log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}
	case constant.UpdateConFaceUrlAndNickName:
		var lc model_struct.LocalConversation
		st := node.Args.(common.SourceIDAndSessionType)
		log.ZInfo(ctx, "UpdateConFaceUrlAndNickName", "st", st)
		switch st.SessionType {
		case constant.SingleChatType:
			lc.UserID = st.SourceID
			lc.ConversationID = c.getConversationIDBySessionType(st.SourceID, constant.SingleChatType)
			lc.ConversationType = constant.SingleChatType
		case constant.SuperGroupChatType:
			conversationID, conversationType, err := c.getConversationTypeByGroupID(ctx, st.SourceID)
			if err != nil {
				// log.Error("internal", "getConversationTypeByGroupID database err:", err.Error())
				return
			}
			lc.GroupID = st.SourceID
			lc.ConversationID = conversationID
			lc.ConversationType = conversationType
		default:
			log.ZError(ctx, "not support sessionType", nil, "sessionType", st.SessionType)
			return
		}
		lc.ShowName = st.Nickname
		lc.FaceURL = st.FaceURL
		err := c.db.UpdateConversation(ctx, &lc)
		if err != nil {
			// log.Error("internal", "setConversationFaceUrlAndNickName database err:", err.Error())
			return
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: lc.ConversationID, Action: constant.ConChange, Args: []string{lc.ConversationID}}})

	case constant.UpdateLatestMessageChange:
		conversationID := node.ConID
		var latestMsg sdk_struct.MsgStruct
		l, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			log.ZError(ctx, "getConversationLatestMsgModel err", err, "conversationID", conversationID)
		} else {
			err := json.Unmarshal([]byte(l.LatestMsg), &latestMsg)
			if err != nil {
				log.ZError(ctx, "latestMsg,Unmarshal err", err)
			} else {
				latestMsg.IsRead = true
				newLatestMessage := utils.StructToJsonString(latestMsg)
				err = c.db.UpdateColumnsConversation(ctx, node.ConID, map[string]interface{}{"latest_msg_send_time": latestMsg.SendTime, "latest_msg": newLatestMessage})
				if err != nil {
					log.ZError(ctx, "updateConversationLatestMsgModel err", err)
				}
			}
		}
	case constant.ConChange:
		conversationIDs := node.Args.([]string)
		conversations, err := c.db.GetMultipleConversationDB(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "getMultipleConversationModel err", err)
		} else {
			var newCList []*model_struct.LocalConversation
			for _, v := range conversations {
				if v.LatestMsgSendTime != 0 {
					newCList = append(newCList, v)
				}
			}
			c.ConversationListener.OnConversationChanged(utils.StructToJsonStringDefault(newCList))
		}
	case constant.NewCon:
		cidList := node.Args.([]string)
		cLists, err := c.db.GetMultipleConversationDB(ctx, cidList)
		if err != nil {
			// log.Error("internal", "getMultipleConversationModel err :", err.Error())
		} else {
			if cLists != nil {
				// log.Info("internal", "getMultipleConversationModel success :", cLists)
				c.ConversationListener.OnNewConversation(utils.StructToJsonString(cLists))
			}
		}
	case constant.ConChangeDirect:
		cidList := node.Args.(string)

		c.ConversationListener.OnConversationChanged(cidList)

	case constant.NewConDirect:
		cidList := node.Args.(string)
		// log.Debug("internal", "NewConversation", cidList)
		c.ConversationListener.OnNewConversation(cidList)

	case constant.ConversationLatestMsgHasRead:
		hasReadMsgList := node.Args.(map[string][]string)
		var result []*model_struct.LocalConversation
		var latestMsg sdk_struct.MsgStruct
		var lc model_struct.LocalConversation
		for conversationID, msgIDList := range hasReadMsgList {
			LocalConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				// log.Error("internal", "get conversation err", err.Error(), conversationID)
				continue
			}
			err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
			if err != nil {
				// log.Error("internal", "JsonStringToStruct err", err.Error(), conversationID)
				continue
			}
			if utils.IsContain(latestMsg.ClientMsgID, msgIDList) {
				latestMsg.IsRead = true
				lc.ConversationID = conversationID
				lc.LatestMsg = utils.StructToJsonString(latestMsg)
				LocalConversation.LatestMsg = utils.StructToJsonString(latestMsg)
				err := c.db.UpdateConversation(ctx, &lc)
				if err != nil {
					// log.Error("internal", "UpdateConversation database err:", err.Error())
					continue
				} else {
					result = append(result, LocalConversation)
				}
			}
		}
		if result != nil {
			// log.Info("internal", "getMultipleConversationModel success :", result)
			c.ConversationListener.OnNewConversation(utils.StructToJsonString(result))
		}
	case constant.SyncConversation:
		// operationID := node.Args.(string)
		// log.Debug(operationID, "reconn sync conversation start")
		c.SyncConversations(ctx)
		err := c.SyncConversationUnreadCount(ctx)
		if err != nil {
			// log.Error(operationID, "reconn sync conversation unread count err", err.Error())
		}
		totalUnreadCount, err := c.db.GetTotalUnreadMsgCountDB(ctx)
		if err != nil {
			// log.Error("internal", "TotalUnreadMessageChanged database err:", err.Error())
		} else {
			c.ConversationListener.OnTotalUnreadMessageCountChanged(totalUnreadCount)
		}

	}
}

func (c *Conversation) doUpdateMessage(c2v common.Cmd2Value) {
	if c.ConversationListener == nil {
		// log.Error("internal", "not set conversationListener")
		return
	}

	node := c2v.Value.(common.UpdateMessageNode)
	ctx := c2v.Ctx
	switch node.Action {
	case constant.UpdateMsgFaceUrlAndNickName:
		args := node.Args.(common.UpdateMessageInfo)
		if args.GroupID == "" {
			if args.UserID == c.loginUserID {
				conversationIDList, err := c.db.GetAllSingleConversationIDList(ctx)
				if err != nil {
					log.ZError(ctx, "GetAllSingleConversationIDList err", err)
					return
				} else {
					log.ZDebug(ctx, "get single conversationID list", "conversationIDList", conversationIDList)
					for _, conversationID := range conversationIDList {
						err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
						if err != nil {
							log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
							continue
						}
					}

				}
			} else {
				conversationID := c.getConversationIDBySessionType(args.UserID, constant.SingleChatType)
				err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
				if err != nil {
					log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
				}

			}
		} else {
			conversationID := c.getConversationIDBySessionType(args.GroupID, constant.SuperGroupChatType)
			err := c.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, conversationID, args.UserID, args.FaceURL, args.Nickname)
			if err != nil {
				log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname err", err)
			}
		}

	}

}

func (c *Conversation) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	if msg.SendTime < c.LoginTime() || c.LoginTime() == 0 {
		log.ZWarn(ctx, "ignore notification", nil, "clientMsgID", msg.ClientMsgID, "serverMsgID",
			msg.ServerMsgID, "seq", msg.Seq, "contentType", msg.ContentType,
			"sendTime", msg.SendTime, "loginTime", c.full.Group().LoginTime())
		return
	}
	if c.msgListener == nil {
		log.ZError(ctx, "msgListner is nil", nil)
		return
	}
	go func() {
		c.SyncConversations(ctx)
	}()
}

func (c *Conversation) doNotificationNew(c2v common.Cmd2Value) {
	ctx := c2v.Ctx
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).Msgs
	syncFlag := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).SyncFlag
	switch syncFlag {
	case constant.MsgSyncBegin:
		//开始同步数据
		c.ConversationListener.OnSyncServerStart()
		if err := c.SyncConversationHashReadSeqs(ctx); err != nil {
			log.ZError(ctx, "SyncConversationHashReadSeqs err", err)
		}
		// 同步数据
		for _, syncFunc := range []func(c context.Context) error{c.user.SyncLoginUserInfo, c.SyncConversations} {
			go func(syncFunc func(c context.Context) error) {
				_ = syncFunc(ctx)
			}(syncFunc)
		}
	case constant.MsgSyncFailed:
		// 同步数据错误
		c.ConversationListener.OnSyncServerFailed()
	case constant.MsgSyncEnd:
		// 同步数据结束
		defer c.ConversationListener.OnSyncServerFinish()
		//同步其他数据
		c.syncOtherInformation(ctx)
	}

	for conversationID, msgs := range allMsg {
		if len(msgs.Msgs) != 0 {
			lastMsg := msgs.Msgs[len(msgs.Msgs)-1]
			log.ZDebug(ctx, "SetNotificationSeq", "conversationID", conversationID, "seq", lastMsg.Seq)
			if lastMsg.Seq != 0 {
				if err := c.db.SetNotificationSeq(ctx, conversationID, lastMsg.Seq); err != nil {
					log.ZError(ctx, "SetNotificationSeq err", err, "conversationID", conversationID, "lastMsg", lastMsg)
				}
			}
		}
		for _, v := range msgs.Msgs {
			switch {
			case v.ContentType == constant.ConversationChangeNotification || v.ContentType == constant.ConversationPrivateChatNotification:
				// 会话改变和私聊通知
				c.DoNotification(ctx, v)
			case v.ContentType == constant.ConversationUnreadNotification:
				//会话未读通知
				var tips sdkws.ConversationHasReadTips
				_ = json.Unmarshal(v.Content, &tips)
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: tips.ConversationID, Action: constant.UnreadCountSetZero}})
				c.db.DeleteConversationUnreadMessageList(ctx, tips.ConversationID, tips.UnreadCountTime)
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{tips.ConversationID}}})
				continue
			case v.ContentType == constant.BusinessNotification:
				//业务通知
				c.business.DoNotification(ctx, v)
				continue
			case v.ContentType == constant.RevokeNotification:
				//撤回通知
				c.doRevokeMsg(ctx, v)
			case v.ContentType == constant.ClearConversationNotification:
				//清空会话通知
				c.doClearConversations(ctx, v)
			case v.ContentType == constant.DeleteMsgsNotification:
				//删除消息通知
				c.doDeleteMsgs(ctx, v)
			case v.ContentType == constant.HasReadReceipt:
				//已读回执
				c.doReadDrawing(ctx, v)
			}
			bChan := make(chan bool, 1)
			go common.ListenerUserInfoChange(ctx, bChan, func() {
				c.group.SyncJoinedGroupMember(ctx)
			})
			switch v.SessionType {
			case constant.SingleChatType:
				//单聊
				if v.ContentType > constant.FriendNotificationBegin && v.ContentType < constant.FriendNotificationEnd {
					c.friend.DoNotification(ctx, v, bChan)
				} else if v.ContentType > constant.UserNotificationBegin && v.ContentType < constant.UserNotificationEnd {
					c.user.DoNotification(ctx, v)
				} else if utils2.Contain(v.ContentType, constant.GroupApplicationRejectedNotification, constant.GroupApplicationAcceptedNotification, constant.JoinGroupApplicationNotification) {
					c.group.DoNotification(ctx, v)
				} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {

					continue
				}
			case constant.GroupChatType, constant.SuperGroupChatType:
				//群聊
				if v.ContentType > constant.GroupNotificationBegin && v.ContentType < constant.GroupNotificationEnd {
					c.group.DoNotification(ctx, v)
				} else if v.ContentType > constant.SignalingNotificationBegin && v.ContentType < constant.SignalingNotificationEnd {
					continue
				}
			}
		}
	}

}

// syncOtherInformation 同步其他信息
func (c *Conversation) syncOtherInformation(ctx context.Context) {
	// 同步数据
	for _, syncFunc := range []func(c context.Context) error{c.friend.SyncFriendList, c.group.SyncJoinedGroup, c.friend.SyncFriendApplication, c.friend.SyncSelfFriendApplication, c.group.SyncAdminGroupApplication, c.group.SyncSelfGroupApplication,
		c.group.SyncJoinedGroupMember, c.friend.SyncBlackList,
	} {
		go func(syncFunc func(c context.Context) error) {
			_ = syncFunc(ctx)
		}(syncFunc)
	}
}
