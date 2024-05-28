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
	APIADDR = "http://47.108.68.161:9099"
	WSADDR  = "ws://47.108.68.161:10001"
	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	//UserID = "14743920172863488"
	//token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxNDc0MzkyMDE3Mjg2MzQ4OFwiLFwiY2VudGVyX3VzZXJfaWRcIjpcIjE0NzEyOTIyMTU2NTY4NTc2XCIsXCJwbGF0Zm9ybVwiOlwiV2luZG93c1wiLFwidGVuYW50SWRcIjpcIjUyMjQ1MDIzMjE1MjA2NFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI1MjI0NTAyMzIxNTIwNjRcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxNjg3NTAzNSwibmJmIjoxNzE2NTE1MDM1LCJpYXQiOjE3MTY1MTUwMzV9.RLQvUupF-7xRdzL73xG8MVwwrHCv4Ywo90HUsaoOdRdXDS9B2D6lhCap85I61pew9UtnVMbslc5-Xmrhhi7yWQ"
	UserID = "133374647734272"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxMzMzNzQ2NDc3MzQyNzJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxNDcxMzM4OTkzODkwNTA4OFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInRlbmFudElkXCI6XCI1MjI0NTAyMzIxNTIwNjRcIixcInNlcnZlcl9jb2RlXCI6XCJcIixcInJvbGVcIjpcIlVzZXJcIixcInNjb3BlXCI6XCJcIixcIm5vZGVJZFwiOlwiNTIyNDUwMjMyMTUyMDY0XCIsXCJvcHRpb25zXCI6bnVsbH0iLCJleHAiOjE3MTcxNTQ1MDMsIm5iZiI6MTcxNjc5NDUwMywiaWF0IjoxNzE2Nzk0NTAzfQ.4RqYyraSd8EJV4B2d-3GhI0YL9nxWJqnUWbuOJj-EQnHuQ_eX3S0ItkCU6xmxm1xfpeegqDIyQoeqmWa4SZ0hQ"
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
