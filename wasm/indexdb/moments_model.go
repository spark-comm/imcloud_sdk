package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

type LocalMoments struct {
}

func NewLocalMoments() *LocalMoments {
	return &LocalMoments{}
}

func (d *LocalMoments) InsertMomentsComments(ctx context.Context, moments *model_struct.LocalMomentsComments) error {
	return nil
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
	return nil
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
	return nil
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
	return nil
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
func (d *LocalMoments) UpdateMoments(ctx context.Context, moments *model_struct.LocalMoments) error {
	return nil
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
	return &model_struct.LocalMoments{}, nil
}

const (
	MOMENTS_FIND_TIMESTAMPS_TYPE_FIRST = 0
	MOMENTS_FIND_TIMESTAMPS_TYPE_LAST  = 1
)

func (d *LocalMoments) GetMomentTimestamps(ctx context.Context, t int) (int64, error) { // 获取最后/最前同步同步时间戳（created_at） 0 最后 1 最前
	return 0, nil
}
