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
	//APIADDR = "http://8.137.13.1:9099"
	//WSADDR  = "ws://8.137.13.1:10001"
	APIADDR = "http://0.0.0.0:9099"
	WSADDR  = "ws://0.0.0.0:10001"
	UserID  = "49394919499894784"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI0OTM4MzY3NTU5NDI4MDk2MFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlwifSIsImV4cCI6MTY5MTc1MjYzNiwibmJmIjoxNjkxMzkyNjM2LCJpYXQiOjE2OTEzOTI2MzZ9.FyBCcmzThoD9RFZCu9fwr_W8pvQ5VhMa-lisKOo5fv_gNqhGAFtpK_DnmOBgWC47JVIRvu2d6mksytGS-LykFw"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "../"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = 1
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
