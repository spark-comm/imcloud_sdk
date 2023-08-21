package friend

import (
	"context"
	"errors"
	"fmt"
	"github.com/imCloud/im/pkg/proto/sdkws"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

func (f *Friend) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if f.friendListener == nil {
		return errors.New("f.friendListener == nil")
	}
	if msg.SendTime < f.loginTime || f.loginTime == 0 {
		return errors.New("ignore notification")
	}
	switch msg.ContentType {
	case constant.FriendApplicationNotification:
		//好友请求增加
		tips := sdkws.FriendApplicationTips{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncApplicationByNotification(ctx, tips.FromToUserID)
	case constant.FriendApplicationApprovedNotification:
		//发起的好友请求被同意
		var tips sdkws.FriendApplicationApprovedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if err := f.syncApplicationByNotification(ctx, tips.FromToUserID); err != nil {
			return err
		}
		return f.syncFriendByNotification(ctx, f.loginUserID, tips.FromToUserID.ToUserID)
	case constant.FriendApplicationRejectedNotification:
		//发起的好友请求被拒绝
		var tips sdkws.FriendApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncApplicationByNotification(ctx, tips.FromToUserID)
	case constant.FriendAddedNotification:
		//新增好友通知
		var tips sdkws.FriendAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncFriendByNotification(ctx, f.loginUserID, tips.Friend.OwnerUserID)
	case constant.FriendDeletedNotification:
		//好友被删除通知
		var tips sdkws.FriendDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncFriendByNotification(ctx, f.loginUserID, tips.FromToUserID.FromUserID)
	case constant.FriendRemarkSetNotification:
		// 好友给设置备注
		var tips sdkws.FriendInfoChangedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.syncFriendByNotification(ctx, tips.FromToUserID.FromUserID, tips.FromToUserID.ToUserID)
		}
		return nil
	case constant.FriendInfoUpdatedNotification:
		//好友信息变更
		var tips sdkws.UserInfoUpdatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		//blChan <- true
		return f.syncFriendByNotification(ctx, f.loginUserID, tips.UserID)
	case constant.BlackAddedNotification:
		//被好友拉黑
		var tips sdkws.BlackAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	case constant.BlackDeletedNotification:
		//被好友移出黑名单
		var tips sdkws.BlackDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
}
