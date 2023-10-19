package testv2

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"testing"
)

func TestLikeMoments(t *testing.T) {
	err := open_im_sdk.UserForSDK.Moments().Like(ctx, &sdk_params_callback.LikeRequest{
		MomentId:     "484927858544640",
		UserNickname: "10000622",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}
func TestCommentMoments(t *testing.T) {
	err := open_im_sdk.UserForSDK.Moments().Comment(ctx, &sdk_params_callback.CommentRequest{
		MomentId:       "484927858544640",
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
