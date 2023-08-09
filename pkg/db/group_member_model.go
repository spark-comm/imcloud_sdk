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
	"gorm.io/gorm"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMember model_struct.LocalGroupMember
	return &groupMember, utils.Wrap(d.conn.WithContext(ctx).Where("group_id = ? AND user_id = ?",
		groupID, userID).Take(&groupMember).Error, "GetGroupMemberInfoByGroupIDUserID failed")
}

func (d *DataBase) GetAllGroupMemberList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	return groupMemberList, utils.Wrap(d.conn.WithContext(ctx).Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}
func (d *DataBase) GetAllGroupMemberUserIDList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	return groupMemberList, utils.Wrap(d.conn.WithContext(ctx).Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}

func (d *DataBase) GetGroupMemberCount(ctx context.Context, groupID string) (int32, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var count int64
	err := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).Where("group_id = ? ", groupID).Count(&count).Error
	return int32(count), utils.Wrap(err, "GetGroupMemberCount failed")
}

func (d *DataBase) GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? And user_id IN ? ", groupID, userIDList).Find(&groupMemberList).Error
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}
func (d *DataBase) GetGroupAdminID(ctx context.Context, groupID string) ([]string, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var adminIDList []string
	return adminIDList, utils.Wrap(d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).Select("user_id").Where("group_id = ? And role_level = ?", groupID, constant.GroupAdmin).Find(&adminIDList).Error, "")
}

func (d *DataBase) GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? ", groupID).Find(&groupMemberList).Error
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}
func (d *DataBase) GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	var err error
	switch filter {
	case constant.GroupFilterAll:
		err = d.conn.WithContext(ctx).Where("group_id = ?", groupID).Order("role_level DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
		//Order("CASE role_level WHEN 2 THEN 1 WHEN 3 TEN 2 WHEN 1 THEN 3 ELSE 4 END").
	case constant.GroupFilterOwner:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupOwner).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupAdmin).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupOrdinaryUsers).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterAdminAndOrdinaryUsers:
		err = d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? or role_level = ?)", groupID, constant.GroupAdmin, constant.GroupOrdinaryUsers).Order("role_level DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	case constant.GroupFilterOwnerAndAdmin:
		err = d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? or role_level = ?)", groupID, constant.GroupOwner, constant.GroupAdmin).Order("role_level DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	default:
		return nil, fmt.Errorf("filter args failed %d", filter)
	}
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListSplit failed ")
}

func (d *DataBase) GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? And (role_level = ? OR role_level = ?)", groupID, constant.GroupOwner, constant.GroupAdmin).Order("join_time DESC").Find(&groupMemberList).Error
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListSplit failed ")
}

func (d *DataBase) GetGroupMemberOwner(ctx context.Context, groupID string) (*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMember model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ? And role_level = ?", groupID, constant.GroupOwner).Find(&groupMember).Error
	return &groupMember, utils.Wrap(err, "GetGroupMemberListSplit failed ")
}

func (d *DataBase) GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	var err error
	if len(userIDList) == 0 {
		err = d.conn.WithContext(ctx).Where("group_id = ? And join_time  between ? and ? ", groupID, joinTimeBegin, joinTimeEnd).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	} else {
		err = d.conn.WithContext(ctx).Where("group_id = ? And join_time  between ? and ? And user_id NOT IN ?", groupID, joinTimeBegin, joinTimeEnd, userIDList).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	}
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListSplitByJoinTimeFilter failed ")
}

func (d *DataBase) GetGroupOwnerAndAdminByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupMemberList []model_struct.LocalGroupMember
	err := d.conn.WithContext(ctx).Where("group_id = ?  AND (role_level = ? Or role_level = ?)", groupID, constant.GroupOwner, constant.GroupAdmin).Find(&groupMemberList).Error
	var transfer []*model_struct.LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupMemberUIDListByGroupID(ctx context.Context, groupID string) (result []string, err error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var g model_struct.LocalGroupMember
	g.GroupID = groupID
	err = d.conn.WithContext(ctx).Model(&g).Where("group_id = ?", groupID).Pluck("user_id", &result).Error
	return result, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(groupMember).Error, "")
}

//funcation (d *DataBase) BatchInsertMessageList(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
//	if MessageList == nil {
//		return nil
//	}
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	return utils.Wrap(d.conn.WithContext(ctx).Create(MessageList).Error, "BatchInsertMessageList failed")
//}

func (d *DataBase) BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	if groupMemberList == nil {
		return errors.New("nil")
	}
	return utils.Wrap(d.conn.WithContext(ctx).Create(groupMemberList).Error, "BatchInsertMessageList failed")
}

func (d *DataBase) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	groupMember := model_struct.LocalGroupMember{}
	return d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (d *DataBase) DeleteGroupAllMembers(ctx context.Context, groupID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	groupMember := model_struct.LocalGroupMember{}
	return d.conn.WithContext(ctx).Where("group_id=? ", groupID).Delete(&groupMember).Error
}

func (d *DataBase) UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(groupMember).Select("*").Updates(*groupMember)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) UpdateGroupMemberField(ctx context.Context, groupID, userID string, args map[string]interface{}) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	c := model_struct.LocalGroupMember{GroupID: groupID, UserID: userID}
	t := d.conn.WithContext(ctx).Model(&c).Updates(args)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateGroupMemberField failed")
}

func (d *DataBase) GetGroupMemberInfoIfOwnerOrAdmin(ctx context.Context) ([]*model_struct.LocalGroupMember, error) {
	var ownerAndAdminList []*model_struct.LocalGroupMember
	groupList, err := d.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	for _, v := range groupList {
		memberList, err := d.GetGroupOwnerAndAdminByGroupID(ctx, v.GroupID)
		if err != nil {
			return nil, utils.Wrap(err, "")
		}
		ownerAndAdminList = append(ownerAndAdminList, memberList...)
	}
	return ownerAndAdminList, nil
}

func (d *DataBase) SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	if !isSearchMemberNickname && !isSearchUserID {
		return nil, errors.New("args failed")
	}

	var countCon int
	var condition string
	if isSearchUserID {
		condition = fmt.Sprintf("user_id like %q ", "%"+keyword+"%")
		countCon++
	}
	if isSearchMemberNickname {
		if countCon > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("nickname like %q ", "%"+keyword+"%")
	}

	var groupMemberList []model_struct.LocalGroupMember
	if groupID != "" {
		condition = "( " + condition + " ) "
		condition += " and group_id IN ? "
		log.Debug("", "subCondition SearchGroupMembers ", condition)
		err = d.conn.WithContext(ctx).Where(condition, []string{groupID}).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
	} else {
		log.Debug("", "subCondition SearchGroupMembers ", condition)
		err = d.conn.WithContext(ctx).Where(condition).Order("join_time DESC").Offset(offset).Limit(count).Find(&groupMemberList).Error
		log.Debug("", "subCondition SearchGroupMembers ", condition, len(groupMemberList))
	}

	for _, v := range groupMemberList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}

func (d *DataBase) GetGroupMemberAllGroupIDs(ctx context.Context) ([]string, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	// SELECT DISTINCT group_id FROM local_group_members;
	var groupIDs []string
	err := d.conn.WithContext(ctx).Select("DISTINCT group_id").Model(&model_struct.LocalGroupMember{}).Pluck("group_id", &groupIDs).Error
	if err != nil {
		return nil, err
	}
	return groupIDs, nil
}

func (d *DataBase) SearchKickMemberList(ctx context.Context, params sdk_params_callback.GetKickGroupListReq) ([]*sdk_params_callback.KickGroupList, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	tx1 := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).
		Select([]string{
			"group_id", "user_id", "nickname", "role_level", "join_time", "face_url",
			"code", "phone", "gender", "group_user_name",
		}).Where("group_id = ? AND user_id != ?", params.GroupID, params.UserID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.IsManger { //管理员只可踢普通用户
				db.Where("role_level = ?", 1)
			}
			db.Where("nickname LIKE ?", "%"+params.Name+"%")
			return db
		})
	tx2 := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).Select([]string{
		"group_id", "user_id", "nickname", "role_level", "join_time", "face_url",
		"code", "phone", "gender", "group_user_name",
	}).Where("group_id = ? AND user_id != ?", params.GroupID, params.UserID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.IsManger { //管理员只可踢普通用户
				db.Where("role_level = ?", 1)
			}
			db.Where("group_user_name LIKE ?", "%"+params.Name+"%")
			return db
		})
	tx3 := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupMember{}).
		Select([]string{
			"group_id", "user_id", "nickname", "role_level", "join_time", "face_url",
			"code", "phone", "gender", "group_user_name",
		}).Where("group_id = ? AND user_id != ?", params.GroupID, params.UserID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.IsManger { //管理员只可踢普通用户
				db.Where("role_level = ?", 1)
			}
			db.Where("phone LIKE ?", "%"+params.Name+"%")
			return db
		})
	total := int64(0)
	result := make([]*sdk_params_callback.KickGroupList, 0)
	err := d.conn.WithContext(ctx).Table("(?) AS t", d.conn.Raw("? UNION ? UNION ?", tx1, tx2, tx3)).
		Count(&total).Order("join_time DESC").
		Offset((params.PageNum - 1) * params.PageSize).Limit(params.PageSize).Find(&result).Error
	return result, err
}
