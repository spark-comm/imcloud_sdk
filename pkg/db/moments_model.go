package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

// InsertMomentsComments ， 插入朋友圈评论
// 参数：
//      ctx ： desc
//      moments ： desc
// 返回值：
//      error ：desc
func (d *DataBase) InsertMomentsComments(ctx context.Context, moments *model_struct.LocalMomentsComments) error {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	return utils.Wrap(d.conn.WithContext(ctx).Create(moments).Error, "Insert LocalMomentsComments failed")
}

// DeleteMomentsComments ， 删除朋友圈评论
// 参数：
//      ctx ： desc
//      commentId ： desc
// 返回值：
//      error ：desc
func (d *DataBase) DeleteMomentsComments(ctx context.Context, commentId string) error {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	local := model_struct.LocalMomentsComments{CommentId: commentId}
	return utils.Wrap(d.conn.WithContext(ctx).Delete(&local).Error, "Delete LocalMomentsComments failed")
}

// UpdateMomentsComments ， 更新朋友圈评论
// 参数：
//      ctx ： desc
//      moments ： desc
// 返回值：
//      error ：desc
func (d *DataBase) UpdateMomentsComments(ctx context.Context, moments *model_struct.LocalMomentsComments) error {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	t := d.conn.WithContext(ctx).Model(&model_struct.LocalMomentsComments{}).Updates(moments)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

// GetMomentsComments ， 获取某个朋友圈所有的评论
// 参数：
//      ctx ： desc
//      momentId ： desc
// 返回值：
//      error ：desc
func (d *DataBase) GetMomentsComments(ctx context.Context, momentId string) ([]*model_struct.LocalMomentsComments, error) {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	var res []*model_struct.LocalMomentsComments
	t := d.conn.WithContext(ctx).
		Model(&model_struct.LocalMomentsComments{}).
		Where("moment_id = ?", momentId).
		Where("deleted_at = 0").
		Order("created_at asc"). // 按时间顺序排列
		Find(res)
	if err := t.Error; err != nil {
		return nil, utils.Wrap(t.Error, "")
	}
	return res, nil
}

// InsertMoments ， 插入朋友圈
// 参数：
//      ctx ： desc
//      moments ： desc
// 返回值：
//      error ：desc
func (d *DataBase) InsertMoments(ctx context.Context, moments *model_struct.LocalMoments) error {
	//TODO implement me
	d.momentsMtx.Lock()
	defer d.momentsMtx.Unlock()

	return utils.Wrap(d.conn.WithContext(ctx).Create(moments).Error, "Insert LocalMoments failed")
}

// InsertBatchMoments ， 批量插入朋友圈
// 参数：
//      ctx ： desc
//      moments ： desc
// 返回值：
//      error ：desc
func (d *DataBase) InsertBatchMoments(ctx context.Context, moments []*model_struct.LocalMoments) error {
	//TODO implement me
	d.momentsMtx.Lock()
	defer d.momentsMtx.Unlock()

	return utils.Wrap(d.conn.WithContext(ctx).Create(moments).Error, "Insert LocalMoments failed")
}

// DeleteMoments ， 删除朋友圈
// 参数：
//      ctx ： desc
//      momentId ： desc
// 返回值：
//      error ：desc
func (d *DataBase) DeleteMoments(ctx context.Context, momentId string) error {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	local := model_struct.LocalMoments{MomentId: momentId}
	return utils.Wrap(d.conn.WithContext(ctx).Delete(&local).Error, "Delete LocalMoments failed")
}

// UpdateMoments ， 更新朋友圈
// 参数：
//      ctx ： desc
//      moments ： desc
// 返回值：
//      error ：desc
func (d *DataBase) UpdateMoments(ctx context.Context, moments *model_struct.LocalMoments) error {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	t := d.conn.WithContext(ctx).Model(&model_struct.LocalMoments{}).Updates(moments)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

// GetMoments ， 获取朋友圈
// 参数：
//      ctx ： desc
//      momentId ： desc
// 返回值：
//      error ：desc
func (d *DataBase) GetMoments(ctx context.Context, momentId string) (*model_struct.LocalMoments, error) {
	//TODO implement me
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	var res *model_struct.LocalMoments
	t := d.conn.WithContext(ctx).
		Model(&model_struct.LocalMoments{}).
		Where("moment_id = ?", momentId).
		Where("deleted_at = 0").
		Order("created_at asc"). // 按时间顺序排列
		First(res)
	if err := t.Error; err != nil {
		return nil, utils.Wrap(t.Error, "")
	}
	return res, nil
}

func (d *DataBase) GetMomentsList(ctx context.Context, page, size int, isSelf bool) ([]*model_struct.LocalMoments, error) {
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	var res []*model_struct.LocalMoments
	t := d.conn.WithContext(ctx).
		Model(&model_struct.LocalMoments{}).
		Where("deleted_at = 0").
		Order("created_at desc") // 按时间顺序排列

	if page != 0 && size != 0 {
		t.Limit(size).Offset((page - 1) * size)
	}

	if err := t.Find(&res).Error; err != nil {
		return nil, utils.Wrap(t.Error, "")
	}
	return res, nil
}

const (
	MOMENTS_FIND_TIMESTAMPS_TYPE_FIRST = 0
	MOMENTS_FIND_TIMESTAMPS_TYPE_LAST  = 1
)

func (d *DataBase) GetMomentTimestamps(ctx context.Context, t int) (int64, error) { // 获取最后/最前同步同步时间戳（created_at） 0 最后 1 最前
	d.momentsCommentsMtx.Lock()
	defer d.momentsCommentsMtx.Unlock()

	var res *model_struct.LocalMoments
	tx := d.conn.WithContext(ctx).
		Model(&model_struct.LocalMoments{}).
		Where("deleted_at = 0").
		Limit(1)

	if t == MOMENTS_FIND_TIMESTAMPS_TYPE_FIRST {
		tx.Order("created_at asc") // 按时间顺序排列
	} else {
		tx.Order("created_at desc") // 按时间顺序排列
	}
	tx.First(&res)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	if tx.Error != nil {
		return 0, tx.Error
	}

	return res.CreatedAt, nil
}
