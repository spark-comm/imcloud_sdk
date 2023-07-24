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

package user

import (
	imUserPb "github.com/imCloud/api/user/v1"
	"open_im_sdk/pkg/db/model_struct"
)

func ServerUserToLocalUser(user *imUserPb.ProfileReply) *model_struct.LocalUser {
	return &model_struct.LocalUser{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		FaceURL:          user.FaceURL,
		CreateTime:       user.CreatedAt.Seconds,
		AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalMsgReceive,
		//AttachedInfo: user.AttachedInfo,
	}
}
