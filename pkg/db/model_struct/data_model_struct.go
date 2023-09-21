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

package model_struct

//
//message FriendInfo{
//string OwnerUserID = 1;
//string Remark = 2;
//int64 CreateTime = 3;
//UserInfo FriendUser = 4;
//int32 AddSource = 5;
//string OperatorUserID = 6;
//string Ex = 7;
//}
//open_im_sdk.FriendInfo(FriendUser) != imdb.Friend(FriendUserID)
//	table = ` CREATE TABLE IF NOT EXISTS friends(
//     owner_user_id CHAR (64) NOT NULL,
//     friend_user_id CHAR (64) NOT NULL ,
//     name varchar(64) DEFAULT NULL ,
//	 face_url varchar(100) DEFAULT NULL ,
//     remark varchar(255) DEFAULT NULL,
//     gender int DEFAULT NULL ,
//   	 phone_number varchar(32) DEFAULT NULL ,
//	 birth INTEGER DEFAULT NULL ,
//	 email varchar(64) DEFAULT NULL ,
//	 create_time INTEGER DEFAULT NULL ,
//	 add_source int DEFAULT NULL ,
//	 operator_user_id CHAR(64) DEFAULT NULL,
//  	 ex varchar(1024) DEFAULT NULL,
//  	 PRIMARY KEY (owner_user_id,friend_user_id)
// 	)`

type LocalFriend struct {
	OwnerUserID    string `json:"ownerUserID" gorm:"column:owner_user_id;primary_key;type:varchar(64);comment:所属用户"`
	FriendUserID   string `json:"friendUserID" gorm:"column:friend_user_id;primary_key;type:varchar(64);comment:好友id"`
	Remark         string `json:"remark" gorm:"column:remark;type:varchar(255);comment:备注"`
	CreateAt       int64  `json:"createdAt" gorm:"column:create_at;comment:创建时间"`
	AddSource      int32  `json:"addSource" gorm:"column:add_source;comment:添加来源"`
	OperatorUserID string `json:"operatorUserID" gorm:"column:operator_user_id;type:varchar(64)"`
	FaceURL        string `json:"faceURL" gorm:"column:face_url;size:44;comment:头像"`
	Nickname       string `json:"nickname" gorm:"column:nickname;size:50;comment:用户昵称"`
	Message        string `json:"message" gorm:"column:message;size:50;comment:个性签名"`
	Code           string `json:"code" gorm:"column:code;size:6;uniqueIndex;comment:用户ID"`
	Phone          string `json:"phone" gorm:"column:phone;size:16;uniqueIndex;comment:手机号码"`
	Email          string `json:"email" gorm:"column:email;size:36;comment:邮箱"`
	Birth          int64  `json:"birth" gorm:"column:birth;comment:生日"`
	Gender         int32  `json:"gender" gorm:"column:gender;default:0;comment:性别，默认女"`
	Ex             string `json:"ex" gorm:"column:ex;type:varchar(1024)"`
	BackgroundURL  string `json:"backgroundURL" gorm:"column:background_url;type:varchar(255)"`
	AttachedInfo   string `json:"attachedInfo" gorm:"column:attached_info;type:varchar(1024)"`
	SortFlag       string `gorm:"column:sort_flag;type:varchar(2);comment:排序标志" json:"sortFlag"`
	NotPeersFriend int32  `json:"notPeersFriend" gorm:"column:not_peers_friend;default:0;comment:关系链是否完整；1303断裂"`
	IsComplete     int32  `gorm:"column:is_complete;default:1;comment:同步完成" json:"isComplete"`
	IsDestroyMsg   int32  `gorm:"column:is_destroy_msg;not null;tinyint;default:0;comment:阅后即焚开关0-关闭 1-开启" json:"isDestroyMsg" `
}

// message FriendRequest{
// string  FromUserID = 1;
// string ToUserID = 2;
// int32 HandleResult = 3;
// string ReqMsg = 4;
// int64 CreateTime = 5;
// string HandlerUserID = 6;
// string HandleMsg = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }
// open_im_sdk.FriendRequest == imdb.FriendRequest
type LocalFriendRequest struct {
	FromUserID    string `gorm:"column:from_user_id;primary_key;type:varchar(64);comment:来源用户ID" json:"fromUserID"`
	FromNickname  string `gorm:"column:from_nickname;type:varchar;type:varchar(255);comment:来源用户昵称" json:"fromNickname"`
	FromFaceURL   string `gorm:"column:from_face_url;type:varchar;type:varchar(255);comment:来源用户头像" json:"fromFaceURL"`
	FromGender    int32  `gorm:"column:from_gender;comment:来源用户性别" json:"fromGender"`
	FromCode      string `json:"fromCode" gorm:"column:from_code;size:6;comment:来源用户code"`
	FromMessage   string `json:"fromMessage" gorm:"column:from_message;size:6;comment:来源用户个性签名"`
	FromPhone     string `json:"fromPhone" gorm:"column:from_phone;size:6;comment:来源用户手机号"`
	ToUserID      string `gorm:"column:to_user_id;primary_key;type:varchar(64);comment:接收方用户ID" json:"toUserID"`
	ToNickname    string `gorm:"column:to_nickname;type:varchar;type:varchar(255);comment:接收方昵称" json:"toNickname"`
	ToFaceURL     string `gorm:"column:to_face_url;type:varchar;type:varchar(255);comment:接收方头像" json:"toFaceURL"`
	ToGender      int32  `gorm:"column:to_gender;comment:接收方性别" json:"toGender"`
	ToCode        string `json:"toCode" gorm:"column:to_code;size:6;comment:接收方的code"`
	ToMessage     string `json:"toMessage" gorm:"column:to_message;size:6;comment:接收方个性签名"`
	TomPhone      string `json:"toPhone" gorm:"column:tom_phone;size:6;comment:接收方手机号"`
	HandleResult  int32  `gorm:"column:handle_result;comment:处理0:待处理;1:同意;2:拒绝" json:"handleResult"`
	ReqMsg        string `gorm:"column:req_msg;type:varchar(255);comment:请求消息" json:"reqMsg"`
	CreateAt      int64  `gorm:"column:create_at;comment:创建时间" json:"createAt"`
	HandlerUserID string `gorm:"column:handler_user_id;type:varchar(64);comment:处理的用户id" json:"handlerUserID"`
	HandleMsg     string `gorm:"column:handle_msg;type:varchar(255);comment:处理备注" json:"handleMsg"`
	HandleTime    int64  `gorm:"column:handle_time;comment:处理时间" json:"handleTime"`
	SortFlag      string `gorm:"column:sort_flag;type:varchar(2);comment:排序标志" json:"sortFlag"`
	Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo  string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

//message GroupInfo{
//  string GroupID = 1;
//  string GroupName = 2;
//  string NotificationCmd = 3;
//  string Introduction = 4;
//  string FaceUrl = 5;
//  string OwnerUserID = 6;
//  uint32 MemberCount = 8;
//  int64 CreateTime = 7;
//  string Ex = 9;
//  int32 Status = 10;
//  string CreatorUserID = 11;
//  int32 GroupType = 12;
//}
//  open_im_sdk.GroupInfo (OwnerUserID ,  MemberCount )> imdb.Group
//    	group_id char(64) NOT NULL,
//		name varchar(64) DEFAULT NULL ,
//    	introduction varchar(255) DEFAULT NULL,
//    	notification varchar(255) DEFAULT NULL,
//    	face_url varchar(100) DEFAULT NULL,
//    	group_type int DEFAULT NULL,
//    	status int DEFAULT NULL,
//    	creator_user_id char(64) DEFAULT NULL,
//    	create_time INTEGER DEFAULT NULL,
//    	ex varchar(1024) DEFAULT NULL,
//    	PRIMARY KEY (group_id)
//	)`

const (
	FinishMemberSync = 2
)

type LocalGroup struct {
	GroupID                string `json:"groupID" gorm:"column:group_id;primary_key;type:varchar(64)"`
	CreateTime             int64  `json:"createTime" gorm:"column:create_time"`
	Ex                     string `json:"ex" gorm:"column:ex;type:varchar(1024)"`
	Code                   string `json:"code" gorm:"column:code;size:20;comment:群编码,对应前台的群id"`
	GroupName              string `json:"groupName" gorm:"column:name;size:255;comment:群名称" `
	Notification           string `json:"notification" gorm:"column:notification;size:255;comment:群公告" `
	Introduction           string `json:"introduction" gorm:"column:introduction;size:2000;comment:群简介"`
	FaceURL                string `json:"faceURL" gorm:"column:face_url;size:255;comment:群头像"`
	Status                 int32  `json:"status" gorm:"column:status;default:0;comment:群状态"`
	CreatorUserID          string `json:"creatorUserID" gorm:"column:creator_user_id;size:64;comment:建群用户Id"`
	GroupType              int32  `json:"groupType" gorm:"column:group_type;default:0;comment:群类型"`
	OwnerUserID            string `json:"ownerUserID" gorm:"column:owner_user_id;size:64;comment:群主id"`
	MaxMemberCount         int32  `json:"maxMemberCount" gorm:"column:max_member_count;comment:群的最大用户数"`
	MemberCount            int32  `json:"memberCount" gorm:"column:member_count;comment:群用户数"`
	AttachedInfo           string `json:"attachedInfo" gorm:"column:attached_info;size:1024;comment:附加信息"`
	NeedVerification       int32  `json:"needVerification" gorm:"column:need_verification;default:0;comment:加群是否需要验证:0需要,1不需要"`
	LookMemberInfo         int32  `json:"lookMemberInfo" gorm:"column:look_member_info;default:0;comment:支持成员信息查看:0支持,1不支持"`
	ApplyMemberFriend      int32  `json:"applyMemberFriend" gorm:"column:apply_member_friend;default:0;comment:是否可通过群加好友:0支持,1不支持" `
	OnlyManageUpdateName   int32  `json:"onlyManageUpdateName" gorm:"column:only_manage_update_name;default:0;comment:是否仅管理员或群主能够更新群名称:"`
	NotificationUpdateTime int64  `json:"notificationUpdateTime" gorm:"column:notification_update_time;comment:群公告更新时间"`
	NotificationUserID     string `json:"notificationUserID" gorm:"column:notification_user_id;size:64;comment:更新公告用户id"`
	IsReal                 int32  `json:"isReal" gorm:"column:is_real;default:0;comment:是否实名认证，默认否"`
	IsOpen                 uint   `json:"isOpen" gorm:"column:is_open;default:1;comment:是否公开群组,1:开启;2:关闭"`
	AllowPrivateChat       uint   `json:"allowPrivateChat" gorm:"column:allow_private_chat;default:2;comment:允许成员私聊,1:开启;2:关闭"`
	IsComplete             int32  `gorm:"column:is_complete;default:1;comment:同步完成" json:"isComplete"`
	MemberIsComplete       int32  `gorm:"column:member_is_complete;default:0;comment:成员同步完成" json:"member_is_complete"`
}

//message GroupMemberFullInfo {
//string GroupID = 1 ;
//string UserID = 2 ;
//int32 roleLevel = 3;
//int64 JoinTime = 4;
//string NickName = 5;
//string FaceUrl = 6;
//int32 JoinSource = 8;
//string OperatorUserID = 9;
//string Ex = 10;
//int32 AppMangerLevel = 7; //if >0
//}  open_im_sdk.GroupMemberFullInfo(AppMangerLevel) > imdb.GroupMember
//  group_id char(64) NOT NULL,
//   user_id char(64) NOT NULL,
//   nickname varchar(64) DEFAULT NULL,
//   user_group_face_url varchar(64) DEFAULT NULL,
//   role_level int DEFAULT NULL,
//   join_time INTEGER DEFAULT NULL,
//   join_source int DEFAULT NULL,
//   operator_user_id char(64) NOT NULL,

type LocalGroupMember struct {
	GroupID        string `gorm:"column:group_id;primary_key;type:varchar(64);comment:群id" json:"groupID"`
	UserID         string `gorm:"column:user_id;primary_key;type:varchar(64);comment:用户id" json:"userID"`
	FaceURL        string `gorm:"column:face_url;size:255;comment:头像" json:"faceURL"`
	Nickname       string `gorm:"column:nickname;size:255;comment:用户昵称" json:"nickname"`
	SortFlag       string `gorm:"column:sort_flag;size:10;comment:用户昵称排序" json:"sortFlag"`
	GroupUserName  string `gorm:"column:group_user_name;size:255;comment:用户群中昵称" json:"groupUserName"`
	Message        string `gorm:"column:message;size:255;comment:个性签名" json:"message"`
	Code           string `gorm:"column:code;size:6;comment:用户ID" json:"code"`
	Phone          string `gorm:"column:phone;size:16;comment:手机号码" json:"phone"`
	Email          string `gorm:"column:email;size:36;comment:邮箱" json:"email"`
	Birth          int64  `gorm:"column:birth;comment:生日" json:"birth"`
	Gender         int32  `gorm:"column:gender;default:0;comment:性别，默认女" json:"gender"`
	RoleLevel      int32  `gorm:"column:role_level;comment:角色20:普通用户100:群主；60:管理员" json:"roleLevel"`
	JoinTime       int64  `gorm:"column:join_time;index:index_join_tim;comment:加入时间" json:"joinTime"`
	JoinSource     int32  `gorm:"column:join_source;comment:来源" json:"joinSource"`
	InviterUserID  string `gorm:"column:inviter_user_id;size:64;comment:邀请用户id"  json:"inviterUserID"`
	MuteEndTime    int64  `gorm:"column:mute_end_time;default:0;comment:禁言时间" json:"muteEndTime"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	BackgroundURL  string `json:"backgroundURL" gorm:"column:background_url;type:varchar(255)"`
}

// message GroupRequest{
// string UserID = 1;
// string GroupID = 2;
// string HandleResult = 3;
// string ReqMsg = 4;
// string  HandleMsg = 5;
// int64 ReqTime = 6;
// string HandleUserID = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }open_im_sdk.GroupRequest == imdb.GroupRequest
type LocalGroupRequest struct {
	GroupID       string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID"`
	GroupName     string `gorm:"column:group_name;size:255" json:"groupName"`
	Notification  string `gorm:"column:notification;type:varchar(255)" json:"notification"`
	Introduction  string `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	GroupFaceURL  string `gorm:"column:face_url;type:varchar(255)" json:"groupFaceURL"`
	CreateTime    int64  `gorm:"column:create_time" json:"createTime"`
	Status        int32  `gorm:"column:status" json:"status"`
	CreatorUserID string `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	GroupType     int32  `gorm:"column:group_type" json:"groupType"`
	OwnerUserID   string `gorm:"column:owner_user_id;type:varchar(64)" json:"ownerUserID"`
	MemberCount   int32  `gorm:"column:member_count" json:"memberCount"`
	GroupCode     string `gorm:"column:group_code" json:"groupCode"`

	UserID      string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname    string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	UserFaceURL string `gorm:"column:user_face_url;type:varchar(255)" json:"userFaceURL"`
	Gender      int32  `gorm:"column:gender;type:tinyint" json:"gender"`
	Code        string `gorm:"column:code;type:varchar(32)" json:"code"`

	HandleResult int32  `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg       string `gorm:"column:req_msg;type:varchar(255)" json:"reqMsg"`
	HandledMsg   string `gorm:"column:handle_msg;type:varchar(255)" json:"handledMsg"`
	ReqTime      int64  `gorm:"column:req_time" json:"reqTime"`
	HandleUserID string `gorm:"column:handle_user_id;type:varchar(64)" json:"handleUserID"`
	HandledTime  int64  `gorm:"column:handle_time" json:"handledTime"`
	//Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	//AttachedInfo  string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	JoinSource    int32  `gorm:"column:join_source" json:"joinSource"`
	InviterUserID string `gorm:"column:inviter_user_id;size:64"  json:"inviterUserID"`
}

// string UserID = 1;
// string Nickname = 2;
// string FaceUrl = 3;
// int32 Gender = 4;
// string PhoneNumber = 5;
// string Birth = 6;
// string Email = 7;
// string Ex = 8;
// int64 CreateTime = 9;
// int32 AppMangerLevel = 10;
// open_im_sdk.User == imdb.User
type LocalUser struct {
	UserID           string `json:"userID" gorm:"column:user_id;primary_key;type:varchar(64)"`
	FaceURL          string `json:"faceURL" gorm:"column:face_url;size:255;comment:头像文件名"`
	Nickname         string `json:"nickname" gorm:"column:nickname;size:50;comment:用户昵称"`
	Message          string `json:"message" gorm:"column:message;size:50;comment:个性签名"`
	Code             string `json:"code" gorm:"column:code;size:6;uniqueIndex;comment:用户ID"`
	Phone            string `json:"phone" gorm:"column:phone;size:16;uniqueIndex;comment:手机号码"`
	Email            string `json:"email" gorm:"column:email;size:36;comment:邮箱"`
	Birth            int64  `json:"birth" gorm:"column:birth;comment:生日"`
	Gender           int32  `gorm:"column:gender;comment:性别" json:"gender"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	ShareCode        string `json:"shareCode" gorm:"column:share_code;size:20;not null;default:'';comment:分享码"`
	LastLogin        int64  `json:"lastLogin" gorm:"column:last_login;default:0;comment:上次登陆时间"`
	Options          string `json:"options" gorm:"column:optionst;default:'';size:3000;comment:用户配置项"`
	AppMangerLevel   int32  `gorm:"column:app_manger_level" json:"-"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	GlobalRecvMsgOpt int32  `gorm:"column:global_recv_msg_opt" json:"globalRecvMsgOpt"`
}

// message BlackInfo{
// string OwnerUserID = 1;
// int64 CreateTime = 2;
// PublicUserInfo BlackUserInfo = 4;
// int32 AddSource = 5;
// string OperatorUserID = 6;
// string Ex = 7;
// }
// open_im_sdk.BlackInfo(BlackUserInfo) != imdb.Black (BlockUserID)
type LocalBlack struct {
	OwnerUserID    string `gorm:"column:owner_user_id;primary_key;type:varchar(64)" json:"ownerUserID"`
	BlackUserID    string `gorm:"column:black_user_id;primary_key;type:varchar(64)" json:"blackUserID"`
	Nickname       string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	Gender         int32  `gorm:"column:gender;comment:性别" json:"gender"`
	Message        string `json:"message" gorm:"column:message;size:50;comment:个性签名"`
	Code           string `json:"code" gorm:"column:code;size:6;comment:用户ID"`
	CreateAt       int64  `gorm:"column:createAt" json:"createAt"`
	AddSource      int32  `gorm:"column:add_source" json:"addSource"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	SortFlag       string `gorm:"column:sort_flag;type:varchar(2);comment:排序标志" json:"sortFlag"`
}

type LocalSeqData struct {
	UserID string `gorm:"column:user_id;primary_key;type:varchar(64)"`
	Seq    uint32 `gorm:"column:seq"`
}

type LocalSeq struct {
	ID     string `gorm:"column:id;primary_key;type:varchar(64)"`
	MinSeq uint32 `gorm:"column:min_seq"`
}

// `create table if not exists  chat_log (
//
//	     client_msg_id char(64)   NOT NULL,
//	     server_msg_id char(64)   DEFAULT NULL,
//		  send_id char(64)   NOT NULL ,
//		  is_read int NOT NULL ,
//		  seq INTEGER DEFAULT NULL ,
//		  status int NOT NULL ,
//		  session_type int NOT NULL ,
//		  recv_id char(64)   NOT NULL ,
//		  content_type int NOT NULL ,
//	     sender_face_url varchar(100) DEFAULT NULL,
//	     sender_nick_name varchar(64) DEFAULT NULL,
//		  msg_from int NOT NULL ,
//		  content varchar(1000)   NOT NULL ,
//		  sender_platform_id int NOT NULL ,
//		  send_time INTEGER DEFAULT NULL ,
//		  create_time INTEGER  DEFAULT NULL,
//	     ex varchar(1024) DEFAULT NULL,
//		  PRIMARY KEY (client_msg_id)
//		)`

// 删除会话，可能会话没有
// 确认删除，告诉会话 ID
// 清空聊天记录的发，会话有，但是聊天记录没有
// DeleteMlessageFromlocalAndSvr
// db

// 不同的会话本地有一个单独的表，其中单聊的话也是这样，有一个单聊的表

// 删除的话，先删除表，在删除本地的 seq ，最后清楚这个表。
// 删除所有的消息的话，全部都是服务器来做，调用接口，然后客户端收到回调，然后删除本地的所有的信息。
// 删除一条信息，删除最新的话，会话上有一条最新的消息，删除这条消息，会话上就没有消息了，此时显示的是第二条。
// 和微信一样，我们 Go get error 分支，然后调用最新的 APi

type LocalChatLog struct {
	ClientMsgID          string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID          string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID               string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID               string `gorm:"column:recv_id;index:index_recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID     int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname       string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL        string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType          int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom              int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType          int32  `gorm:"column:content_type;index:content_type_alone" json:"contentType"`
	Content              string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead               bool   `gorm:"column:is_read" json:"isRead"`
	Status               int32  `gorm:"column:status" json:"status"`
	Seq                  int64  `gorm:"column:seq;index:index_seq;default:0" json:"seq"`
	SendTime             int64  `gorm:"column:send_time;index:index_send_time;" json:"sendTime"`
	CreateTime           int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo         string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                   string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	LocalEx              string `gorm:"column:local_ex;type:varchar(1024)" json:"localEx"`
	IsReact              bool   `gorm:"column:is_react" json:"isReact"`
	IsExternalExtensions bool   `gorm:"column:is_external_extensions" json:"isExternalExtensions"`
	MsgFirstModifyTime   int64  `gorm:"column:msg_first_modify_time" json:"msgFirstModifyTime"`
}

type LocalErrChatLog struct {
	Seq              int64  `gorm:"column:seq;primary_key" json:"seq"`
	ClientMsgID      string `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	SendTime         int64  `gorm:"column:send_time" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}
type TempCacheLocalChatLog struct {
	ClientMsgID      string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	Seq              int64  `gorm:"column:seq;default:0" json:"seq"`
	SendTime         int64  `gorm:"column:send_time;" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

// `create table if not exists  conversation (
//
//	conversation_id char(128) NOT NULL,
//	conversation_type int(11) NOT NULL,
//	user_id varchar(128)  DEFAULT NULL,
//	group_id varchar(128)  DEFAULT NULL,
//	show_name varchar(128)  NOT NULL,
//	face_url varchar(128)  NOT NULL,
//	recv_msg_opt int(11) NOT NULL ,
//	unread_count int(11) NOT NULL ,
//	latest_msg varchar(255)  NOT NULL ,
//	latest_msg_send_time INTEGER(255)  NOT NULL ,
//	draft_text varchar(255)  DEFAULT NULL ,
//	draft_timestamp INTEGER(255)  DEFAULT NULL ,
//	is_pinned int(10) NOT NULL ,
//	PRIMARY KEY (conversation_id)
//
// )
type LocalConversation struct {
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	ShowName              string `gorm:"column:show_name;type:varchar(255)" json:"showName"`
	FaceURL               string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	LatestMsg             string `gorm:"column:latest_msg;type:varchar(1000)" json:"latestMsg"`
	LatestMsgSendTime     int64  `gorm:"column:latest_msg_send_time;index:index_latest_msg_send_time" json:"latestMsgSendTime"`
	DraftText             string `gorm:"column:draft_text" json:"draftText"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	MaxSeq                int64  `gorm:"column:max_seq" json:"maxSeq"`
	MinSeq                int64  `gorm:"column:min_seq" json:"minSeq"`
	HasReadSeq            int64  `gorm:"column:has_read_seq" json:"hasReadSeq"`
	MsgDestructTime       int64  `gorm:"column:msg_destruct_time;default:604800" json:"msgDestructTime"`
	IsMsgDestruct         bool   `gorm:"column:is_msg_destruct;default:false" json:"isMsgDestruct"`
	BackgroundURL         string `json:"backgroundURL" gorm:"column:background_url;type:varchar(255)"`
}
type LocalConversationUnreadMessage struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ClientMsgID    string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	SendTime       int64  `gorm:"column:send_time" json:"sendTime"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

// message GroupRequest{
// string UserID = 1;
// string GroupID = 2;
// string HandleResult = 3;
// string ReqMsg = 4;
// string  HandleMsg = 5;
// int64 ReqTime = 6;
// string HandleUserID = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }open_im_sdk.GroupRequest == imdb.GroupRequest
type LocalAdminGroupRequest struct {
	LocalGroupRequest
}

type LocalChatLogReactionExtensions struct {
	ClientMsgID             string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	LocalReactionExtensions []byte `gorm:"column:local_reaction_extensions" json:"localReactionExtensions"`
}
type LocalWorkMomentsNotification struct {
	JsonDetail string `gorm:"column:json_detail"`
	CreateTime int64  `gorm:"create_time"`
}

type WorkMomentNotificationMsg struct {
	NotificationMsgType int32  `json:"notificationMsgType"`
	ReplyUserName       string `json:"replyUserName"`
	ReplyUserID         string `json:"replyUserID"`
	Content             string `json:"content"`
	ContentID           string `json:"contentID"`
	WorkMomentID        string `json:"workMomentID"`
	UserID              string `json:"userID"`
	UserName            string `json:"userName"`
	FaceURL             string `json:"faceURL"`
	WorkMomentContent   string `json:"workMomentContent"`
	CreateTime          int32  `json:"createTime"`
}

func (LocalWorkMomentsNotification) TableName() string {
	return "local_work_moments_notification"
}

type LocalWorkMomentsNotificationUnreadCount struct {
	UnreadCount int `gorm:"unread_count" json:"unreadCount"`
}

func (LocalWorkMomentsNotificationUnreadCount) TableName() string {
	return "local_work_moments_notification_unread_count"
}

type NotificationSeqs struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	Seq            int64  `gorm:"column:seq" json:"seq"`
}

func (NotificationSeqs) TableName() string {
	return "local_notification_seqs"
}

type LocalUpload struct {
	PartHash   string `gorm:"column:part_hash;primary_key" json:"partHash"`
	UploadID   string `gorm:"column:upload_id;type:varchar(1000)" json:"uploadID"`
	UploadInfo string `gorm:"column:info;type:varchar(2000)" json:"info"`
	ExpireTime int64  `gorm:"column:expire_time" json:"expireTime"`
	CreateTime int64  `gorm:"column:create_time" json:"createTime"`
}

func (LocalUpload) TableName() string {
	return "local_uploads"
}

type LocalStranger struct {
	UserID           string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname         string `gorm:"column:name;type:varchar(255)" json:"nickname"`
	FaceURL          string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AppMangerLevel   int32  `gorm:"column:app_manger_level" json:"-"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	GlobalRecvMsgOpt int32  `gorm:"column:global_recv_msg_opt" json:"globalRecvMsgOpt"`
}

func (LocalStranger) TableName() string {
	return "local_stranger"
}
