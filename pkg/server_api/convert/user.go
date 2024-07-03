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

package convert

import (
	"encoding/json"

	"github.com/brian-god/imcloud_sdk/pkg/db/model_struct"
	usermodel "github.com/spark-comm/spark-api/api/common/model/user/v2"
)

func ServerUserToLocalUser(user *usermodel.UserProfile) (*model_struct.LocalUser, error) {
	op, err := json.Marshal(user.Options)
	if err != nil {
		return nil, err
	}
	return &model_struct.LocalUser{
		UserID:     user.UserId,
		FaceURL:    user.FaceURL,
		Nickname:   user.Nickname,
		Message:    user.Message,
		Code:       user.Code,
		Phone:      user.Phone,
		Email:      user.Email,
		Birth:      user.Birth,
		Gender:     user.Gender,
		CreateTime: user.CreatedAt.Seconds,
		ShareCode:  user.ShareCode,
		Options: func() string {
			if op == nil {
				return ""
			}
			return string(op)
		}(),
		Account:     user.Account,
		AccountType: user.AccountType,
	}, nil
}
