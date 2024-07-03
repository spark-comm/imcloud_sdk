package module

import (
	"github.com/spark-comm/imcloud_sdk/pkg/constant"

	"github.com/OpenIMSDK/protocol/friend"
)

type TestFriendManager struct {
	*MetaManager
}

func (t *TestFriendManager) ImportFriends(ownerUserID string, friendUserIDs []string) error {
	req := &friend.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	}
	return t.postWithCtx(constant.ImportFriendListRouter, &req, nil)
}
