package server_api

import (
	"context"
	"github.com/brian-god/imcloud_sdk/pkg/ccontext"
	"github.com/brian-god/imcloud_sdk/sdk_struct"
)

const (
	APIADDR = "http://127.0.0.1:9099"
	//APIADDR = "http://8.137.13.1:9099"
	WSADDR = "ws://8.137.13.1:10001"
	UserID = "931422227402752"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MzE0MjIyMjc0MDI3NTJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM3ODEyNDczODU2XCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIwMjc3ODkxLCJuYmYiOjE3MTk5MTc4OTEsImlhdCI6MTcxOTkxNzg5MX0.CCEAE16Tuk-MXM2oPR7EqVUqM6P6um8gWzUC0HNpNUuZ-92tzqCZU3ix-0ciKbXPQy4UORC4vqyELprWOmUmfg"
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
