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

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/exec"
)

type LocalMoments struct {
}

func NewLocalMoments() *LocalMoments {
	return &LocalMoments{}
}

// InsertMomentsComments ， 插入朋友圈评论
// 参数：
//
//	ctx ： desc
//	comment ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) InsertMomentsComments(ctx context.Context, comment *model_struct.LocalMomentsComments) error {
	_, err := exec.Exec(utils.StructToJsonString(comment))
	return err
}

// InsertBatchMomentsComments ， 批量插入朋友圈评论
// 参数：
//
//	ctx ： desc
//	comments ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) InsertBatchMomentsComments(ctx context.Context, comments []*model_struct.LocalMomentsComments) error {
	_, err := exec.Exec(utils.StructToJsonString(comments))
	return err
}

// DeleteMomentsComments ， 删除朋友圈评论
// 参数：
//
//	ctx ： desc
//	commentId ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) DeleteMomentsComments(ctx context.Context, commentId string) error {
	//_, err := exec.Exec(moments)
	//return err
	return nil
}

// UpdateMomentsComments ， 更新朋友圈评论
// 参数：
//
//	ctx ： desc
//	moments ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) UpdateMomentsComments(ctx context.Context, moments *model_struct.LocalMomentsComments) error {
	_, err := exec.Exec(utils.StructToJsonString(moments))
	return err
}

// GetMomentsComments ， 获取某个朋友圈所有的评论
// 参数：
//
//	ctx ： desc
//	momentId ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) GetMomentsComments(ctx context.Context, momentId string) ([]*model_struct.LocalMomentsComments, error) {
	//return exec.Exec(momentId)
	return []*model_struct.LocalMomentsComments{}, nil
}

// InsertMoments ， 插入朋友圈
// 参数：
//
//	ctx ： desc
//	moments ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) InsertMoments(ctx context.Context, moments *model_struct.LocalMoments) error {
	_, err := exec.Exec(utils.StructToJsonString(moments))
	return err
}

// InsertBatchMoments ， 批量插入朋友圈
// 参数：
//
//	ctx ： desc
//	moments ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) InsertBatchMoments(ctx context.Context, moments []*model_struct.LocalMoments) error {
	_, err := exec.Exec(utils.StructToJsonString(moments))
	return err
}

// DeleteMoments ， 删除朋友圈
// 参数：
//
//	ctx ： desc
//	momentId ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) DeleteMoments(ctx context.Context, momentId string) error {
	_, err := exec.Exec(momentId)
	return err
}

// UpdateMoments ， 更新朋友圈
// 参数：
//
//	ctx ： desc
//	moments ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) UpdateMoments(ctx context.Context, momentId string, moments interface{}) error {
	_, err := exec.Exec(momentId, utils.StructToJsonString(moments))
	return err
}

// GetMoments ， 获取朋友圈
// 参数：
//
//	ctx ： desc
//	momentId ： desc
//
// 返回值：
//
//	error ：desc
func (d *LocalMoments) GetMoments(ctx context.Context, momentId string) (*model_struct.LocalMoments, error) {
	//return exec.Exec(momentId)
	return &model_struct.LocalMoments{}, nil
}

const (
	MOMENTS_FIND_TIMESTAMPS_TYPE_FIRST = 0
	MOMENTS_FIND_TIMESTAMPS_TYPE_LAST  = 1
)

// GetMomentTimestamps ， 获取时间戳
// 参数：
//
//	ctx ： desc
//	t ： desc
//
// 返回值：
//
//	int64 ：desc
//	error ：desc
func (d *LocalMoments) GetMomentTimestamps(ctx context.Context, t int) (int64, error) { // 获取最后/最前同步同步时间戳（created_at） 0 最后 1 最前
	//return exec.Exec(t)
	return 0, nil
}

// GetMomentsList ， 获取朋友圈列表
// 参数：
//
//	ctx ： desc
//	page ： desc
//	size ： desc
//	isSelf ： desc
//	userId ： desc
//
// 返回值：
//
//	[]*model_struct.LocalMoments ：desc
//	error ：desc
func (d *LocalMoments) GetMomentsList(ctx context.Context, page, size int, isSelf bool, userId string) ([]*model_struct.LocalMoments, error) {
	_, err := exec.Exec(page, size, isSelf, userId)
	return []*model_struct.LocalMoments{}, err
}
