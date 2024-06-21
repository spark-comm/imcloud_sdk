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
	//UserID = "14743920172863488"
	//token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxNDc0MzkyMDE3Mjg2MzQ4OFwiLFwiY2VudGVyX3VzZXJfaWRcIjpcIjE0NzEyOTIyMTU2NTY4NTc2XCIsXCJwbGF0Zm9ybVwiOlwiV2luZG93c1wiLFwidGVuYW50SWRcIjpcIjUyMjQ1MDIzMjE1MjA2NFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI1MjI0NTAyMzIxNTIwNjRcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxNjg3NTAzNSwibmJmIjoxNzE2NTE1MDM1LCJpYXQiOjE3MTY1MTUwMzV9.RLQvUupF-7xRdzL73xG8MVwwrHCv4Ywo90HUsaoOdRdXDS9B2D6lhCap85I61pew9UtnVMbslc5-Xmrhhi7yWQ"
	UserID = "922670631751680"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MjI2NzA2MzE3NTE2ODBcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM2NjEyOTAyOTEyXCIsXCJwbGF0Zm9ybVwiOlwiV2luZG93c1wiLFwidGVuYW50SWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI5MTEzNTU3NjI3MDg0ODBcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxOTI1NTI0MiwibmJmIjoxNzE4ODk1MjQyLCJpYXQiOjE3MTg4OTUyNDJ9.7C4YVl8JFLAQsHRf47RASQM9rEZJZCR2uSCFPJfMoNQufk-UefQwsxBTo1fdJpD1RosLsKMxq_KvQSaaGrbUHg"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 4
	cf.IsExternalExtensions = true
	cf.PlatformID = 1
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	cf.Language = "en"
	return cf
}
