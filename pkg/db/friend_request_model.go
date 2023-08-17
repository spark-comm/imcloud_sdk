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

//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/db/pg"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(friendRequest).Error, "InsertFriendRequest failed")
}

func (d *DataBase) DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&model_struct.LocalFriendRequest{}).Error, "DeleteFriendRequestBothUserID failed")
}

func (d *DataBase) UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(friendRequest).Select("*").Updates(*friendRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetRecvFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("to_user_id = ?", d.loginUserID).Order("create_at DESC").Find(&friendRequestList).Error, "GetRecvFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetRecvFriendApplication failed")
}

func (d *DataBase) GetSendFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id = ?", d.loginUserID).Order("create_at DESC").Find(&friendRequestList).Error, "GetSendFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSendFriendApplication failed")
}

// GetRecvFriendApplicationList 分页获取我收到的好友请求
func (d *DataBase) GetRecvFriendApplicationList(ctx context.Context, page *pg.Page) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	transfer := make([]*model_struct.LocalFriendRequest, 0)
	err := utils.Wrap(d.conn.WithContext(ctx).Where("to_user_id = ?", d.loginUserID).Scopes(pg.Operation(page)).Order("handle_result asc,create_at DESC").Find(&transfer).Error, "GetRecvFriendApplication failed")
	return transfer, utils.Wrap(err, "GetRecvFriendApplication failed")
}

// GetSendFriendApplicationList 分页获取我发送的好友请求
func (d *DataBase) GetSendFriendApplicationList(ctx context.Context, page *pg.Page) ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	transfer := make([]*model_struct.LocalFriendRequest, 0)
	err := utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id = ?", d.loginUserID).Scopes(pg.Operation(page)).Order("create_at DESC").Find(&transfer).Error, "GetSendFriendApplication failed")
	return transfer, utils.Wrap(err, "GetSendFriendApplication failed")
}
func (d *DataBase) GetFriendApplicationByBothID(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()

	var friendRequest model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Take(&friendRequest).Error, "GetFriendApplicationByBothID failed")

	return &friendRequest, utils.Wrap(err, "GetFriendApplicationByBothID failed")
}

// GetUnprocessedNum 获取未处理的好友申请数
func (d *DataBase) GetUnprocessedNum(ctx context.Context) (int64, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var count int64
	err := utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalFriendRequest{}).Where("to_user_id = ? and handle_result=0", d.loginUserID).Count(&count).Error, "GetUnprocessedNum failed")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DataBase) GetNotInListFriendInfo(ctx context.Context, cond, user string, userIDs []string, pageSize, pageNum int) ([]sdk_params_callback.SearchNotInGroupUserResp, int64, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	total := int64(0)
	result := []sdk_params_callback.SearchNotInGroupUserResp{}
	err := d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{}).Select([]string{
		"friend_user_id", "face_url", "nickname", "code", "phone", "gender", "remark",
	}).Where("owner_user_id = ?", user).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(userIDs) > 0 {
			db.Where("friend_user_id NOT IN (?)", userIDs)
		}
		if cond != "" {
			db.Where("nickname LIKE ? OR remark LIKE ?", "%"+cond+"%", "%"+cond+"%")
		}
		return db
	}).Count(&total).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&result).Error
	return result, total, err
}
