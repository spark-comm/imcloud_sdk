package server_api

import (
	"github.com/OpenIMSDK/protocol/msg"
	pbMsg "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openimsdk/openim-sdk-core/internal/util"
	"github.com/openimsdk/openim-sdk-core/pkg/constant"
	"golang.org/x/net/context"
)

// DeleteAllMessageFromSvr Delete all server messages
func DeleteAllMessageFromSvr(ctx context.Context, loginUserID string) error {
	var apiReq msg.UserClearAllMsgReq
	apiReq.UserID = loginUserID
	_, err := util.ProtoApiPost[msg.UserClearAllMsgReq, empty.Empty](
		ctx,
		constant.ClearAllMsgRouter,
		&apiReq,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMessageFromSvr The user deletes part of the message from the server
func DeleteMessageFromSvr(ctx context.Context, loginUserID string, conversationID string, seqs ...int64) error {
	_, err := util.ProtoApiPost[msg.DeleteMsgsReq, empty.Empty](
		ctx,
		constant.DeleteMsgsRouter,
		&msg.DeleteMsgsReq{
			UserID:         loginUserID,
			Seqs:           seqs,
			ConversationID: conversationID,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// PullMessageBySeqs 根据seq拉取消息
func PullMessageBySeqs(ctx context.Context, loginUserID string, seqs []*sdkws.SeqRange) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{}
	err := util.CallPostApi[*sdkws.PullMessageBySeqsReq, *sdkws.PullMessageBySeqsResp](
		ctx,
		constant.PullUserMsgBySeqRouter,
		&sdkws.PullMessageBySeqsReq{UserID: loginUserID, SeqRanges: seqs},
		resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RevokeOneMessage revoke one message
func RevokeOneMessage(ctx context.Context, conversationID, loginUserID string, seq int64) error {
	if _, err := util.ProtoApiPost[msg.RevokeMsgReq, empty.Empty](
		ctx,
		constant.RevokeMsgRouter,
		&msg.RevokeMsgReq{
			ConversationID: conversationID,
			Seq:            seq,
			UserID:         loginUserID},
	); err != nil {
		return err
	}
	return nil
}

// MarkMsgAsRead2Svr mark msg as read
func MarkMsgAsRead2Svr(ctx context.Context, loginUserID, conversationID string, seqs []int64) error {
	req := &pbMsg.MarkMsgsAsReadReq{UserID: loginUserID, ConversationID: conversationID, Seqs: seqs}
	_, err := util.ProtoApiPost[pbMsg.MarkMsgsAsReadReq, empty.Empty](
		ctx,
		constant.MarkMsgsAsReadRouter,
		req,
	)
	if err != nil {
		return err
	}
	return nil
}

// MarkConversationAsReadSvr mark conversation as read
func MarkConversationAsReadSvr(ctx context.Context, loginUserID, conversationID string, hasReadSeq int64, seqs []int64) error {
	req := &pbMsg.MarkConversationAsReadReq{UserID: loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq, Seqs: seqs}

	_, err := util.ProtoApiPost[pbMsg.MarkConversationAsReadReq, empty.Empty](
		ctx,
		constant.MarkConversationAsRead,
		req,
	)
	if err != nil {
		return err
	}
	return nil
}
