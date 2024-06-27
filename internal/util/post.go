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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/brian-god/imcloud_sdk/pkg/ccontext"
	"github.com/brian-god/imcloud_sdk/pkg/network"
	"github.com/brian-god/imcloud_sdk/pkg/sdkerrs"
	"github.com/golang/protobuf/proto"
	v2 "github.com/miliao_apis/api/common/net/v2"
	"golang.org/x/net/context"
)

// apiClient is a global HTTP client with a timeout of one minute.
var apiClient = &http.Client{
	Timeout: time.Second * 30,
}

// ApiResponse represents the standard structure of an API response.
type ApiResponse struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	ErrDlt  string          `json:"errDlt"`
	Data    json.RawMessage `json:"data"`
}

// ProtoApiPost 使用proto格式的数据传输
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
	contextInfo := ccontext.Info(ctx)
	reqUrl := contextInfo.ApiAddr() + url
	content, err := network.PostWithTimeOutByte(reqUrl, marshal, contextInfo.Token(), contextInfo.Language(), time.Second*300)
	if err != nil {
		return nil, sdkerrs.ErrSdkInternal.Wrap("json.Marshal(req) failed " + err.Error())
	}
	var result v2.Result
	err = proto.Unmarshal(content, &result)
	if err != nil {
		return nil, sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	if result.Code != http.StatusOK {
		return nil, sdkerrs.New(int(result.Code), result.Msg, result.ErrMsg)
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
	//log.Error("", fmt.Sprintf("请求接口地址:%s", url))
	marshal, err := proto.Marshal(req)
	if err != nil {
		return sdkerrs.ErrArgs.Wrap("marshal args err")
	}
	contextInfo := ccontext.Info(ctx)
	reqUrl := contextInfo.ApiAddr() + url
	content, err := network.PostWithTimeOutByte(reqUrl, marshal, contextInfo.Token(), contextInfo.Language(), time.Second*300)
	if err != nil {
		return sdkerrs.ErrSdkInternal.Wrap("call api  failed " + err.Error())
	}
	var result v2.Result
	err = proto.Unmarshal(content, &result)
	if err != nil {
		return sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	if result.Code != http.StatusOK {
		return sdkerrs.New(int(result.Code), result.Msg, result.ErrMsg)
	}
	if len(result.Data) == 0 || string(result.Data) == "null" {
		return nil
	}
	if err := proto.Unmarshal(result.Data, resp); err != nil {
		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("proto.Unmarshal(%q, %T) failed %s", string(result.Data), resp, err.Error()))
	}
	return nil
}
func GetPageAll[A interface {
	GetPagination() *v2.RequestPagination
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

func GetFirstPage[A interface {
	GetPagination() *v2.RequestPagination
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
