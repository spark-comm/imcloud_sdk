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

func (d *DataBase) InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var localFriend model_struct.LocalFriend
	if d.conn.WithContext(ctx).Where("owner_user_id = ? and friend_user_id = ?", friend.OwnerUserID, friend.FriendUserID).First(&localFriend).RowsAffected == 0 {
		// 记录不存在，创建新记录
		return utils.Wrap(d.conn.WithContext(ctx).Create(friend).Error, "InsertFriend failed")
	} else {
		// 记录存在，更新记录
		return utils.Wrap(d.conn.WithContext(ctx).Model(&localFriend).Updates(friend).Error, "InsertFriend failed")
	}
}

func (d *DataBase) DeleteFriendDB(ctx context.Context, friendUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&model_struct.LocalFriend{}).Error, "DeleteFriend failed")
}

func (d *DataBase) UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()

	t := d.conn.WithContext(ctx).Model(friend).Select("*").Updates(*friend)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}
func (d *DataBase) GetAllFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	err := utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"GetFriendList failed")
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) GetPageFriendList(ctx context.Context, offset, count int) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []*model_struct.LocalFriend
	err := utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ?", d.loginUserID).Offset(offset).Limit(count).Order("sort_flag").Find(&friendList).Error,
		"GetFriendList failed")
	return friendList, err
}
func (d *DataBase) GetFriendsByPage(ctx context.Context, page, size int) ([]*model_struct.LocalFriend, int64, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var total int64
	var friendList []*model_struct.LocalFriend
	tx := d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{}).Where("owner_user_id = ?", d.loginUserID)
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, utils.Wrap(err, "GetFriendList failed")
	}
	err := utils.Wrap(tx.Order("sort_flag").Scopes(SqlDataLimit(size, page)).Scan(&friendList).Error,
		"GetFriendList failed")
	return friendList, total, err
}
func (d *DataBase) SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var count int
	var friendList []model_struct.LocalFriend
	var condition string
	if isSearchUserID {
		condition = fmt.Sprintf("friend_user_id like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchNickname {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("nickname like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchRemark {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("remark like %q ", "%"+keyword+"%")
	}
	err := d.conn.WithContext(ctx).Where(condition).Order("create_time DESC").Find(&friendList).Error
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "SearchFriendList failed ")

}
func (d *DataBase) SearchFriends(ctx context.Context, keyword string, notPeersFriend bool, page, size int) ([]*model_struct.LocalFriend, int64, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var total int64
	tx := d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{})
	if keyword != "" {
		var condition string
		fields := []string{"remark", "nickname", "code", "phone", "email"}
		for i, field := range fields {
			if i == 0 {
				condition = fmt.Sprintf("%s like %q ", field, "%"+keyword+"%")
			} else {
				condition += fmt.Sprintf("or %s like %q ", field, "%"+keyword+"%")
			}
		}
		tx = tx.Where(condition)
	}
	if notPeersFriend {
		tx = tx.Where(" not_peers_friend = 0")
	}
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, utils.Wrap(err, "SearchFriends failed")
	}
	var friendList []*model_struct.LocalFriend
	err := tx.Order("sort_flag ASC").Scopes(SqlDataLimit(size, page)).Scan(&friendList).Error
	return friendList, total, utils.Wrap(err, "SearchFriendList failed ")
}
func (d *DataBase) GetFriendsNotInGroup(ctx context.Context, groupID, keyword string, page, size int) ([]*model_struct.LocalFriend, int64, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []*model_struct.LocalFriend
	var totalCount int64
	tx := d.conn.Table("local_friends").
		Joins("LEFT JOIN local_group_members ON local_friends.friend_user_id = local_group_members.user_id AND local_group_members.group_id = ?", groupID).
		Where("local_group_members.user_id IS NULL and local_friends.not_peers_friend = 0")
	if keyword != "" {
		var condition string
		fields := []string{"remark", "nickname", "code", "phone", "email"}
		for i, field := range fields {
			if i == 0 {
				condition = fmt.Sprintf("local_friends.%s like %q ", field, "%"+keyword+"%")
			} else {
				condition += fmt.Sprintf("or local_friends.%s like %q ", field, "%"+keyword+"%")
			}
		}
		tx = tx.Where(condition)
	}
	if err := tx.Count(&totalCount).Error; err != nil {
		return nil, 0, utils.Wrap(err, "GetFriendsNotInGroup failed")
	}
	err := tx.WithContext(ctx).Select("local_friends.*").Order("sort_flag ASC").Scopes(SqlDataLimit(size, page)).Scan(&friendList).Error
	return friendList, totalCount, utils.Wrap(err, "GetFriendsNotInGroup failed ")
}
func (d *DataBase) GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friend model_struct.LocalFriend
	return &friend, utils.Wrap(d.conn.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Take(&friend).Error, "GetFriendInfoByFriendUserID failed")
}

func (d *DataBase) GetFriendInfoList(ctx context.Context, friendUserIDList []string, filterNotPeersFriend bool) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	var err error
	if filterNotPeersFriend {
		err = utils.Wrap(d.conn.WithContext(ctx).Where("not_peers_friend = 0 and friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	} else {
		err = utils.Wrap(d.conn.WithContext(ctx).Where("friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	}
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
func (d *DataBase) UpdateColumnsFriend(ctx context.Context, friendIDs []string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	// Update records where FriendUserID is in the friendIDs slice
	t := d.conn.WithContext(ctx).Model(&model_struct.LocalFriend{}).Where("friend_user_id IN ?", friendIDs).Updates(args)

	return utils.Wrap(t.Error, "UpdateColumnsFriend failed")
}
