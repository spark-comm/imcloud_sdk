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

package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brian-god/xy-apis/api/common/net/v2"
	"github.com/golang/protobuf/proto"
	"github.com/imCloud/api/common"
	"net/http"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/sdkerrs"
	"reflect"
	"time"
)

//var (
//	BaseURL = ""
//	Token   = ""
//)

type apiResponse struct {
	Code   int             `json:"code"`
	ErrMsg string          `json:"errMsg"`
	Msg    string          `json:"msg"`
	Reason string          `json:"reason"`
	Data   json.RawMessage `json:"data"`
}

//func ApiPost(ctx context.Context, api string, req, resp any) (err error) {
//	operationID, _ := ctx.Value("operationID").(string)
//	if operationID == "" {
//		err := sdkerrs.ErrArgs.Wrap("call api operationID is empty")
//		log.ZError(ctx, "ApiRequest", err, "type", "ctx not set operationID")
//		return err
//	}
//	defer func(start time.Time) {
//		end := time.Now()
//		if err == nil {
//			log.ZDebug(ctx, "CallApi", "api", api, "use", "state", "success", time.Duration(end.UnixNano()-start.UnixNano()))
//		} else {
//			log.ZError(ctx, "CallApi", err, "api", api, "use", "state", "failed", time.Duration(end.UnixNano()-start.UnixNano()))
//		}
//	}(time.Now())
//	reqBody, err := json.Marshal(req)
//	if err != nil {
//		log.ZError(ctx, "ApiRequest", err, "type", "json.Marshal(req) failed")
//		return sdkerrs.ErrSdkInternal.Wrap("json.Marshal(req) failed " + err.Error())
//	}
//	ctxInfo := ccontext.Info(ctx)
//	reqUrl := ctxInfo.ApiAddr() + api
//	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(reqBody))
//	if err != nil {
//		log.ZError(ctx, "ApiRequest", err, "type", "http.NewRequestWithContext failed")
//		return sdkerrs.ErrSdkInternal.Wrap("sdk http.NewRequestWithContext failed " + err.Error())
//	}
//	log.ZDebug(ctx, "ApiRequest", "url", reqUrl, "body", string(reqBody))
//	request.ContentLength = int64(len(reqBody))
//	request.Header.Set("Content-Type", "application/json")
//	request.Header.Set("operationID", operationID)
//	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctxInfo.Token()))
//	response, err := new(http.Client).Do(request)
//	if err != nil {
//		log.ZError(ctx, "ApiRequest", err, "type", "network error")
//		return sdkerrs.ErrNetwork.Wrap("ApiPost http.Client.Do failed " + err.Error())
//	}
//	defer response.Body.Close()
//	respBody, err := io.ReadAll(response.Body)
//	if err != nil {
//		log.ZError(ctx, "ApiResponse", err, "type", "read body", "status", response.Status)
//		return sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
//	}
//	log.ZDebug(ctx, "ApiResponse", "url", reqUrl, "status", response.Status, "body", string(respBody))
//	var baseApi apiResponse
//	if err := json.Unmarshal(respBody, &baseApi); err != nil {
//		log.ZError(ctx, "ApiResponse", err, "type", "api code parse")
//		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("api %s json.Unmarshal(%q, %T) failed %s", api, string(respBody), &baseApi, err.Error()))
//	}
//	if baseApi.Code != http.StatusOK {
//		err := sdkerrs.New(baseApi.Code, baseApi.Msg, baseApi.Reason)
//		log.ZError(ctx, "ApiResponse", err, "type", "api code error", "msg", baseApi.ErrMsg, "dlt", baseApi.Reason)
//		return err
//	}
//	if resp == nil || len(baseApi.Data) == 0 || string(baseApi.Data) == "null" {
//		return nil
//	}
//	if err := json.Unmarshal(baseApi.Data, resp); err != nil {
//		log.ZError(ctx, "ApiResponse", err, "type", "api data parse", "data", string(baseApi.Data), "bind", fmt.Sprintf("%T", resp))
//		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("json.Unmarshal(%q, %T) failed %s", string(baseApi.Data), resp, err.Error()))
//	}
//	return nil
//}

// 使用proto格式的数据传输
func ProtoApiPost[T any, R any](ctx context.Context, url string, data *T) (*R, error) {
	resp := new(R)
	message, ok := interface{}(data).(proto.Message)
	if !ok {
		return nil, sdkerrs.ErrArgs.Wrap("req data not proto message")
	}
	res, ok := interface{}(resp).(proto.Message)
	if !ok {
		return nil, sdkerrs.ErrArgs.Wrap("response data not proto message")
	}
	marshal, err := proto.Marshal(message)
	if err != nil {
		return nil, sdkerrs.ErrArgs.Wrap("marshal args err")
	}
	reqUrl := ccontext.Info(ctx).ApiAddr() + url
	content, err := network.PostWithTimeOutByte(reqUrl, marshal, ccontext.Info(ctx).Token(), time.Second*300)
	if err != nil {
		return nil, sdkerrs.ErrSdkInternal.Wrap("json.Marshal(req) failed " + err.Error())
	}
	var result v2.Result
	err = proto.Unmarshal(content, &result)
	if err != nil {
		return nil, sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	if result.Code != http.StatusOK {
		return nil, sdkerrs.New(
			int(result.Code),
			result.Msg,
			result.ErrMsg,
			result.Reason)
	}
	if reflect.TypeOf(resp).Elem().Name() == "Empty" {
		return resp, nil
	}
	if len(result.Data) == 0 || string(result.Data) == "null" {
		return resp, nil
	}
	if err := proto.Unmarshal(result.Data, res); err != nil {
		return nil, sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("proto.Unmarshal(%q, %T) failed %s", string(result.Data), resp, err.Error()))
	}
	return interface{}(res).(*R), nil
}

//func CallApi[T any](ctx context.Context, api string, req any) (*T, error) {
//	resp := req.(T)
//	if err := ApiPost(ctx, api, req, &resp); err != nil {
//		return nil, err
//	}
//	return &resp, nil
//}

func CallPostApi[R, T proto.Message](ctx context.Context, api string, req R, res T) error {
	return CallPostProtoApi(
		ctx,
		api,
		req,
		res,
	)
}

// CallPostProtoApi
// 使用proto格式的数据传输
func CallPostProtoApi[T proto.Message, R proto.Message](ctx context.Context, url string, req T, resp R) error {
	log.Error("", fmt.Sprintf("请求接口地址:%s", url))
	marshal, err := proto.Marshal(req)
	if err != nil {
		return sdkerrs.ErrArgs.Wrap("marshal args err")
	}
	reqUrl := ccontext.Info(ctx).ApiAddr() + url
	content, err := network.PostWithTimeOutByte(reqUrl, marshal, ccontext.Info(ctx).Token(), time.Second*300)
	if err != nil {
		return sdkerrs.ErrSdkInternal.Wrap("call api  failed " + err.Error())
	}
	var result v2.Result
	err = proto.Unmarshal(content, &result)
	if err != nil {
		return sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	if result.Code != http.StatusOK {
		return sdkerrs.New(int(result.Code), result.Msg, result.ErrMsg, result.Reason)
	}
	if len(result.Data) == 0 || string(result.Data) == "null" {
		return nil
	}
	if err := proto.Unmarshal(result.Data, resp); err != nil {
		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("proto.Unmarshal(%q, %T) failed %s", string(result.Data), resp, err.Error()))
	}
	return nil
}

// GetPageAll 获取所有配置
//
//	func GetPageAll[A interface {
//		GetPagination() *sdkws.RequestPagination
//	}, B, C any](ctx context.Context, api string, req A, fn func(resp *B) []C) ([]C, error) {
//
//		if req.GetPagination().ShowNumber <= 0 {
//			req.GetPagination().ShowNumber = 50
//		}
//		var res []C
//		for i := int32(0); ; i++ {
//			req.GetPagination().PageNumber = i + 1
//			memberResp, err := CallApi[B](ctx, api, req)
//			if err != nil {
//				return nil, err
//			}
//			list := fn(memberResp)
//			res = append(res, list...)
//			if len(list) < int(req.GetPagination().ShowNumber) {
//				break
//			}
//		}
//		return res, nil
//	}
func GetPageAll[A interface {
	GetPagination() *common.RequestPagination
	proto.Message
}, B, C proto.Message](ctx context.Context, api string, req A, resp B, fn func(resp B) []C) ([]C, error) {
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 100
	}
	var res []C
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
		err := CallPostApi[A, B](ctx, api, req, resp)
		if err != nil {
			return nil, err
		}
		list := fn(resp)
		res = append(res, list...)
		if len(list) < int(req.GetPagination().ShowNumber) {
			break
		}
	}
	return res, nil
}

// GetFirstPage 获取第一页数据
//
//	func GetFirstPage[A interface {
//		GetPagination() *sdkws.RequestPagination
//	}, B, C any](ctx context.Context, api string, req A, fn func(resp *B) []C) ([]C, error) {
//
//		if req.GetPagination().ShowNumber <= 0 {
//			req.GetPagination().ShowNumber = 10
//		}
//		var res []C
//		memberResp, err := CallApi[B](ctx, api, req)
//		if err != nil {
//			return nil, err
//		}
//		list := fn(memberResp)
//		res = append(res, list...)
//		return res, nil
//	}
func GetFirstPage[A interface {
	GetPagination() *common.RequestPagination
	proto.Message
}, B, C proto.Message](ctx context.Context, api string, req A, resp B, fn func(resp B) []C) ([]C, error) {
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 10
	}
	var res []C
	err := CallPostApi[A, B](ctx, api, req, resp)
	if err != nil {
		return nil, err
	}
	list := fn(resp)
	res = append(res, list...)
	return res, nil
}
