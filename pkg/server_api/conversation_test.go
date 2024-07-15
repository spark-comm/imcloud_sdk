package server_api

import "testing"

func Test_ClearConversationMsgFromSvr(t *testing.T) {
	err := ClearConversationMsgFromSvr(getCtx(), UserID, "sg_1334607223984128", true)
	if err != nil {
		t.Error(err)
	}
}
