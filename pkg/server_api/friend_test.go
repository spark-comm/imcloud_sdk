package server_api

import "testing"

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
