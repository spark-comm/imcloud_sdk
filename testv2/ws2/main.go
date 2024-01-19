package main

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/sdk_struct"
	"time"
)

func main() {
	fmt.Println("------------------------>>>>>>>>>>>>>>>>>>> test v2 init funcation <<<<<<<<<<<<<<<<<<<------------------------")
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	configData, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	isInit := open_im_sdk.InitSDK(listner, "test", string(configData))
	if !isInit {
		panic("init sdk failed")
	}
	ctx := open_im_sdk.UserForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID")
	//token, err := GetUserToken(ctx, UserID)
	//if err != nil {
	//	panic(err)
	//}
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	open_im_sdk.UserForSDK.SetListenerForService(&onListenerForService{ctx: ctx})
	open_im_sdk.UserForSDK.SetConversationListener(&onConversationListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetGroupListener(&onGroupListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetAdvancedMsgListener(&onAdvancedMsgListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetFriendListener(&onFriendListener{ctx: ctx})
	for true {
		time.Sleep(time.Second * 60)
	}
}

const (
	APIADDR = "http://8.137.13.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	UserID = "102803301208064"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxMDI4MDMzMDEyMDgwNjRcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxMDI4MDMzMDEyMDgwNjRcIixcInBsYXRmb3JtXCI6XCJJT1NcIixcInRlbmFudElkXCI6XCIxMjIyOTgxNzg3MzYxMjhcIixcInNlcnZlcl9jb2RlXCI6XCJcIixcInJvbGVcIjpcIlVzZXJcIixcInNjb3BlXCI6XCJcIixcIm5vZGVJZFwiOlwiMTIyMjk4MTc4NzM2MTI4XCIsXCJvcHRpb25zXCI6bnVsbH0iLCJleHAiOjE3MDYwMTYwNTgsIm5iZiI6MTcwNTY1NjA1OCwiaWF0IjoxNzA1NjU2MDU4fQ.Hr_qul4novWcKJcJb4U_05_oRnZxoSGYmXVUTabPzPvSg1rClUfdvtUz9818SfAx3f8asJyLHY6wiIhuERnBKg"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 2
	cf.IsExternalExtensions = true
	cf.PlatformID = 3
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
