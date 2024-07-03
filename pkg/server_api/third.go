package server_api

import (
	"github.com/OpenIMSDK/protocol/third"
	"github.com/spark-comm/imcloud_sdk/internal/util"
	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"golang.org/x/net/context"
)

// InitiateMultipartUpload 初始化分片上传
func InitiateMultipartUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	resp := &third.InitiateMultipartUploadResp{}
	err := util.CallPostApi[*third.InitiateMultipartUploadReq, *third.InitiateMultipartUploadResp](
		ctx,
		constant.ObjectInitiateMultipartUpload,
		req,
		resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AuthSign 签名
func AuthSign(ctx context.Context, req *third.AuthSignReq) (*third.AuthSignResp, error) {
	resp := &third.AuthSignResp{}
	err := util.CallPostApi[*third.AuthSignReq, *third.AuthSignResp](ctx, constant.ObjectAuthSign, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CompleteMultipartUpload 完成分片上传
func CompleteMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	resp := &third.CompleteMultipartUploadResp{}
	err := util.CallPostApi[*third.CompleteMultipartUploadReq, *third.CompleteMultipartUploadResp](
		ctx, constant.ObjectCompleteMultipartUpload, req, resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AccessURL  获取访问链接
func AccessURL(ctx context.Context, req *third.AccessURLReq) (*third.AccessURLResp, error) {
	resp := &third.AccessURLResp{}
	err := util.CallPostApi[*third.AccessURLReq, *third.AccessURLResp](
		ctx, constant.ObjectAccessURL, req, resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PartLimit 分片限制
func PartLimit(ctx context.Context) (*third.PartLimitResp, error) {
	resp := &third.PartLimitResp{}
	err := util.CallPostApi[*third.PartLimitReq, *third.PartLimitResp](
		ctx, constant.ObjectPartLimit, &third.PartLimitReq{}, resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
