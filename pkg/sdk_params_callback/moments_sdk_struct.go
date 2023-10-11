package sdk_params_callback

import "open_im_sdk/pkg/db/model_struct"

// PublishRequest ， 发布朋友圈请求参数
type PublishRequest struct {
	UserName  string   `json:"userName,required"` // 用户名称
	Avatar    string   `json:"avatar"`            // 头像
	Content   string   `json:"content"`           // 内容
	Images    []string `json:"images"`            // 图片
	VideoUrl  string   `json:"videoUrl"`          // 视屏地址
	VideoImg  string   `json:"videoImg"`          // 视屏图片
	Location  string   `json:"location"`          // 位置
	Longitude float64  `json:"longitude"`         // 经度
	Latitude  float64  `json:"latitude"`          // 纬度
}

// CommentRequest ， 评论请求参数
type CommentRequest struct {
	MomentId       string `json:"momentId,required"` // 朋友圈 id
	Type           int32  `json:"type,required"`     // 评论的类型 1:评论；2回复
	Avatar         string `json:"avatar"`            // 头像
	Nickname       string `json:"nickname,required"` // 昵称
	Content        string `json:"content,required"`  // 内容
	SourceUserId   string `json:"sourceUserId"`      // 来源评论用户头像（仅回复）
	SourceAvatar   string `json:"sourceAvatar"`      // 来源评论用户头像（仅回复）
	SourceNickname string `json:"sourceNickname"`    // 来源频率用户昵称（仅回复）
}

// DeleteRequest ， 删除朋友圈请求参数
type DeleteRequest struct {
	MomentId string `json:"momentId,required"` // 朋友圈 id
}

// LikeRequest ， 点赞朋友圈请求参数
type LikeRequest struct {
	MomentId     string `json:"momentId,required"`     // 朋友圈 id
	UserNickname string `json:"userNickname,required"` // 用户昵称
}

// UnlikeRequest ， 取消点赞朋友圈返回参数
type UnlikeRequest struct {
	MomentId string `json:"momentId,required"` // 朋友圈 id
}

// V2ListRequest ， 朋友圈列表请求参数
type V2ListRequest struct {
	IsSelf bool  `json:"isSelf"`        // 是否为仅查询该用户自己发的
	Page   int32 `json:"page,required"` // 页数
	Size   int32 `json:"size,required"` // 每页大小
}

// V2ListReply ， 朋友圈列表返回参数
type V2ListReply struct {
	model_struct.LocalMoments
	Images   []string                             `json:"images"`
	Comments []*model_struct.LocalMomentsComments `json:"comments"`
}
