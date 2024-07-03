package module

import (
	"fmt"

	"github.com/spark-comm/imcloud_sdk/pkg/constant"
	"github.com/spark-comm/imcloud_sdk/pkg/utils"
	"github.com/spark-comm/imcloud_sdk/sdk_struct"

	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
)

type ApiMsgSender struct {
	*MetaManager
}

type SendMsgReq struct {
	RecvID string `json:"recvID" binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	SendMsg
}

type SendMsg struct {
	SendID           string                 `json:"sendID"           binding:"required"`
	GroupID          string                 `json:"groupID"          binding:"required_if=SessionType 2|required_if=SessionType 3"`
	SenderNickname   string                 `json:"senderNickname"`
	SenderFaceURL    string                 `json:"senderFaceURL"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]interface{} `json:"content"          binding:"required"                                            swaggerignore:"true"`
	ContentType      int32                  `json:"contentType"      binding:"required"`
	SessionType      int32                  `json:"sessionType"      binding:"required"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

func (a *ApiMsgSender) SendMsg(sendID, recvID string, index int) error {
	content := fmt.Sprintf("this is test msg user %s to user %s, index: %d", sendID, recvID, index)
	text := sdk_struct.TextElem{Content: content}
	req := &SendMsgReq{
		RecvID: recvID,
		SendMsg: SendMsg{
			SendID:           sendID,
			SenderPlatformID: constant.WindowsPlatformID,
			ContentType:      constant.Text,
			SessionType:      constant.SingleChatType,
			Content:          map[string]interface{}{"content": utils.StructToJsonString(text)},
		},
	}
	var resp msg.SendMsgResp
	return a.postWithCtx(constant.SendMsgRouter, req, &resp)
}
