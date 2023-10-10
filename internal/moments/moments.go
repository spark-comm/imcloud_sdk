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
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
)

type Moments struct {
	listener              open_im_sdk_callback.OnMomentsListener // TODO: 待实现
	loginUserID           string
	db                    db_interface.DataBase
	momentsSyncer         *syncer.Syncer[*model_struct.LocalMoments, string]
	momentsCommentsSyncer *syncer.Syncer[*model_struct.LocalMomentsComments, [2]string]
}

func (m *Moments) initSyncer() {
	m.momentsSyncer = syncer.New(m.db.InsertMoments, func(ctx context.Context, value *model_struct.LocalMoments) error {
		return m.db.DeleteMoments(ctx, value.MomentId)
	}, func(ctx context.Context, server, local *model_struct.LocalMoments) error {
		return m.db.UpdateMoments(ctx, server)
	}, func(value *model_struct.LocalMoments) string {
		return value.MomentId
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalMoments) error {
		return nil
	})
}

// SyncNewMomentsFromSvr ， 从服务器同步新朋友圈
// 参数：
//      ctx ： desc
// 返回值：
func (m *Moments) SyncNewMomentsFromSvr(ctx context.Context) {
	// 获取本地最新时间戳
	//timestamps, err := m.db.GetMomentTimestamps(ctx, db.MOMENTS_FIND_TIMESTAMPS_TYPE_FIRST)
	//if err != nil {
	//	return
	//}
}

func (m *Moments) Sync() {

}

func (m *Moments) getMomentsFromSvr(ctx context.Context, timestamps int64) ([]*model_struct.LocalMoments, int64, error) {
	resp := &momentsv1.V2ListReply{}
	err := util.CallPostApi[*momentsv1.V2ListRequest, *momentsv1.V2ListReply](
		ctx, constant.GetGroupsInfoRouter, &momentsv1.V2ListRequest{}, resp,
	)
	if err != nil {
		return nil, 0, err
	}

	return ServerMomentsToLocalMoments(resp.List), resp.Timestamp, nil
}

// getGroupsInfoFromSvr 从服务端获取群数据
//func (g *Group) getGroupsInfoFromSvr(ctx context.Context, groupIDs []string) ([]*groupv1.GroupInfo, error) {
//	//resp, err := util.CallApi[groupv1.GetGroupInfoResponse](ctx, constant.GetGroupsInfoRouter, &groupv1.GetGroupInfoReq{GroupID: groupIDs})
//	resp := &groupv1.GetGroupInfoResponse{}
//	err := util.CallPostApi[*groupv1.GetGroupInfoReq, *groupv1.GetGroupInfoResponse](
//		ctx, constant.GetGroupsInfoRouter, &groupv1.GetGroupInfoReq{GroupID: groupIDs}, resp,
//	)
//	if err != nil {
//		return nil, err
//	}
//	return resp.Data, nil
//}
