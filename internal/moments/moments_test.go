package moments

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/sdk_params_callback"
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

func TestPublishMoments(t *testing.T) {
	err := moments.Publish(getCtx(), &sdk_params_callback.PublishRequest{
		UserName:  "asfasd",
		Avatar:    "https://dasfa.com/a.png",
		Content:   "测认识",
		Images:    []string{"dafafa", "111"},
		VideoUrl:  "https://v.png",
		VideoImg:  "https://ii.png",
		Location:  "afdsfas",
		Longitude: 1.11,
		Latitude:  2.22,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestDeleteMoments(t *testing.T) {
	err := moments.Delete(getCtx(), &sdk_params_callback.DeleteRequest{
		MomentId: "29685941538816",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestCommentMoments(t *testing.T) {
	err := moments.Comment(getCtx(), &sdk_params_callback.CommentRequest{
		MomentId:       "31723265986560",
		Type:           1,
		Avatar:         "dsafda",
		Nickname:       "sdfadfa",
		Content:        "afasdfas",
		SourceUserId:   "11111",
		SourceAvatar:   "1111",
		SourceNickname: "11111",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestLikeMoments(t *testing.T) {
	err := moments.Like(getCtx(), &sdk_params_callback.LikeRequest{
		MomentId:     "31723265986560",
		UserNickname: "dasfads",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestUnLikeMoments(t *testing.T) {
	err := moments.UnLike(getCtx(), &sdk_params_callback.UnlikeRequest{
		MomentId: "31723265986560",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestGetMomentsList(t *testing.T) {
	list, err := moments.GetMomentsList(getCtx(), &sdk_params_callback.V2ListRequest{
		IsSelf: false,
		Page:   1,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success", list)
}
