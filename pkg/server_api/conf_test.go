package server_api

import (
	"context"

	"github.com/spark-comm/imcloud_sdk/pkg/ccontext"
	"github.com/spark-comm/imcloud_sdk/sdk_struct"
)

const (
	//APIADDR = "http://127.0.0.1:9099"
	APIADDR = "http://8.137.13.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	UserID  = "922670631751680"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MjI2NzA2MzE3NTE2ODBcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM2NjEyOTAyOTEyXCIsXCJwbGF0Zm9ybVwiOlwiV2luZG93c1wiLFwidGVuYW50SWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI5MTEzNTU3NjI3MDg0ODBcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcyMDQ2MzQxMCwibmJmIjoxNzIwMTAzNDEwLCJpYXQiOjE3MjAxMDM0MTB9.OenbkK8XPE_aYhM8Toi4Q0jhIEhmufCtUt8ek9pjf1_vpYWnlddSU3kiyZ3fFr-czpk0SmgX9CdOJjxle3eB8Q"
)

func getCtx() context.Context {
	info := &ccontext.GlobalConfig{
		UserID: UserID,
		Token:  token,
		IMConfig: sdk_struct.IMConfig{
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
			PlatformID: 3,
			DataDir:    "./",
			LogLevel:   1,
			Language:   "en",
		},
	}
	ctx := ccontext.WithInfo(context.Background(), info)
	return ctx
}
