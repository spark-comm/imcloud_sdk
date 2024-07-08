package server_api

import (
	friendPb "github.com/spark-comm/spark-api/api/im_cloud/friend/v2"
	"testing"
)

func Test_GetAllFriendList(t *testing.T) {
	list, err := GetAllFriendList(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func Test_GetAllBlackList(t *testing.T) {
	list, err := GetAllBlackList(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
func Test_GetSendFriendApplication(t *testing.T) {
	list, err := GetSendFriendApplication(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
func Test_GetReceiveFriendApplication(t *testing.T) {
	list, err := GetReceiveFriendApplication(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func Test_BothFriendRequest(t *testing.T) {
	data, err := BothFriendRequest(getCtx(), "1319567527776256", "934495075176448")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func Test_SetFriendInfo(t *testing.T) {
	err := SetFriendInfo(getCtx(), &friendPb.SetFriendInfoReq{FromUserID: "922670631751680", ToUserID: "931422227402752", Remark: "你好30"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}
func Test_AddBlack(t *testing.T) {
	err := AddBlack(getCtx(), UserID, "931422227402752")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func Test_RemoveBlack(t *testing.T) {
	err := RemoveBlack(getCtx(), UserID, "931422227402752")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}
