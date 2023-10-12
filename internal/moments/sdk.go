package moments

import (
	"context"
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdk_params_callback"
	"strings"
)

// TODO: 事务处理

// Publish ， 发布朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) Publish(ctx context.Context, params *sdk_params_callback.PublishRequest) error {

	// 发布
	svr, err := m.publishMoments2Svr(ctx, params)
	if err != nil {
		return err
	}

	var img string
	if len(params.Images) > 0 {
		img = strings.Join(params.Images, ",")
	}

	// 直接插入
	if err := m.db.InsertMoments(ctx, &model_struct.LocalMoments{
		CreatedAt: svr.CreatedAt,
		MomentId:  svr.MomentId,
		UserId:    m.loginUserID,
		UserName:  params.UserName,
		Avatar:    params.Avatar,
		Content:   params.Content,
		Images:    img,
		VideoUrl:  params.VideoUrl,
		VideoImg:  params.VideoImg,
		Location:  params.Location,
		Longitude: params.Longitude,
		Latitude:  params.Latitude,
	}); err != nil {
		return err
	}

	// 同步最新的数据至本地
	//if err := m.SyncNewMomentsFromSvr(ctx); err != nil {
	//	return err
	//}

	return nil
}

// Delete ， 删除朋友圈
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) Delete(ctx context.Context, params *sdk_params_callback.DeleteRequest) error {
	// 删除服务器
	if err := m.deletedMoments2Svr(ctx, params); err != nil {
		return err
	}

	// 删除本地
	if err := m.db.DeleteMoments(ctx, params.MomentId); err != nil {
		return err
	}

	return nil
}

// Comment ， 评论
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) Comment(ctx context.Context, params *sdk_params_callback.CommentRequest) error {

	// 评论
	svr, err := m.commentMoments2Svr(ctx, params)
	if err != nil {
		return err
	}

	// 插入本地数据库
	if err := m.db.InsertMomentsComments(ctx, &model_struct.LocalMomentsComments{
		CreatedAt:      svr.CreatedAt,
		CommentId:      svr.CommentId,
		MomentId:       params.MomentId,
		UserId:         m.loginUserID,
		Type:           int(params.Type),
		Avatar:         params.Avatar,
		NickName:       params.Nickname,
		SourceUserId:   params.SourceUserId,
		SourceAvatar:   params.SourceAvatar,
		SourceNickName: params.SourceNickname,
		Content:        params.Content,
	}); err != nil {
		return err
	}

	return nil
}

// Like ， 点赞
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
const (
	MomentsIsLike = 1 // 已点赞
	MomentsUnLike = 2 // 未点赞
)

func (m *Moments) Like(ctx context.Context, params *sdk_params_callback.LikeRequest) error {
	// 点赞
	if err := m.likeMoments2Svr(ctx, params); err != nil {
		return err
	}

	// 同步本地
	if err := m.db.UpdateMoments(ctx, params.MomentId, map[string]interface{}{
		"is_like": MomentsIsLike,
		"likes":   gorm.Expr("likes + 1"),
	}); err != nil {
		return err
	}

	return nil
}

// UnLike ， 取消点赞
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
//      error ：desc
func (m *Moments) UnLike(ctx context.Context, params *sdk_params_callback.UnlikeRequest) error {
	// 取消点赞
	if err := m.unlikeMoments2Svr(ctx, params); err != nil {
		return err
	}

	// 同步本地
	if err := m.db.UpdateMoments(ctx, params.MomentId, map[string]interface{}{
		"is_like": MomentsUnLike,
		"likes":   gorm.Expr("likes - 1"),
	}); err != nil {
		return err
	}

	return nil
}

// GetMomentsList ， 获取朋友圈列表
// 参数：
//      ctx ： desc
//      params ： desc
// 返回值：
func (m *Moments) GetMomentsList(ctx context.Context, params *sdk_params_callback.V2ListRequest) ([]*sdk_params_callback.V2ListReply, error) {

	// 如果是第一页，先同步最新消息
	if params.Page <= 1 {
		err := m.SyncNewMomentsFromSvr(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 从数据库获取朋友圈
	list, err := m.db.GetMomentsList(ctx, int(params.Page), int(params.Size), params.IsSelf, m.loginUserID)
	if err != nil {
		return nil, err
	}

	var res []*sdk_params_callback.V2ListReply
	// 获取评论
	for _, v := range list {
		one := sdk_params_callback.V2ListReply{
			LocalMoments: *v,
		}

		comments, err := m.db.GetMomentsComments(ctx, v.MomentId)
		if err != nil {
			return nil, err
		}

		one.Comments = comments

		if v.Images != "" {
			one.Images = strings.Split(v.Images, ",")
		}

		res = append(res, &one)
	}

	return res, err
}
