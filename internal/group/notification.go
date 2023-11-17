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

package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/imCloud/api/im/v1"
	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/sdkws"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

func (g *Group) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := g.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "DoGroupNotification failed", err)
		}
	}()
}

func (g *Group) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if g.listener == nil {
		return errors.New("listener is nil")
	}
	//如果通知消息的时间小于登录时间则不进行通知
	if msg.SendTime < g.loginTime {
		return nil
	}
	switch msg.ContentType {
	//创建群消息通知
	case constant.GroupCreatedNotification: // 1501
		var detail sdkws.GroupCreatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		//同步群信息
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		//同步群成员
		return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
		//群信息变更通知
	case constant.GroupInfoSetNotification: // 1502
		var detail v1.GroupInfoSetNotificationReq
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
		//加群请求通知
	case constant.JoinGroupApplicationNotification: // 1503
		var detail sdkws.JoinGroupApplicationTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.Applicant.UserID == g.loginUserID {
			//自己主动加群申请信息
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
		} else {
			//(以管理员或群主身份)获取群的加群申请
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
		}
		//接受群请求通知
	case constant.GroupApplicationAcceptedNotification: // 1505
		var detail sdkws.GroupApplicationAcceptedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			//return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
			//只同步单个群的申请情况
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
		} else if detail.ReceiverAs == 1 {
			//只同步单个群的申请情况
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
			//return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
		g.syncGroupMembers(ctx, detail.Group.GroupID, detail.OpUser.InviterUserID)
		return g.SyncGroups(ctx, detail.Group.GroupID)
		//群请求拒绝通知
	case constant.GroupApplicationRejectedNotification: // 1506
		var detail sdkws.GroupApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			//return g.SyncAdminGroupApplication(ctx)
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
		} else if detail.ReceiverAs == 1 {
			//return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
			return g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID)
		}
		return nil
		//return g.SyncSelfGroupApplications(ctx, detail.Group.GroupID)
		//转让群主通知
	case constant.GroupOwnerTransferredNotification: // 1507
		var detail sdkws.GroupOwnerTransferredTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		if detail.Group == nil {
			return errors.New(fmt.Sprintf("group is nil, groupID: %s", detail.Group.GroupID))
		}
		//同步指定的群成员信息
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.OpUser.UserID, detail.NewGroupOwner.UserID)
		//踢人通知
	case constant.MemberKickedNotification: // 1508
		var detail sdkws.MemberKickedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		var self bool
		for _, info := range detail.KickedUserList {
			if info.UserID == g.loginUserID {
				self = true
				break
			}
		}
		//踢掉自己
		if self {
			// 删除群
			return g.deleteGroup(ctx, detail.Group.GroupID)
		} else { //仅群成员信息变动
			var userIDs []string
			for _, info := range detail.KickedUserList {
				userIDs = append(userIDs, info.UserID)
			}
			//删除指定成员
			return g.deleteGroupMembers(ctx, detail.Group.GroupID, userIDs...)
		}
		//退群通知
	case constant.MemberQuitNotification: // 1504
		var detail sdkws.MemberQuitTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.QuitUser.UserID == g.loginUserID {
			//退群为自己  删除操作
			return g.deleteGroup(ctx, detail.Group.GroupID)
		} else { //群成员同步
			return g.deleteGroupMembers(ctx, detail.Group.GroupID, detail.QuitUser.UserID)
		}
		//成员邀请通知
	case constant.MemberInvitedNotification: // 1509
		var detail sdkws.MemberInvitedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		var userIDs []string
		for _, info := range detail.InvitedUserList {
			userIDs = append(userIDs, info.UserID)
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		g.SyncOneGroupApplicationsFunc(ctx, detail.Group.GroupID) //同步邀请需要申请的请求
		return g.syncGroupMembers(ctx, detail.Group.GroupID, userIDs...)
		//if utils.IsContain(g.loginUserID, userIDs) {
		//	return g.SyncGroups(ctx, detail.Group.GroupID)
		//} else {
		//	//登录用户不在列表中
		//	if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
		//		return err
		//	}
		//	return g.syncGroupMembers(ctx, detail.Group.GroupID, userIDs...)
		//}
		//进群通知
	case constant.MemberEnterNotification: // 1510
		var detail sdkws.MemberEnterTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.EntrantUser.UserID)
		//群解散通知
	case constant.GroupDismissedNotification: // 1511
		var detail sdkws.GroupDismissedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.DismissedNotification(ctx, &detail)
		//群成员禁言通知
	case constant.GroupMemberMutedNotification: // 1512
		var detail sdkws.GroupMemberMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.MutedUser.GroupID)
		//群成员禁言解除通知
	case constant.GroupMemberCancelMutedNotification: // 1513
		var detail sdkws.GroupMemberCancelMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.MutedUser.GroupID)
		//群禁言通知
	case constant.GroupMutedNotification: // 1514
		var detail sdkws.GroupMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
		//群禁言解除通知
	case constant.GroupCancelMutedNotification: // 1515
		var detail sdkws.GroupCancelMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
		//设置群成员信息
	case constant.GroupMemberInfoSetNotification: // 1516
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID)
		//设置管理员通知
	case constant.GroupMemberSetToAdminNotification: // 1517
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID)
		//管理员更改为普通成员通知
	case constant.GroupMemberSetToOrdinaryUserNotification: // 1518
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.syncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID)
	case 1519: // 1519  constant.GroupInfoSetAnnouncementNotification
		var detail sdkws.GroupInfoSetTips // sdkws.GroupInfoSetAnnouncementTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	case 1520: // 1520  constant.GroupInfoSetNameNotification
		var detail sdkws.GroupInfoSetTips // sdkws.GroupInfoSetNameTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	default:
		return fmt.Errorf("unknown tips type: %d", msg.ContentType)
	}
}

// DismissedNotification 群解散处理
func (g *Group) DismissedNotification(ctx context.Context, detail *sdkws.GroupDismissedTips) error {
	//删除会话
	g.DelGroupConversation(ctx, detail.Group.GroupID)
	//删除群成员
	if err := g.db.DeleteGroupAllMembers(ctx, detail.Group.GroupID); err != nil {
		return err
	}
	return g.db.DeleteGroup(ctx, detail.Group.GroupID)
}
