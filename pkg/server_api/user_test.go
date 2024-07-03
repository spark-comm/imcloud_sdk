package server_api

import (
	userPb "github.com/spark-comm/spark-api/api/im_cloud/user/v2"
	"testing"
)

func Test_GetSelfUserInfoFromSvr(t *testing.T) {
	data, err := GetSelfUserInfoFromSvr(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func Test_GetServerUserInfo(t *testing.T) {
	data, err := GetServerUserInfo(getCtx(), []string{UserID})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func Test_GetUsersOption(t *testing.T) {
	data, err := GetUsersOption(getCtx(), UserID, userPb.UserOption_globalRecvMsgOpt.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func Test_SetUserOption(t *testing.T) {
	err := SetUsersOption(getCtx(), UserID, userPb.UserOption_globalRecvMsgOpt.String(), 1)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_SearchUser(t *testing.T) {
	data, err := SearchUser(getCtx(), &userPb.SearchProfileReq{
		SearchValue: "18800001041",
		Type:        1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func Test_FindFullProfileByUserId(t *testing.T) {
	data, err := FindFullProfileByUserId(getCtx(), UserID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
