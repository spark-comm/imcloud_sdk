package server_api

import (
	v2 "github.com/spark-comm/spark-api/api/im_cloud/group/v2"
	"testing"
)

func Test_GetSpecifiedGroupsInfo(t *testing.T) {
	data, err := GetSpecifiedGroupsInfo(getCtx(), []string{"237875325046784"})
	if err != nil {
		t.Error(err)

	}
	t.Log(data) //t.Log(data["data"].([]interface{})[0].(map[string]interface{})["group_id"])}
}

func Test_GetServerJoinGroup(t *testing.T) {
	list, err := GetServerJoinGroup(getCtx(), UserID)
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func Test_GetGroupsInfo(t *testing.T) {
	list, err := GetGroupsInfo(getCtx(), "237875325046784")
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func Test_SearchGroupByCode(t *testing.T) {
	data, err := SearchGroupByCode(getCtx(), UserID, "h6db2y")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func Test_GetServerSelfGroupApplication(t *testing.T) {
	data, err := GetServerSelfGroupApplication(getCtx(), UserID)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}
func Test_GetServerAdminGroupApplicationList(t *testing.T) {
	data, err := GetServerAdminGroupApplicationList(getCtx(), UserID)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func Test_GetServerGroupMembers(t *testing.T) {
	data, err := GetServerGroupMembers(getCtx(), "237875325046784")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func Test_GetGroupAbstractInfo(t *testing.T) {
	data, err := GetGroupAbstractInfo(getCtx(), "237875325046784")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func Test_CreateGroup(t *testing.T) {
	data, err := CreateGroup(getCtx(), &v2.CrateGroupReq{
		GroupName:  "sdk创建",
		GroupType:  2,
		MemberList: []string{"1331954574168064", "1332177157492736", "1332053626851328", "922670631751680"},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func Test_DismissGroup(t *testing.T) {
	err := DismissGroup(getCtx(), "209711102169088", UserID)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetGroupSwitchInfo(t *testing.T) {
	err := SetGroupSwitchInfo(getCtx(), "586691660222464", UserID, v2.GroupSwitchOption_needVerification.String(), 1)
	if err != nil {
		t.Error(err)
	}
}
func Test_JoinGroup(t *testing.T) {
	err := JoinGroup(getCtx(), UserID, "586691660222464", "我想进群", 1)
	if err != nil {
		t.Error(err)
	}
}
