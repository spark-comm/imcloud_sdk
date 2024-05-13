package friend

import (
	"context"
	"errors"
	"fmt"
	"github.com/imCloud/api/common/notice"
	"github.com/imCloud/im/pkg/common/log"
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
		//好友请求通知
		tips := notice.FriendApplicationTips{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		log.ZInfo(ctx, "收到好友请求通知，开始同步数据,tips data:%+v", tips)
		return f.syncApplication(ctx, tips.FromToUserID)
	case constant.FriendApplicationApprovedNotification:
		//发起的好友请求被同意
		var tips sdkws.FriendApplicationApprovedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if err := f.syncApplicationByNotification(ctx, tips.FromToUserID); err != nil {
			return err
		}
		return f.syncFriendByNotification(ctx, tips.FromToUserID.ToUserID)
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
		//如果是后台处理这里的好友请求需要重新同步，否则移动端还能处理
		//err := f.SyncUntreatedFriendReceiveFriendApplication(ctx)
		//if err != nil {
		//	log.Println(fmt.Sprintf("SyncUntreatedFriendReceiveFriendApplication err:%+v", err))
		//}
		return f.syncFriendByNotification(ctx, tips.Friend.OwnerUserID)
	case constant.FriendDeletedNotification:
		//好友被删除通知
		var tips sdkws.FriendDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.SyncDelFriend(ctx, tips.FromToUserID.FromUserID)
	case constant.FriendRemarkSetNotification:
		// 好友给设置备注
		var tips sdkws.FriendInfoChangedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.syncFriendByNotification(ctx, tips.FromToUserID.ToUserID)
		}
		return nil
	case constant.FriendInfoUpdatedNotification:
		//好友信息变更
		var tips sdkws.UserInfoUpdatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		//blChan <- true
		return f.syncFriendByNotification(ctx, tips.UserID)
	case constant.BlackAddedNotification:
		//被好友拉黑
		var tips sdkws.BlackAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		} else {
			//被对方拉黑
			return f.SyncBlackList(ctx)
		}
	case constant.BlackDeletedNotification:
		//被好友移出黑名单
		var tips sdkws.BlackDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		} else {
			//被对方拉黑
			return f.SyncBlackList(ctx)
		}
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
}
