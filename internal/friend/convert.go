package friend

import (
	"open_im_sdk/pkg/db/model_struct"

	friendPb "github.com/imCloud/api/friend/v1"
)

func ServerFriendRequestToLocalFriendRequest(info *friendPb.FriendRequests) *model_struct.LocalFriendRequest {
	return &model_struct.LocalFriendRequest{
		FromUserID:    info.FromUserId,
		FromNickname:  info.FromNickname,
		FromFaceURL:   info.FromFaceUrl,
		FromGender:    info.FromGender,
		FromCode:      info.FromCode,
		FromPhone:     info.FromPhone,
		FromMessage:   info.FromMessage,
		ToUserID:      info.ToUserId,
		ToNickname:    info.ToNickname,
		ToFaceURL:     info.ToFaceUrl,
		ToGender:      info.ToGender,
		ToMessage:     info.ToMessage,
		ToCode:        info.ToCode,
		TomPhone:      info.ToPhone,
		HandleResult:  info.HandleResult,
		ReqMsg:        info.ReqMsg,
		CreateTime:    info.CreateTime,
		HandlerUserID: info.HandlerUserId,
		HandleMsg:     info.HandleMsg,
		HandleTime:    info.HandleTime,
		Ex:            info.Ex,
	}
}
func ServerFriendToLocalFriend(info *friendPb.ListFriendForSdkFriendInfo) *model_struct.LocalFriend {
	return &model_struct.LocalFriend{
		OwnerUserID:    info.OwnerUserId,
		FriendUserID:   info.FriendUserId,
		Remark:         info.Remark,
		CreateTime:     info.CreateTime,
		AddSource:      info.AddSource,
		OperatorUserID: info.OperatorUserId,
		Nickname:       info.Nickname,
		FaceURL:        info.FaceUrl,
		Ex:             info.Ex,
		Phone:          info.Phone,
		Code:           info.Code,
		Message:        info.Message,
		Email:          info.Email,
		Birth:          info.Birth,
		Gender:         info.Gender,
	}
}

func ServerBlackToLocalBlack(info *friendPb.BlackList) *model_struct.LocalBlack {
	return &model_struct.LocalBlack{
		OwnerUserID:    info.OwnerUserId,
		BlockUserID:    info.BlackUserId,
		CreateTime:     info.CreatedAt,
		OperatorUserID: info.OwnerUserId,
		Nickname:       info.Nickname,
		FaceURL:        info.FaceUrl,
		Gender:         info.Gender,
		Message:        info.Message,
		Code:           info.Code,
		Ex:             info.Ex,
	}
}
