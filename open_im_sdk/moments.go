package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

// GetMomentsList ， 获取朋友圈列表
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func GetMomentsList(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().GetMomentsList, params)
}

// PublishMoments ， 发布朋友圈
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func PublishMoments(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().Publish, params)
}

// DeleteMoments ， 删除朋友圈
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func DeleteMoments(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().Delete, params)
}

// CommentMoments ， 评论朋友圈
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func CommentMoments(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().Comment, params)
}

// LikeMoments ， 点赞朋友圈
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func LikeMoments(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().Like, params)
}

// UnLikeMoments ， 取消点赞朋友圈
// 参数：
//      callback ： desc
//      operationID ： desc
//      params ： desc
// 返回值：
func UnLikeMoments(callback open_im_sdk_callback.Base, operationID string, params string) {
	call(callback, operationID, UserForSDK.Moments().UnLike, params)
}
