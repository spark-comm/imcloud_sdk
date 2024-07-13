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
	//30
	UserID = "931422227402752"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MzE0MjIyMjc0MDI3NTJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM3ODEyNDczODU2XCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIzNDczMDgxLCJuYmYiOjE3MjA4ODEwODEsImlhdCI6MTcyMDg4MTA4MX0.UgZTLwkTTq9mLziSUm1patGSxeJgxarrJA1Z2FGgGicOXTMss7eD9PAEm-nYx7zk0vQMy_f65Pnm7WPQWFg8Ng"
	// 26
	//UserID = "922670631751680"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MjI2NzA2MzE3NTE2ODBcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM2NjEyOTAyOTEyXCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIzNDcxMzU0LCJuYmYiOjE3MjA4NzkzNTQsImlhdCI6MTcyMDg3OTM1NH0.32IEp5vo0Dxpu4SksPsn44I8FM99k7ydOyrr1yWzTV7LoEkOPHi8f4RsiVNDeglsAe4vEdVTKq2jipJWRAvXTg"
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
