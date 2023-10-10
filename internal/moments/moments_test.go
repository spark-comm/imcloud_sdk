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
	//APIADDR = "http://47.109.111.139:4368"
	APIADDR = "http://127.0.0.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	UserID  = "55236290680983552"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTIzNjI5MDY4MDk4MzU1MlwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwicm9sZVwiOlwiVVNFUlwifSIsImV4cCI6MTY5NzI2Mzg0MywibmJmIjoxNjk2OTAzODQzLCJpYXQiOjE2OTY5MDM4NDN9.irii4-M5EVOiRN0f_YQ-bAC0K4_c4OlNa_FVI8F40z0UbI1Gd9sYWcdSyialZHSmunICou5_3Y4T6z4MoIyJlQ"
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

func TestSyncNewMomentsFromSvr(t *testing.T) {
	err := moments.SyncNewMomentsFromSvr(getCtx())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestSyncHistoryMomentsFromSvr(t *testing.T) {
	err := moments.SyncHistoryMomentsFromSvr(getCtx())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}
