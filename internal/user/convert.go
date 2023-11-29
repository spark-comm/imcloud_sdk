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
	"encoding/json"
	"open_im_sdk/pkg/db/model_struct"

	imUserPb "github.com/imCloud/api/user/v1"
	"github.com/jinzhu/copier"
)

func ServerUserToLocalUser(user *imUserPb.ProfileReply) (*model_struct.LocalUser, error) {
	loginUser := model_struct.LocalUser{}
	if err := copier.Copy(&loginUser, user); err != nil {
		return nil, err
	}
	loginUser.UserID = user.UserId
	b, err := json.Marshal(user.Options)
	if err != nil {
		return nil, err
	}
	loginUser.Options = string(b)
	loginUser.AccountType = user.AccountType
	loginUser.Account = user.Account
	return &loginUser, nil
}
