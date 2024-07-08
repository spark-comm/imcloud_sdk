package server_api

import "testing"

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
