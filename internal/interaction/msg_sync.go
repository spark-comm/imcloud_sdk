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

package interaction

import (
	"context"
	"github.com/imCloud/im/pkg/proto/msg"
	"math"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/user"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
	"sync"
	"time"

	"github.com/imCloud/im/pkg/common/log"
	"github.com/imCloud/im/pkg/proto/sdkws"
)

const (
	connectPullNums = 1
	defaultPullNums = 30
	SplitPullMsgNum = 100
)

// The callback synchronization starts. The reconnection ends
type MsgSyncer struct {
	loginUserID            string                // login user ID
	longConnMgr            *LongConnMgr          // long connection manager
	PushSeqCh              chan common.Cmd2Value // channel for receiving push messages and the maximum SEQ number
	conversationCh         chan common.Cmd2Value // storage and session triggering
	msgSync                chan *sdkws.SeqRange  // channel sync message
	syncedMaxSeqs          map[string]int64      // map of the maximum synced SEQ numbers for all group IDs
	synceMinSeqs           map[string]int64      // 消息加载到的最小seq
	conversationInitStatus map[string]bool       // 是否完成了会话的初始化加载
	db                     db_interface.DataBase // data store
	syncTimes              int                   // times of sync
	reinstalled            bool                  //true if the app was uninstalled and reinstalled
	isLoginInit            bool                  // tre is init
	ctx                    context.Context       // context
	friend                 *friend.Friend
	group                  *group.Group
	user                   *user.User
}

// NewMsgSyncer creates a new instance of the message synchronizer.
func NewMsgSyncer(ctx context.Context, conversationCh, PushSeqCh chan common.Cmd2Value,
	loginUserID string, longConnMgr *LongConnMgr, db db_interface.DataBase, syncTimes int, friend *friend.Friend, group *group.Group, user *user.User) (*MsgSyncer, error) {
	m := &MsgSyncer{
		loginUserID:            loginUserID,
		longConnMgr:            longConnMgr,
		PushSeqCh:              PushSeqCh,
		conversationCh:         conversationCh,
		ctx:                    ctx,
		syncedMaxSeqs:          make(map[string]int64),
		synceMinSeqs:           make(map[string]int64),
		conversationInitStatus: make(map[string]bool),
		msgSync:                make(chan *sdkws.SeqRange, 1000),
		db:                     db,
		syncTimes:              syncTimes,
		friend:                 friend,
		group:                  group,
		user:                   user,
		isLoginInit:            true,
	}
	if err := m.loadSeq(ctx); err != nil {
		log.ZError(ctx, "loadSeq err", err)
		return nil, err
	}
	go m.DoListener()
	return m, nil
}

// seq The db reads the data to the memory,set syncedMaxSeqs
// seq数据库将数据读取到内存，设置syncedMaxSeqs
func (m *MsgSyncer) loadSeq(ctx context.Context) error {
	conversationIDList, err := m.db.GetAllConversationIDList(ctx)
	if err != nil {
		log.ZError(ctx, "get conversation id list failed", err)
		return err
	}
	for _, v := range conversationIDList {
		maxSyncedSeq, err := m.db.GetConversationNormalMsgSeq(ctx, v)
		if err != nil {
			log.ZError(ctx, "get group normal seq failed", err, "conversationID", v)
		} else {
			m.syncedMaxSeqs[v] = maxSyncedSeq
		}
	}
	notificationSeqs, err := m.db.GetNotificationAllSeqs(ctx)
	if err != nil {
		log.ZError(ctx, "get notification seq failed", err)
		return err
	}
	for _, notificationSeq := range notificationSeqs {
		m.syncedMaxSeqs[notificationSeq.ConversationID] = notificationSeq.Seq
	}
	log.ZDebug(ctx, "loadSeq", "syncedMaxSeqs", m.syncedMaxSeqs)
	return nil
}

// DoListener Listen to the message pipe of the message synchronizer
// and process received and pushed messages
func (m *MsgSyncer) DoListener() {
	for {
		select {
		case cmd := <-m.PushSeqCh:
			m.handlePushMsgAndEvent(cmd)
		case seqRange := <-m.msgSync:
			m.handleSync(seqRange)
		case <-m.ctx.Done():
			log.ZInfo(m.ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
}

// get seqs need sync interval
// 获取seqs需要同步间隔
func (m *MsgSyncer) getSeqsNeedSync(syncedMaxSeq, maxSeq int64) []int64 {
	var seqs []int64
	for i := syncedMaxSeq + 1; i <= maxSeq; i++ {
		seqs = append(seqs, i)
	}
	return seqs
}

// 处理消息推送事件
func (m *MsgSyncer) handlePushMsgAndEvent(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdConnSuccesss:
		//链接成功处理
		log.ZDebug(cmd.Ctx, "recv long conn mgr connected", "cmd", cmd.Cmd, "value", cmd.Value)
		m.doConnected(cmd.Ctx)
	case constant.CmdPushSeq:
		log.ZDebug(cmd.Ctx, "recv max seqs from long conn mgr, start sync msgs", "cmd", cmd.Cmd, "value", cmd.Value)
		wsSeqResp := cmd.Value.(*sdkws.GetMaxSeqResp)
		//同步消息
		m.compareSeqsAndBatchSync(cmd.Ctx, wsSeqResp.MaxSeqs, defaultPullNums)
	case constant.CmdPushMsg:
		m.doPushMsg(cmd.Ctx, cmd.Value.(*sdkws.PushMessages))
	}
}

// // compareSeqsAndBatchSync 批量同步消息
//
//	func (m *MsgSyncer) compareSeqsAndBatchSync(ctx context.Context, maxSeqToSync map[string]int64, minSeqToSync map[string]int64, pullNums int64) {
//		needSyncSeqMap := make(map[string][2]int64)
//		for conversationID, maxSeq := range maxSeqToSync {
//			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
//				if maxSeq > syncedMaxSeq {
//					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq, maxSeq}
//				}
//			} else {
//				//使用服务端的最小seq
//				syncedMinSeq := int64(0)
//				if minSeqToSync != nil {
//					if i, ok1 := minSeqToSync[conversationID]; ok1 {
//						syncedMinSeq = i
//					}
//				}
//				needSyncSeqMap[conversationID] = [2]int64{syncedMinSeq, maxSeq}
//			}
//		}
//		_ = m.syncAndTriggerMsgs(m.ctx, needSyncSeqMap, pullNums)
//	}
func (m *MsgSyncer) compareSeqsAndBatchSync(ctx context.Context, maxSeqToSync map[string]int64, pullNums int64) {
	needSyncSeqMap := make(map[string][2]int64)
	//when app reinstalled do not pull notifications messages.
	if m.reinstalled {
		notificationsSeqMap := make(map[string]int64)
		messagesSeqMap := make(map[string]int64)
		for conversationID, seq := range maxSeqToSync {
			if IsNotification(conversationID) {
				notificationsSeqMap[conversationID] = seq
			} else {
				messagesSeqMap[conversationID] = seq
			}
		}
		for conversationID, seq := range notificationsSeqMap {
			err := m.db.SetNotificationSeq(ctx, conversationID, seq)
			if err != nil {
				log.ZWarn(ctx, "SetNotificationSeq err", err, "conversationID", conversationID, "seq", seq)
				continue
			} else {
				m.syncedMaxSeqs[conversationID] = seq
			}
		}
		for conversationID, maxSeq := range messagesSeqMap {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
			}
		}
		m.reinstalled = false
	} else {
		for conversationID, maxSeq := range maxSeqToSync {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
			}
		}
	}
	// 同步数据
	_ = m.syncAndTriggerMsgs(m.ctx, needSyncSeqMap, pullNums)
}

// compareSeqsAndSync 比较seq和同步消息
func (m *MsgSyncer) compareSeqsAndSync(maxSeqToSync map[string]int64) {
	for conversationID, maxSeq := range maxSeqToSync {
		if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
			if maxSeq > syncedMaxSeq {
				_ = m.syncAndTriggerMsgs(m.ctx, map[string][2]int64{conversationID: {syncedMaxSeq, maxSeq}}, defaultPullNums)
			}
		} else {
			_ = m.syncAndTriggerMsgs(m.ctx, map[string][2]int64{conversationID: {syncedMaxSeq, maxSeq}}, defaultPullNums)
		}
	}
}

// doPushMsg 处理在在线推送
func (m *MsgSyncer) doPushMsg(ctx context.Context, push *sdkws.PushMessages) {
	log.ZDebug(ctx, "push msgs", "push", push, "syncedMaxSeqs", m.syncedMaxSeqs)
	m.pushTriggerAndSync(ctx, push.Msgs, m.triggerConversation)
	m.pushTriggerAndSync(ctx, push.NotificationMsgs, m.triggerNotification)
}

// 推送和同步消息
func (m *MsgSyncer) pushTriggerAndSync(ctx context.Context, pullMsgs map[string]*sdkws.PullMsgs, triggerFunc func(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error) {
	if len(pullMsgs) == 0 {
		return
	}
	needSyncSeqMap := make(map[string][2]int64)
	var lastSeq int64
	var storageMsgs []*sdkws.MsgData
	for conversationID, msgs := range pullMsgs {
		for _, msg := range msgs.Msgs {
			if msg.Seq == 0 {
				_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: []*sdkws.MsgData{msg}}})
				continue
			}
			lastSeq = msg.Seq
			storageMsgs = append(storageMsgs, msg)
		}
		if lastSeq == m.syncedMaxSeqs[conversationID]+int64(len(storageMsgs)) && lastSeq != 0 {
			log.ZDebug(ctx, "trigger msgs", "msgs", storageMsgs)
			_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: storageMsgs}})
			m.syncedMaxSeqs[conversationID] = lastSeq
		} else if lastSeq != 0 { //为0就是全是通知
			needSyncSeqMap[conversationID] = [2]int64{m.syncedMaxSeqs[conversationID], lastSeq}
		}
	}
	m.syncAndTriggerMsgs(ctx, needSyncSeqMap, defaultPullNums)
}

// Called after successful reconnection to synchronize the latest message
// 在成功重新连接后调用以同步最新消息
func (m *MsgSyncer) doConnected(ctx context.Context) {
	// 同步触发通知
	common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.MsgSyncBegin}, m.conversationCh)
	// 同步数据
	go m.syncBaseInformation(ctx)
	var resp sdkws.GetMaxSeqResp
	// 获取会话的最大和最小seq
	if err := m.longConnMgr.SendReqWaitResp(m.ctx, &sdkws.GetMaxSeqReq{UserID: m.loginUserID}, constant.GetNewestSeq, &resp); err != nil {
		log.ZError(m.ctx, "获取最大的seq", err)
		common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.MsgSyncFailed}, m.conversationCh)
		return
	} else {
		log.ZDebug(m.ctx, "get max seq success", "resp", resp)
	}
	wctx, cancelFunc := context.WithTimeout(ctx, time.Second*60*5)
	defer cancelFunc()
	//根据seq同步消息
	m.compareSeqsAndBatchSync(wctx, resp.MaxSeqs, connectPullNums)
	//最小seq
	m.synceMinSeqs = resp.MinSeqs
	common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.MsgSyncEnd}, m.conversationCh)
}

// IsNotification 是否通知
func IsNotification(conversationID string) bool {
	return strings.HasPrefix(conversationID, "n_")
}

// Fragment synchronization message, seq refresh after successful trigger
func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) error {
	if len(seqMap) > 0 {
		tempSeqMap := make(map[string][2]int64, 50)
		msgNum := 0
		for k, v := range seqMap {
			oneConversationSyncNum := v[1] - v[0]
			if (oneConversationSyncNum/SplitPullMsgNum) > 1 && IsNotification(k) {
				nSeqMap := make(map[string][2]int64, 1)
				nSeqMap[k] = [2]int64{v[0], v[0] + oneConversationSyncNum/2}
				for i := 0; i < 2; i++ {
					resp, err := m.pullMsgBySeqRange(ctx, nSeqMap, syncMsgNum)
					if err != nil {
						log.ZError(ctx, "syncMsgFromSvr err", err, "nSeqMap", nSeqMap)
						return err
					}
					_ = m.triggerConversation(ctx, resp.Msgs)
					_ = m.triggerNotification(ctx, resp.NotificationMsgs)
					for conversationID, seqs := range nSeqMap {
						m.syncedMaxSeqs[conversationID] = seqs[1]
					}
					nSeqMap[k] = [2]int64{v[0] + oneConversationSyncNum/2 + 1, v[1]}
				}
				continue
			}
			tempSeqMap[k] = v
			if oneConversationSyncNum > 0 {
				msgNum += int(oneConversationSyncNum)
			}
			if msgNum >= SplitPullMsgNum {
				resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
				if err != nil {
					log.ZError(ctx, "syncMsgFromSvr err", err, "tempSeqMap", tempSeqMap)
					return err
				}
				_ = m.triggerConversation(ctx, resp.Msgs)
				_ = m.triggerNotification(ctx, resp.NotificationMsgs)
				for conversationID, seqs := range tempSeqMap {
					m.syncedMaxSeqs[conversationID] = seqs[1]
				}
				tempSeqMap = make(map[string][2]int64, 50)
				msgNum = 0
			}
		}

		resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
		if err != nil {
			log.ZError(ctx, "syncMsgFromSvr err", err, "seqMap", seqMap)
			return err
		}
		_ = m.triggerConversation(ctx, resp.Msgs)
		_ = m.triggerNotification(ctx, resp.NotificationMsgs)
		for conversationID, seqs := range seqMap {
			m.syncedMaxSeqs[conversationID] = seqs[1]
		}
	}
	return nil
}

// Fragment synchronization message, seq refresh after successful trigger
// 片段同步消息，成功触发后seq刷新
//func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) error {
//	if len(seqMap) > 0 {
//		resp, err := m.pullMsgBySeqRange(ctx, seqMap, syncMsgNum)
//		if err != nil {
//			log.ZError(ctx, "syncMsgFromSvr err", err, "seqMap", seqMap)
//			return err
//		}
//		_ = m.triggerConversation(ctx, resp.Msgs)
//		_ = m.triggerNotification(ctx, resp.NotificationMsgs)
//		//更新最大wz
//		for conversationID, seqs := range seqMap {
//			m.syncedMaxSeqs[conversationID] = seqs[1]
//			//同步到的最小seq
//			m.synceMinSeqs[conversationID] = seqs[0]
//			//师傅初始化加载完成
//			m.conversationInitStatus[conversationID] = true
//		}
//		return err
//	}
//	return nil
//}

func (m *MsgSyncer) splitSeqs(split int, seqsNeedSync []int64) (splitSeqs [][]int64) {
	if len(seqsNeedSync) <= split {
		splitSeqs = append(splitSeqs, seqsNeedSync)
		return
	}
	for i := 0; i < len(seqsNeedSync); i += split {
		end := i + split
		if end > len(seqsNeedSync) {
			end = len(seqsNeedSync)
		}
		splitSeqs = append(splitSeqs, seqsNeedSync[i:end])
	}
	return
}

// 根据seq区间拉取消息
func (m *MsgSyncer) pullMsgBySeqRange(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) (resp *sdkws.PullMessageBySeqsResp, err error) {
	log.ZDebug(ctx, "pullMsgBySeqRange", "seqMap", seqMap, "syncMsgNum", syncMsgNum)

	req := sdkws.PullMessageBySeqsReq{UserID: m.loginUserID}
	for conversationID, seqs := range seqMap {
		var pullNums = syncMsgNum
		minSeq := seqs[0]
		maxSeq := seqs[1]
		if pullNums > maxSeq-minSeq {
			pullNums = maxSeq - minSeq
		}
		// 计算需要同步的数量
		needSyncNum := maxSeq - minSeq
		if needSyncNum < syncMsgNum {
			req.SeqRanges = append(req.SeqRanges, &sdkws.SeqRange{
				ConversationID: conversationID,
				Begin:          minSeq,
				End:            maxSeq,
				Num:            pullNums,
			})
		} else {
			begin := maxSeq - pullNums
			req.SeqRanges = append(req.SeqRanges, &sdkws.SeqRange{
				ConversationID: conversationID,
				Begin:          begin,
				End:            maxSeq,
				Num:            pullNums,
			})
		}
	}
	resp = &sdkws.PullMessageBySeqsResp{}
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullMsgBySeqList, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// handleSync 处理消息同步
func (m *MsgSyncer) handleSync(seqRange *sdkws.SeqRange) {
	req := sdkws.PullMessageBySeqsReq{UserID: m.loginUserID, SeqRanges: []*sdkws.SeqRange{seqRange}}
	resp := &sdkws.PullMessageBySeqsResp{}
	ctx := ccontext.WithOperationID(m.ctx, utils.OperationIDGenerator())
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullMsgBySeqList, resp); err != nil {
		log.ZError(ctx, "sync message err", err)
	}
	_ = m.triggerConversation(ctx, resp.Msgs)
	_ = m.triggerNotification(ctx, resp.NotificationMsgs)
}

// synchronizes messages by SEQs.
// 根据seq触发同步
func (m *MsgSyncer) syncMsgBySeqs(ctx context.Context, conversationID string, seqsNeedSync []int64) (allMsgs []*sdkws.MsgData, err error) {
	pullMsgReq := sdkws.PullMessageBySeqsReq{}
	pullMsgReq.UserID = m.loginUserID
	split := constant.SplitPullMsgNum
	seqsList := m.splitSeqs(split, seqsNeedSync)
	for i := 0; i < len(seqsList); {
		var pullMsgResp sdkws.PullMessageBySeqsResp
		err := m.longConnMgr.SendReqWaitResp(ctx, &pullMsgReq, constant.PullMsgBySeqList, &pullMsgResp)
		if err != nil {
			log.ZError(ctx, "syncMsgFromSvrSplit err", err, "pullMsgReq", pullMsgReq)
			continue
		}
		i++
		allMsgs = append(allMsgs, pullMsgResp.Msgs[conversationID].Msgs...)
	}

	return allMsgs, nil
}

// triggers a conversation with a new message. 用新消息触发对话
func (m *MsgSyncer) triggerConversation(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	err := common.TriggerCmdNewMsgCome(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationCh)
	if err != nil {
		log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
	}
	log.ZDebug(ctx, "triggerConversation", "msgs", msgs)
	return err
}

func (m *MsgSyncer) triggerNotification(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	err := common.TriggerCmdNotification(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationCh)
	if err != nil {
		log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
	}
	return err
}

// SyncConversationMsg 同步会话消息
func (m *MsgSyncer) SyncConversationMsg(ctx context.Context, conversationID string) error {
	return nil
	////获取当前会话的最大seq
	//maxSeq, err := m.GetConversationMaxSeq(ctx, conversationID)
	//if err != nil {
	//	return err
	//}
	//minSeq := m.synceMinSeqs[conversationID]
	////获取本地消息的seq
	//// todo 拉取不到消息暂时屏蔽
	//seq, err := m.db.GetConversationMessageSeq(ctx, conversationID)
	//if len(seq) == 0 {
	//	seq = []int64{0, maxSeq}
	//}
	//if err == nil && len(seq) > 0 {
	//	//最小左侧补齐
	//	if seq[0] > minSeq {
	//		seq = append([]int64{minSeq}, seq...)
	//	}
	//	if seq[len(seq)-1] < maxSeq {
	//		seq = append(seq, maxSeq)
	//	}
	//	//找出缺省的数据
	//	_, missing := utils.FixNotConnected(seq)
	//	if len(missing) > 0 {
	//		return m.syncMsgBySeqsAndConversation(ctx, conversationID, missing)
	//	}
	//}
	//if maxSeq > 0 {
	//	//有需要同步的数据
	//	spiltList := m.SpiltList(minSeq, maxSeq, constant.SplitPullMsgNum)
	//	// 有需要同步的数据
	//	if len(spiltList) > 0 {
	//		//加入异步队列
	//		for i := 0; i < len(spiltList); i++ {
	//			seqRange := spiltList[i]
	//			m.msgSync <- &sdkws.SeqRange{
	//				ConversationID: conversationID,
	//				Begin:          seqRange[0],
	//				End:            seqRange[1],
	//				Num:            constant.SplitPullMsgNum,
	//			}
	//		}
	//	}
	//}
	return nil
}

// GetConversationMaxSeq 获取最大的seq
func (m *MsgSyncer) GetConversationMaxSeq(ctx context.Context, conversationId string) (int64, error) {
	resp, err := util.ProtoApiPost[msg.GetConversationMaxSeqReq, msg.GetConversationMaxSeqResp](
		ctx,
		constant.GetConversationMaxSeqRouter,
		&msg.GetConversationMaxSeqReq{ConversationID: conversationId},
	)
	if err != nil {
		return 0, err
	}
	return resp.MaxSeq, nil
}

// SpiltList 分片
func (m *MsgSyncer) SpiltList(min, max, size int64) [][]int64 {
	if (max - min) <= size {
		return [][]int64{{min, max}}
	}
	spiltList := make([][]int64, 0)
	mod := math.Ceil(float64((max - min) / size))
	tem := max
	for i := 0; i < int(mod); i++ {
		tmpList := make([]int64, 2)
		tmpList[1] = tem
		start := tem - size
		tmpList[0] = start
		tem = start
		spiltList = append(spiltList, tmpList)
	}
	if max%size > 0 {
		tmpList := []int64{tem, tem - (max % size)}
		spiltList = append(spiltList, tmpList)
	}
	return spiltList
}

// syncBaseInformation 同步基本信息
func (m *MsgSyncer) syncBaseInformation(ctx context.Context) {
	// 同步数据
	var wg sync.WaitGroup
	for _, syncFunc := range []func(c context.Context) error{
		m.user.SyncLoginUserInfo,
		m.friend.SyncFriend,               //同步好友
		m.group.SyncGroup,                 //全量同步群组简单字段
		m.group.SyncUserMemberInfoInGroup, //同步自己在所有群中信息
	} {
		wg.Add(1)
		go func(syncFunc func(c context.Context) error) {
			defer wg.Done()
			_ = syncFunc(ctx)
		}(syncFunc)
	}
	wg.Wait()
	//基础信息同步完成触发通知
	common.TriggerCmdNotification(m.ctx, sdk_struct.CmdNewMsgComeToConversation{SyncFlag: constant.BaseDataSyncComplete}, m.conversationCh)
	log.ZError(ctx, "base info sync success", nil)
}

// syncMsgBySeqsAndConversation 根据seq触发同步
func (m *MsgSyncer) syncMsgBySeqsAndConversation(ctx context.Context, conversationID string, seqsNeedSync []int64) error {
	pullMsgReq := sdkws.PullMessageBySeqsReq{}
	pullMsgReq.UserID = m.loginUserID
	split := constant.SplitPullMsgNum
	seqsList := m.splitSeqs(split, seqsNeedSync)
	allMsgs := make([]*sdkws.MsgData, 0)
	for i := 0; i < len(seqsList); {
		pullMsgReq.SeqRanges = []*sdkws.SeqRange{{
			ConversationID: conversationID,
			Begin:          seqsNeedSync[0],
			End:            seqsNeedSync[len(seqsNeedSync)-1],
			Num:            int64(len(seqsNeedSync)),
		}}
		var pullMsgResp sdkws.PullMessageBySeqsResp
		err := m.longConnMgr.SendReqWaitResp(ctx, &pullMsgReq, constant.PullMsgBySeqList, &pullMsgResp)
		if err != nil {
			log.ZError(ctx, "syncMsgFromSvrSplit err", err, "pullMsgReq", pullMsgReq)
			continue
		}
		i++
		if pullMsgResp.Msgs != nil && pullMsgResp.Msgs[conversationID] != nil {
			allMsgs = append(allMsgs, pullMsgResp.Msgs[conversationID].Msgs...)
		}
	}
	msgs := make(map[string]*sdkws.PullMsgs)
	msgs[conversationID] = &sdkws.PullMsgs{Msgs: allMsgs}
	_ = m.triggerConversation(ctx, msgs)
	return nil
}
