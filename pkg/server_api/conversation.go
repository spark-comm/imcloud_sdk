package server_api

import (
	pbConversation "github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api/convert"
	"golang.org/x/net/context"
)

// SetConversation 设置会话
func SetConversation(ctx context.Context, apiReq *pbConversation.SetConversationsReq) error {
	_, err := util.ProtoApiPost[pbConversation.SetConversationsReq, empty.Empty](ctx, constant.SetConversationsRouter, apiReq)
	if err != nil {
		return err
	}
	return nil
}

// GetServerConversationList 获取会话列表
func GetServerConversationList(ctx context.Context, loginUserID string) ([]*model_struct.LocalConversation, error) {
	resp := &pbConversation.GetAllConversationsResp{}
	err := util.CallPostApi[*pbConversation.GetAllConversationsReq, *pbConversation.GetAllConversationsResp](
		ctx,
		constant.GetAllConversationsRouter,
		&pbConversation.GetAllConversationsReq{OwnerUserID: loginUserID},
		resp,
	)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerConversationToLocal, resp.Conversations), nil
}

// GetServerConversationsByIDs 获取会话列表
func GetServerConversationsByIDs(ctx context.Context, loginUserID string, conversations []string) ([]*model_struct.LocalConversation, error) {
	// todo 服务端未实现
	resp := &pbConversation.GetConversationsResp{}
	err := util.CallPostApi[*pbConversation.GetConversationIDsReq, *pbConversation.GetConversationsResp](
		ctx,
		constant.GetConversationsRouter,
		&pbConversation.GetConversationIDsReq{UserID: loginUserID},
		resp,
	)
	if err != nil {
		return nil, err
	}
	return util.Batch(convert.ServerConversationToLocal, resp.Conversations), nil
}

// GetServerHasReadAndMaxSeqs 获取会话已读和最大未读消息
func GetServerHasReadAndMaxSeqs(ctx context.Context, loginUserID string, conversationIDs ...string) (map[string]*msg.Seqs, error) {
	resp := &msg.GetConversationsHasReadAndMaxSeqResp{}
	req := msg.GetConversationsHasReadAndMaxSeqReq{UserID: loginUserID}
	req.ConversationIDs = conversationIDs
	resp, err := util.ProtoApiPost[msg.GetConversationsHasReadAndMaxSeqReq, msg.GetConversationsHasReadAndMaxSeqResp](
		ctx,
		constant.GetConversationsHasReadAndMaxSeqRouter,
		&req,
	)
	if err != nil {
		log.ZError(ctx, "getServerHasReadAndMaxSeqs err", err)
		return nil, err
	}
	return resp.Seqs, nil
}

// ClearConversationMsgFromSvr 清除服务端
func ClearConversationMsgFromSvr(ctx context.Context, loginUserID string, conversationID string, isDelOther bool) error {
	var apiReq msg.ClearConversationsMsgReq
	apiReq.UserID = loginUserID
	apiReq.ConversationIDs = []string{conversationID}
	if isDelOther {
		apiReq.DeleteSyncOpt = &msg.DeleteSyncOpt{
			IsSyncOther: isDelOther,
		}
	}
	_, err := util.ProtoApiPost[msg.ClearConversationsMsgReq, empty.Empty](
		ctx,
		constant.ClearConversationMsgRouter,
		&apiReq)
	if err != nil {
		return err
	}
	return nil
}

// SetConversationHasReadSeq 设置会话已读
func SetConversationHasReadSeq(ctx context.Context, loginUserID string, conversationID string, hasReadSeq int64) error {
	req := &msg.SetConversationHasReadSeqReq{UserID: loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	_, err := util.ProtoApiPost[msg.SetConversationHasReadSeqReq, empty.Empty](
		ctx,
		constant.SetConversationHasReadSeq,
		req,
	)
	if err != nil {
		return err
	}
	return nil
}
