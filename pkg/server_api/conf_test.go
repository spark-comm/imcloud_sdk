package server_api

import (
	"context"

	"github.com/spark-comm/imcloud_sdk/pkg/ccontext"
	"github.com/spark-comm/imcloud_sdk/sdk_struct"
)

const (
	APIADDR = "http://127.0.0.1:9099"
	//APIADDR = "http://8.137.13.1:9099"
	WSADDR = "ws://8.137.13.1:10001"
	UserID = "922670631751680"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MzE0MjIyMjc0MDI3NTJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM3ODEyNDczODU2XCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIxMDY0NTk5LCJuYmYiOjE3MjA3MDQ1OTksImlhdCI6MTcyMDcwNDU5OX0.T1Ci6vpo0wLXcL-_8VeY86cPYxwBCTO3ZBRRKiJtL0KRf5hP43-AAr3IL1vOjhocXVVBjyRb78B5SouvWxsmPQ"
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
