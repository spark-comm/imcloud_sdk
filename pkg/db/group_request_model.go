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
	"fmt"

	"github.com/spark-comm/imcloud_sdk/pkg/db/model_struct"
	"github.com/spark-comm/imcloud_sdk/pkg/utils"
)

func (d *DataBase) InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	// 记录不存在，创建新记录
	return utils.Wrap(d.conn.WithContext(ctx).Create(groupRequest).Error, "InsertGroupRequest failed")
	//var localGroupReq model_struct.LocalGroupRequest
	//if d.conn.WithContext(ctx).Where("group_id = ? and user_id = ?", groupRequest.GroupID, groupRequest.UserID).First(&localGroupReq).RowsAffected == 0 {
	//	// 记录不存在，创建新记录
	//	return utils.Wrap(d.conn.WithContext(ctx).Create(groupRequest).Error, "InsertGroupRequest failed")
	//} else {
	//	// 记录存在，更新记录
	//	return utils.Wrap(d.conn.WithContext(ctx).Model(&localGroupReq).Updates(groupRequest).Error, "InsertGroupRequest failed")
	//}
}
func (d *DataBase) DeleteGroupRequest(ctx context.Context, groupID, userID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&model_struct.LocalGroupRequest{}).Error, "DeleteGroupRequest failed")
}
func (d *DataBase) UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupRequestList []model_struct.LocalGroupRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Order("create_time DESC").Find(&groupRequestList).Error, "")
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var transfer []*model_struct.LocalGroupRequest
	for _, v := range groupRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}
func (d *DataBase) GetUnProcessGroupRequestNum(ctx context.Context, userId string) (int64, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var num int64
	err := d.conn.WithContext(ctx).Model(&model_struct.LocalAdminGroupRequest{}).Where("handle_result = 0").Count(&num).Error
	return num, utils.Wrap(err, fmt.Sprintf("GetUnProcessGroupRequestNum er"))
}
