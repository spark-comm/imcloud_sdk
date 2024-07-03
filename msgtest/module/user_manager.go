package module

import (
	"fmt"
	"time"

	"github.com/spark-comm/imcloud_sdk/pkg/constant"

	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
)

type TestUserManager struct {
	*MetaManager
}

func (t *TestUserManager) GenUserIDs(num int) (userIDs []string) {
	for i := 0; i < num; i++ {
		userIDs = append(userIDs, fmt.Sprintf("testv3new_%d_%d", time.Now().UnixNano(), i))
	}
	return userIDs
}

func (t *TestUserManager) GenUserIDsWithPrefix(num int, prefix string) (userIDs []string) {
	for i := 0; i < num; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s_testv3new_%d_%d", prefix, time.Now().UnixNano(), i))
	}
	return userIDs
}

func (t *TestUserManager) RegisterUsers(userIDs ...string) error {
	var users []*sdkws.UserInfo
	for _, userID := range userIDs {
		users = append(users, &sdkws.UserInfo{UserID: userID, Nickname: userID})
	}
	return t.postWithCtx(constant.UserRegister, &userPB.UserRegisterReq{
		Secret: t.secret,
		Users:  users,
	}, nil)
}

func (t *TestUserManager) GetToken(userID string, platformID int32) (string, error) {
	return t.getToken(userID, platformID)
}
