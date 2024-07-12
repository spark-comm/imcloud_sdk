package main

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
)

type OnConnListener struct{}

func (c *OnConnListener) OnUserTokenInvalid(errMsg string) {
	//TODO implement me
	panic("implement me")
}

func (c *OnConnListener) OnConnecting() {
	fmt.Println("OnConnecting")
}

func (c *OnConnListener) OnConnectSuccess() {
	fmt.Println("OnConnectSuccess")
}

func (c *OnConnListener) OnConnectFailed(errCode int32, errMsg string) {
	fmt.Println("OnConnectFailed")
}

func (c *OnConnListener) OnKickedOffline() {
	fmt.Println("OnKickedOffline")
}

func (c *OnConnListener) OnUserTokenExpired() {
	fmt.Println("OnUserTokenExpired")
}

type onListenerForService struct {
	ctx context.Context
}

func (o *onListenerForService) OnGroupApplicationAdded(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnGroupApplicationAccepted(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnFriendApplicationAdded(friendApplication string) {
	log.ZDebug(o.ctx, "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onListenerForService) OnFriendApplicationAccepted(groupApplication string) {
	log.ZDebug(o.ctx, "OnFriendApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnRecvNewMessage(message string) {
	log.ZDebug(o.ctx, "OnRecvNewMessage", "message", message)
}

type onConversationListener struct {
	ctx context.Context
}

func (o *onConversationListener) OnDeleteConversation(conversationIds string) {
	//TODO implement me
	panic("implement me")
}

func (o *onConversationListener) OnConversationUserInputStatusChanged(change string) {
	//TODO implement me
	panic("implement me")
}

func (o *onConversationListener) OnDeleteConversation(conversationList string) {
	//TODO implement me
	log.ZDebug(o.ctx, "OnSyncServerProgress")
}

func (o *onConversationListener) OnSyncServerProgress(progress int) {
	//TODO implement me
	log.ZDebug(o.ctx, "OnSyncServerProgress", progress)
}

func (o *onConversationListener) OnSyncServerStart() {
	log.ZDebug(o.ctx, "OnSyncServerStart")
}

func (o *onConversationListener) OnSyncServerFinish() {
	log.ZDebug(o.ctx, "OnSyncServerFinish")
}

func (o *onConversationListener) OnSyncServerFailed() {
	log.ZDebug(o.ctx, "OnSyncServerFailed")
}

func (o *onConversationListener) OnNewConversation(conversationList string) {
	log.ZDebug(o.ctx, "OnNewConversation", "conversationList", conversationList)
}

func (o *onConversationListener) OnConversationChanged(conversationList string) {
	log.ZDebug(o.ctx, "OnConversationChanged", "conversationList", conversationList)
}

func (o *onConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZDebug(o.ctx, "OnTotalUnreadMessageCountChanged", "totalUnreadCount", totalUnreadCount)
}
func (o *onConversationListener) OnClearConversation(conversationId string) {
	log.ZDebug(o.ctx, "OnClearConversation", "conversationId", conversationId)
}

type onGroupListener struct {
	ctx context.Context
}

func (o *onGroupListener) OnGroupDismissed(groupInfo string) {
	log.ZDebug(o.ctx, "OnGroupDismissed", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupAdded(groupInfo string) {
	log.ZDebug(o.ctx, "OnJoinedGroupAdded", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	log.ZDebug(o.ctx, "OnJoinedGroupDeleted", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	log.ZDebug(o.ctx, "OnGroupMemberAdded", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	log.ZDebug(o.ctx, "OnGroupMemberDeleted", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAdded(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationDeleted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupInfoChanged(groupInfo string) {
	log.ZDebug(o.ctx, "OnGroupInfoChanged", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	log.ZDebug(o.ctx, "OnGroupMemberInfoChanged", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationRejected(groupApplication string) {
	log.ZDebug(o.ctx, "OnGroupApplicationRejected", "groupApplication", groupApplication)
}

type onAdvancedMsgListener struct {
	ctx context.Context
}

func (o *onAdvancedMsgListener) OnRecvOnlineOnlyMessage(message string) {
	//TODO implement me
	panic("implement me")
}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessage(message string) {
	//TODO implement me
	panic("implement me")
}

func (o *onAdvancedMsgListener) OnMsgDeleted(message string) {
	log.ZDebug(o.ctx, "OnMsgDeleted", "message", message)
}

//funcation (o *onAdvancedMsgListener) OnMsgDeleted(messageList string) {
//	log.ZDebug(o.ctx, "OnRecvOfflineNewMessages", "messageList", messageList)
//}
//
//funcation (o *onAdvancedMsgListener) OnMsgDeleted(message string) {
//	log.ZDebug(o.ctx, "OnMsgDeleted", "message", message)
//}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessages(messageList string) {
	log.ZDebug(o.ctx, "OnRecvOfflineNewMessages", "messageList", messageList)
}

func (o *onAdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZDebug(o.ctx, "OnRecvNewMessage", "message", message)
}

func (o *onAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.ZDebug(o.ctx, "OnRecvC2CReadReceipt", "msgReceiptList", msgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	log.ZDebug(o.ctx, "OnRecvGroupReadReceipt", "groupMsgReceiptList", groupMsgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvMessageRevoked(msgID string) {
	log.ZDebug(o.ctx, "OnRecvMessageRevoked", "msgID", msgID)
}

func (o *onAdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {
	log.ZDebug(o.ctx, "OnNewRecvMessageRevoked", "messageRevoked", messageRevoked)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.ZDebug(o.ctx, "OnRecvMessageExtensionsChanged", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.ZDebug(o.ctx, "OnRecvMessageExtensionsDeleted", "msgID", msgID, "reactionExtensionKeyList", reactionExtensionKeyList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	log.ZDebug(o.ctx, "OnRecvMessageExtensionsAdded", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

type onFriendListener struct {
	ctx context.Context
}

func (o *onFriendListener) OnFriendApplicationAdded(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationDeleted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationDeleted", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationAccepted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAccepted", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationRejected(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationRejected", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendAdded(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendAdded", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnFriendDeleted(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendDeleted", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnFriendInfoChanged(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendInfoChanged", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnBlackAdded(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackAdded", "blackInfo", blackInfo)
}

func (o *onFriendListener) OnBlackDeleted(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackDeleted", "blackInfo", blackInfo)
}
