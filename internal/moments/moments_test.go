package moments

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"testing"
)

const (
	//预生产
	APIADDR = "http://47.109.111.139:4368"
	//APIADDR = "http://127.0.0.1:4368"
	WSADDR = "ws://8.137.13.1:10001"
	UserID = "48471067332608"
	token  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiI0ODQ3MTA2NzMzMjYwOCIsIlBsYXRmb3JtIjoiSU9TIiwiZXhwIjoxNzEwMDU1MTM1LCJuYmYiOjE2OTQ1MDI4MzUsImlhdCI6MTY5NDUwMzEzNX0.Jp9vQkD2i76SM0KoYmTNWEfCPC6XyMpW_WefZpqpWwI"
	//UserID  = "277"
	//token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIyNzciLCJQbGF0Zm9ybSI6IklPUyIsImV4cCI6MTcwOTk1NTE5NSwibmJmIjoxNjk0NDAyODk1LCJpYXQiOjE2OTQ0MDMxOTV9.uqAq1j7j6jUddfP5DAP3yiKGnprDov-KSHzAQgVQLAQ"
)

func getCtx() context.Context {
	var gl = ccontext.GlobalConfig{}
	gl.ApiAddr = APIADDR
	gl.WsAddr = WSADDR
	gl.DataDir = "../"
	gl.LogLevel = 6
	gl.IsExternalExtensions = true
	gl.Token = token
	gl.UserID = UserID
	return ccontext.WithInfo(context.Background(), &gl)
}

func getDb() db_interface.DataBase {
	sqliteConn, err := db.NewDataBase(getCtx(), UserID, "../")
	if err != nil {
		panic(err)
	}

	return sqliteConn
}

var moments = Moments{loginUserID: UserID, db: getDb()}

func TestGetMomentsFromSvr(t *testing.T) {
	svr, i, err := moments.getMomentsFromSvr(getCtx(), 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success", svr, i)
}
