package moments

import (
	momentsv1 "github.com/imCloud/api/moments/v1"
	"open_im_sdk/pkg/db/model_struct"
	"strings"
)

// ServerMomentsToLocalMoments ， 服务端朋友圈对象转本地对象
// 参数：
//      info ： desc
// 返回值：
//      []*model_struct.LocalMoments ：desc
func ServerMomentsToLocalMoments(list []*momentsv1.ListItem) []*model_struct.LocalMoments {
	var res []*model_struct.LocalMoments
	for _, v := range list {
		item := &model_struct.LocalMoments{
			CreatedAt: v.CreatedAt,
			MomentId:  v.MomentId,
			UserId:    v.UserId,
			UserName:  v.UserName,
			Avatar:    v.Avatar,
			Content:   v.Content,
			Images:    strings.Join(v.Images, ","),
			VideoUrl:  v.VideoUrl,
			VideoImg:  v.VideoImg,
			Location:  v.Location,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
			Likes:     v.Likes,
			IsLike:    v.IsLiked,
		}

		res = append(res, item)
	}
	return res
}

// ServerMomentsCommentsToLocalMomentsComments ， 服务端评论对象转本地评论
// 参数：
//      info ： desc
// 返回值：
//      []*model_struct.LocalMoments ：desc
func ServerMomentsCommentsToLocalMomentsComments(list []*momentsv1.ListItemComment) []*model_struct.LocalMomentsComments {
	var res []*model_struct.LocalMomentsComments
	for _, v := range list {
		item := &model_struct.LocalMomentsComments{
			CreatedAt:      v.CreatedAt,
			CommentId:      v.CommentId,
			MomentId:       v.MomentId,
			UserId:         v.UserId,
			Type:           int(v.Type),
			Avatar:         v.Avatar,
			NickName:       v.Nickname,
			SourceUserId:   v.SourceUserId,
			SourceAvatar:   v.SourceAvatar,
			SourceNickName: v.SourceNickname,
			Content:        v.Content,
		}

		res = append(res, item)
	}

	return res
}
