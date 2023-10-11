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

package moments

import (
	"context"
	momentsv1 "github.com/imCloud/api/moments/v1"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdk_params_callback"
)

const (
	MAX_SIZE = 50
)

type Moments struct {
	listener    open_im_sdk_callback.OnMomentsListener // TODO: 待实现
	loginUserID string
	db          db_interface.DataBase
	//momentsSyncer         *syncer.Syncer[*model_struct.LocalMoments, string]
	//momentsCommentsSyncer *syncer.Syncer[*model_struct.LocalMomentsComments, [2]string]
}

//func (m *Moments) initSyncer() {
//	m.momentsSyncer = syncer.New(m.db.InsertMoments, func(ctx context.Context, value *model_struct.LocalMoments) error {
//		return m.db.DeleteMoments(ctx, value.MomentId)
//	}, func(ctx context.Context, server, local *model_struct.LocalMoments) error {
//		return m.db.UpdateMoments(ctx, local.MomentId, server)
//	}, func(value *model_struct.LocalMoments) string {
//		return value.MomentId
//	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalMoments) error {
//		return nil
//	})
//}

// SyncNewMomentsFromSvr ， 从服务器同步新朋友圈
// 参数：
//      ctx ： desc
// 返回值：
func (m *Moments) SyncNewMomentsFromSvr(ctx context.Context) error {
	timestamps, err := m.db.GetMomentTimestamps(ctx, db.MomentsFindTimestampsTypeLast)
	if err != nil {
		return err
	}

	// 获取本地最新时间戳
	var currentTimestamps int64
	var needSync []*model_struct.LocalMoments
	// 获取本地最新时间戳

out:
	for {
		svr, i, err := m.getMomentsFromSvr(ctx, currentTimestamps)
		if err != nil {
			return err
		}

		// 如若服务器没有数据，或者当前最新的和本地的相等，则无更新，直接退出
		if len(svr) <= 0 || svr[0].CreatedAt <= timestamps {
			break
		}

		// 循环判断 如果时间戳大于本地最大时间戳，则为新圈子 保存，等待插入
		for _, v := range svr {
			if v.CreatedAt > timestamps {
				s2local := ServerMomentsToLocalMoments([]*momentsv1.ListItem{v})
				needSync = append(needSync, s2local...)

				// 同步评论
				s2localComment := ServerMomentsCommentsToLocalMomentsComments(v.Comments)
				err := m.db.InsertBatchMomentsComments(ctx, s2localComment)
				if err != nil {
					return err
				}
			} else {
				// 如果同步完成则退出
				break out
			}
		}

		if len(svr) <= MAX_SIZE {
			break
		}

		currentTimestamps = i
	}

	// 插入
	if len(needSync) > 0 {
		if err := m.db.InsertBatchMoments(ctx, needSync); err != nil {
			return err
		}
	}

	return nil
}

// SyncHistoryMomentsFromSvr ， 从服务器同步历史朋友圈
// 参数：
//      ctx ： desc
// 返回值：
func (m *Moments) SyncHistoryMomentsFromSvr(ctx context.Context) error {
	timestamps, err := m.db.GetMomentTimestamps(ctx, db.MomentsFindTimestampsTypeFirst)
	if err != nil {
		return err
	}

	// 获取本地最新时间戳
	var currentTimestamps = timestamps
	for {
		svr, i, err := m.getMomentsFromSvr(ctx, currentTimestamps)
		if err != nil {
			return err
		}

		// 如若服务器没有数据，或者当前最后一条的和本地的相等，则无更新，直接退出
		if len(svr) <= 0 {
			break
		}

		// 插入
		if err := m.db.InsertBatchMoments(ctx, ServerMomentsToLocalMoments(svr)); err != nil {
			return err
		}

		// 插入评论
		for _, val := range svr {
			// 同步评论
			s2localComment := ServerMomentsCommentsToLocalMomentsComments(val.Comments)
			err := m.db.InsertBatchMomentsComments(ctx, s2localComment)
			if err != nil {
				return err
			}
		}

		// 单次
		if len(svr) <= MAX_SIZE {
			break
		}

		currentTimestamps = i
	}

	return nil
}

// getMomentsFromSvr ， 从服务器获取朋友圈
// 参数：
//      ctx ： desc
//      timestamps ： desc
// 返回值：
//      []*momentsv1.ListItem ：desc
//      int64 ：desc
//      error ：desc
func (m *Moments) getMomentsFromSvr(ctx context.Context, timestamps int64) ([]*momentsv1.ListItem, int64, error) {
	resp := &momentsv1.V2ListReply{}
	err := util.CallPostApi[*momentsv1.V2ListRequest, *momentsv1.V2ListReply](
		ctx, constant.V2ListMomentsRouter, &momentsv1.V2ListRequest{
			UserId:    m.loginUserID,
			IsSelf:    true,
			Page:      1,
			Size:      MAX_SIZE,
			Timestamp: timestamps,
		}, resp,
	)
	if err != nil {
		return nil, 0, err
	}

	return resp.List, resp.Timestamp, nil
}

// publishMoments2Svr ， 发布朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      *momentsv1.PublishReply ：desc
//      error ：desc
func (m *Moments) publishMoments2Svr(ctx context.Context, params *sdk_params_callback.PublishRequest) (*momentsv1.PublishReply, error) {
	resp := &momentsv1.PublishReply{}
	err := util.CallPostApi[*momentsv1.PublishRequest, *momentsv1.PublishReply](
		ctx, constant.V2PublishMomentsRouter, &momentsv1.PublishRequest{
			UserId:    m.loginUserID,
			UserName:  params.UserName,
			Avatar:    params.Avatar,
			Content:   params.Content,
			Images:    params.Images,
			VideoUrl:  params.VideoUrl,
			VideoImg:  params.VideoImg,
			Location:  params.Location,
			Longitude: params.Longitude,
			Latitude:  params.Latitude,
		}, resp,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// commentMoments2Svr ， 评论朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      *momentsv1.CommentReply ：desc
//      error ：desc
func (m *Moments) commentMoments2Svr(ctx context.Context, params *sdk_params_callback.CommentRequest) (*momentsv1.CommentReply, error) {
	resp := &momentsv1.CommentReply{}
	err := util.CallPostApi[*momentsv1.CommentRequest, *momentsv1.CommentReply](
		ctx, constant.V2CommentMomentsRouter, &momentsv1.CommentRequest{
			UserId:         m.loginUserID,
			MomentId:       params.MomentId,
			Type:           params.Type,
			Avatar:         params.Avatar,
			Nickname:       params.Nickname,
			Content:        params.Content,
			SourceUserId:   params.SourceUserId,
			SourceAvatar:   params.SourceAvatar,
			SourceNickname: params.SourceNickname,
		}, resp,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// deletedMoments2Svr ， 删除朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) deletedMoments2Svr(ctx context.Context, params *sdk_params_callback.DeleteRequest) error {
	resp := &momentsv1.DeleteReply{}
	err := util.CallPostApi[*momentsv1.DeleteRequest, *momentsv1.DeleteReply](
		ctx, constant.V2DeleteMomentsRouter, &momentsv1.DeleteRequest{
			UserId:   m.loginUserID,
			MomentId: params.MomentId,
		}, resp,
	)
	if err != nil {
		return err
	}

	return nil
}

// likeMoments2Svr ， 点赞朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) likeMoments2Svr(ctx context.Context, params *sdk_params_callback.LikeRequest) error {
	resp := &momentsv1.LikeReply{}
	err := util.CallPostApi[*momentsv1.LikeRequest, *momentsv1.LikeReply](
		ctx, constant.V2LikeMomentsRouter, &momentsv1.LikeRequest{
			UserId:       m.loginUserID,
			MomentId:     params.MomentId,
			UserNickname: params.UserNickname,
		}, resp,
	)
	if err != nil {
		return err
	}

	return nil
}

// unlikeMoments2Svr ， 取消点赞朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) unlikeMoments2Svr(ctx context.Context, params *sdk_params_callback.UnlikeRequest) error {
	resp := &momentsv1.UnlikeReply{}
	err := util.CallPostApi[*momentsv1.UnlikeRequest, *momentsv1.UnlikeReply](
		ctx, constant.V2UnlikeMomentsRouter, &momentsv1.UnlikeRequest{
			UserId:   m.loginUserID,
			MomentId: params.MomentId,
		}, resp,
	)
	if err != nil {
		return err
	}

	return nil
}
