package msgtest

import (
	"context"
	"testing"

	"github.com/OpenIMSDK/tools/log"
	"github.com/brian-god/imcloud_sdk/msgtest/module"
	"github.com/brian-god/imcloud_sdk/msgtest/sdk_user_simulator"
	"github.com/brian-god/imcloud_sdk/pkg/ccontext"
	"github.com/brian-god/imcloud_sdk/pkg/utils"
)

func Test_SimulateMultiOnline(t *testing.T) {
	ctx := ccontext.WithOperationID(context.Background(), "TEST_ROOT")
	userIDList := []string{"1", "2"}
	metaManager := module.NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	userManager := metaManager.NewUserManager()
	serverTime, err := metaManager.GetServerTime()
	if err != nil {
		t.Fatal(err)
	}
	offset := serverTime - utils.GetCurrentTimestampByMill()
	sdk_user_simulator.SetServerTimeOffset(offset)
	for _, userID := range userIDList {
		token, err := userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(ctx, "get token failed, err: %v", err, "userID", userID)
			continue
		}
		err = sdk_user_simulator.InitSDKAndLogin(userID, token)
		if err != nil {
			log.ZError(ctx, "login failed, err: %v", err, "userID", userID)
		} else {
			log.ZDebug(ctx, "login success, userID: %v", "userID", userID)
		}
	}

}
