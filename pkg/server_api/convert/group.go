package convert

import (
	groupmodel "github.com/miliao_apis/api/common/model/group/v2"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
)

const (
	IsNotComplete = 2
)

func ServerGroupToLocalGroup(info *groupmodel.GroupInfo) *model_struct.LocalGroup {
	return &model_struct.LocalGroup{
		GroupID:                info.GroupID,
		GroupName:              info.GroupName,
		Notification:           info.Notification,
		Introduction:           info.Introduction,
		FaceURL:                info.FaceURL,
		CreateTime:             info.CreateTime,
		Status:                 info.Status,
		CreatorUserID:          info.CreatorUserID,
		GroupType:              info.GroupType,
		OwnerUserID:            info.OwnerUserID,
		MemberCount:            int32(info.MemberCount),
		Ex:                     info.Ex,
		NeedVerification:       info.NeedVerification,
		LookMemberInfo:         info.LookMemberInfo,
		ApplyMemberFriend:      info.ApplyMemberFriend,
		NotificationUpdateTime: info.NotificationUpdateTime,
		NotificationUserID:     info.NotificationUserID,
		AttachedInfo:           info.AttachedInfo,
		Code:                   info.Code,
		MaxMemberCount:         info.MaxMemberCount,
		OnlyManageUpdateName:   info.OnlyManageUpdateName,
		IsReal:                 info.IsReal,
		IsOpen:                 uint(info.IsOpen),
		AllowPrivateChat:       uint(info.AllowPrivateChat),
		IsComplete:             info.IsComplete,
		UpdatedAt:              info.UpdateAt,
	}
}

func ServerGroupMemberToLocalGroupMember(info *groupmodel.MemberInfo) *model_struct.LocalGroupMember {
	return &model_struct.LocalGroupMember{
		GroupID:   info.GroupID,
		UserID:    info.UserID,
		RoleLevel: info.RoleLevel,
		JoinTime:  info.JoinTime,
		Nickname:  info.Nickname,
		SortFlag: func() string {
			return utils.GetChineseFirstLetter(info.Nickname)
		}(),
		GroupUserName:  info.GroupUserName,
		FaceURL:        info.FaceUrl,
		AttachedInfo:   info.AttachedInfo,
		JoinSource:     info.JoinSource,
		OperatorUserID: info.OperatorUserID,
		Ex:             info.Ex,
		MuteEndTime:    info.MuteEndTime,
		Code:           info.Code,
		InviterUserID:  info.InviterUserID,
		BackgroundURL:  info.BackgroundUrl,
		UpdatedAt:      info.UpdateAt,
	}
}

func ServerGroupRequestToLocalGroupRequest(info *groupmodel.GroupRequestInfo) *model_struct.LocalGroupRequest {
	return &model_struct.LocalGroupRequest{
		GroupID:       info.GroupID,
		GroupName:     info.GroupName,
		Notification:  info.Notification,
		Introduction:  info.Introduction,
		GroupFaceURL:  info.GroupFaceURL,
		CreateTime:    info.CreateTime,
		Status:        info.Status,
		CreatorUserID: info.CreatorUserID,
		GroupType:     info.GroupType,
		OwnerUserID:   info.OwnerUserID,
		MemberCount:   info.MemberCount,
		GroupCode:     info.GroupCode,

		UserID:      info.UserID,
		Nickname:    info.Nickname,
		UserFaceURL: info.UserFaceURL,
		Gender:      info.Gender,
		Code:        info.Code,

		HandleResult:  info.HandleResult,
		ReqMsg:        info.ReqMsg,
		HandledMsg:    info.HandledMsg,
		ReqTime:       info.ReqTime,
		HandleUserID:  info.HandleUserID,
		HandledTime:   info.HandledTime,
		JoinSource:    info.JoinSource,
		InviterUserID: info.InviterUserID,
	}
}

func ServerBaseGroupToLocalGroup(info *groupmodel.BaseGroupInfo) *model_struct.LocalGroup {
	return &model_struct.LocalGroup{
		GroupID:       info.GroupID,
		GroupName:     info.NickName,
		FaceURL:       info.FaceURL,
		Status:        info.Status,
		GroupType:     int32(info.GroupType),
		MemberCount:   int32(info.MemberCount),
		Code:          info.Code,
		CreatorUserID: info.CreatorUserID,
		OwnerUserID:   info.OwnerUserID,
		IsComplete:    IsNotComplete,
		UpdatedAt:     info.UpdateAt,
	}
}

func ServerGroupRequestToLocalAdminGroupRequest(info *groupmodel.GroupRequestInfo) *model_struct.LocalAdminGroupRequest {
	return &model_struct.LocalAdminGroupRequest{
		LocalGroupRequest: *ServerGroupRequestToLocalGroupRequest(info),
	}
}
