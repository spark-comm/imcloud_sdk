// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package group

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	groupv1 "github.com/imCloud/api/group/v1"
)

const (
	IsNotComplete = 2
)

func ServerGroupToLocalGroup(info *groupv1.GroupInfo) *model_struct.LocalGroup {
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

func ServerGroupMemberToLocalGroupMember(info *groupv1.MembersInfo) *model_struct.LocalGroupMember {
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

func ServerGroupRequestToLocalGroupRequest(info *groupv1.GroupRequestInfo) *model_struct.LocalGroupRequest {
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

func ServerBaseGroupToLocalGroup(info *groupv1.BaseGroupInfo) *model_struct.LocalGroup {
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

func ServerBaseGroupMemberToLocalGroupMember(info *groupv1.BaseGroupMemberInfo) *model_struct.LocalGroupMember {
	return &model_struct.LocalGroupMember{
		GroupID:   info.GroupID,
		UserID:    info.UserID,
		RoleLevel: info.RoleLevel,
		Nickname:  info.NickName,
		SortFlag: func() string {
			return utils.GetChineseFirstLetter(info.NickName)
		}(),
		GroupUserName: info.GroupUserName,
		FaceURL:       info.FaceURL,
		MuteEndTime:   info.MuteEndTime,
		Code:          info.Code,
		UpdatedAt:     info.UpdateAt,
		BackgroundURL: info.BackgroundUrl,
	}
}

func ServerGroupRequestToLocalAdminGroupRequest(info *groupv1.GroupRequestInfo) *model_struct.LocalAdminGroupRequest {
	return &model_struct.LocalAdminGroupRequest{
		LocalGroupRequest: *ServerGroupRequestToLocalGroupRequest(info),
	}
}
