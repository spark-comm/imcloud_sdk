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
	UserID = "55122332682817536"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjM2NTY2NjgyNDE5MlwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE2OTk0NDI4MzksIm5iZiI6MTY5ODcyMjgzOSwiaWF0IjoxNjk4NzIyODM5fQ.iGmBGdYtMI1E4Tq6wKjZTczhVYqpxQOLaaVT2XbyEnUrs_6rRfan3lURXKaXBOkww4gE4Sk6QyFf19DEr99cTw"
	token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjMzMjY4MjgxNzUzNlwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwidG9rZW5cIjpcIlwiLFwidGVuYW50SWRcIjpcIlwiLFwiZXhwaXJlX3RpbWVfc2Vjb25kc1wiOjAsXCJyb2xlXCI6XCJVU0VSXCJ9IiwiZXhwIjoxNzAwODQ4NjI3LCJuYmYiOjE3MDA0ODg2MjcsImlhdCI6MTcwMDQ4ODYyN30.uSZv4DGyIojJmPny3jIUAPJjX6991QFWRuBW9NLue5Fzp5uajq0g8I7FEVq7n9Em-2tcNSSwJFb5rpm5SJgRAQ"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 4
	cf.IsExternalExtensions = true
	cf.PlatformID = 3
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	cf.Language = "en"
	return cf
}
