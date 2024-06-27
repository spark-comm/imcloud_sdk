package server_api

import (
	"github.com/golang/protobuf/ptypes/empty"
	groupmodel "github.com/miliao_apis/api/common/model/group/v2"
	netmodel "github.com/miliao_apis/api/common/net/v2"
	v2 "github.com/miliao_apis/api/im_cloud/group/v2"
	"github.com/openimsdk/openim-sdk-core/internal/util"
	"github.com/openimsdk/openim-sdk-core/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/pkg/server_api/convert"
	"golang.org/x/net/context"
)

// GetServerAdminGroupApplicationList 获取服务端加群申请
func GetServerAdminGroupApplicationList(ctx context.Context, loginUserID string) ([]*model_struct.LocalAdminGroupRequest, error) {
	fn := func(resp *v2.GetRecvGroupApplicationListReply) []*groupmodel.GroupRequestInfo {
		return resp.List
	}
	req := &netmodel.GetByFormUserListSdk{FromUserID: loginUserID, Pagination: &netmodel.RequestPagination{}}
	resp := &v2.GetRecvGroupApplicationListReply{}
	list, err := util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerGroupRequestToLocalAdminGroupRequest, list), nil
}

// GetServerSelfGroupApplication 获取服务端的自己的加群请求
func GetServerSelfGroupApplication(ctx context.Context, loginUserID string) ([]*model_struct.LocalGroupRequest, error) {
	fn := func(resp *v2.GetRecvGroupApplicationListReply) []*groupmodel.GroupRequestInfo {
		return resp.List
	}
	req := &netmodel.GetByUserListSdk{UserID: loginUserID, Pagination: &netmodel.RequestPagination{}}
	resp := &v2.GetRecvGroupApplicationListReply{}
	list, err := util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, resp, fn)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerGroupRequestToLocalGroupRequest, list), nil
}

// HandlerGroupApplication 处理群申请
func HandlerGroupApplication(ctx context.Context, req *v2.ApplicationResponseReq) error {
	if _, err := util.ProtoApiPost[v2.ApplicationResponseReq, empty.Empty](
		ctx,
		constant.AcceptGroupApplicationRouter,
		req,
	); err != nil {
		return err
	}
	return nil
}
