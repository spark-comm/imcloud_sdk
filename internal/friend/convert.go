package friend

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	friendPb "github.com/imCloud/api/friend/v1"
)

func ServerFriendRequestToLocalFriendRequest(info *friendPb.FriendRequests) *model_struct.LocalFriendRequest {
	return &model_struct.LocalFriendRequest{
		FromUserID:    info.FromUserID,
		FromNickname:  info.FromNickname,
		FromFaceURL:   info.FromFaceURL,
		FromGender:    info.FromGender,
		FromCode:      info.FromCode,
		FromPhone:     info.FromPhone,
		FromMessage:   info.FromMessage,
		ToUserID:      info.ToUserID,
		ToNickname:    info.ToNickname,
		ToFaceURL:     info.ToFaceURL,
		ToGender:      info.ToGender,
		ToMessage:     info.ToMessage,
		ToCode:        info.ToCode,
		TomPhone:      info.ToPhone,
		HandleResult:  info.HandleResult,
		ReqMsg:        info.ReqMsg,
		CreateAt:      info.CreateTime,
		HandlerUserID: info.HandlerUserID,
		HandleMsg:     info.HandleMsg,
		HandleTime:    info.HandleTime,
		Ex:            info.Ex,
		SortFlag:      getSortFlag(info.ToNickname, info.ToNickname),
	}
}
func ServerFriendToLocalFriend(info *friendPb.FriendInfo) *model_struct.LocalFriend {
	return &model_struct.LocalFriend{
		OwnerUserID:    info.OwnerUserID,
		FriendUserID:   info.FriendUserID,
		Remark:         info.Remark,
		CreateAt:       info.CreatedAt,
		AddSource:      info.AddSource,
		OperatorUserID: info.OperatorUserID,
		Nickname:       info.Nickname,
		FaceURL:        info.FaceURL,
		Ex:             info.Ex,
		Phone:          info.Phone,
		Code:           info.Code,
		Message:        info.Message,
		Email:          info.Email,
		Birth:          info.Birth,
		Gender:         info.Gender,
		SortFlag:       getSortFlag(info.Remark, info.Nickname),
		NotPeersFriend: info.NotPeersFriend,
		BackgroundURL:  info.BackgroundUrl,
	}
}

func ServerBlackToLocalBlack(info *friendPb.BlackList) *model_struct.LocalBlack {
	return &model_struct.LocalBlack{
		OwnerUserID:    info.OwnerUserID,
		BlackUserID:    info.BlackUserID,
		CreateAt:       info.CreatedAt,
		OperatorUserID: info.OwnerUserID,
		Nickname:       info.Nickname,
		FaceURL:        info.FaceURL,
		Gender:         info.Gender,
		Message:        info.Message,
		Code:           info.Code,
		Ex:             info.Ex,
		SortFlag:       getSortFlag(info.Nickname, info.Nickname),
	}
}

func getSortFlag(remake, nickname string) string {
	if remake != "" {
		return utils.GetChineseFirstLetter(remake)
	}
	return utils.GetChineseFirstLetter(nickname)
}
