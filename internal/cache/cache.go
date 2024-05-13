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

package cache

import (
	"context"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdkerrs"
	"sync"
)

type BaseInfo struct {
	Nickname      string
	FaceURL       string
	BackgroundURL string
}
type Cache struct {
	user            *user.User
	friend          *friend.Friend
	userMap         sync.Map
	conversationMap sync.Map
	loginUserID     string
	ch              chan common.Cmd2Value
}

func NewCache(user *user.User, friend *friend.Friend, loginUserID string, ch chan common.Cmd2Value) *Cache {
	return &Cache{user: user, friend: friend, loginUserID: loginUserID, ch: ch}
}

func (c *Cache) Update(userID, faceURL, nickname, backgroundURL string) {
	c.userMap.Store(userID, BaseInfo{FaceURL: faceURL, Nickname: nickname, BackgroundURL: backgroundURL})
}
func (c *Cache) UpdateConversation(conversation model_struct.LocalConversation) {
	c.conversationMap.Store(conversation.ConversationID, conversation)
}
func (c *Cache) UpdateConversations(conversations []*model_struct.LocalConversation) {
	for _, conversation := range conversations {
		c.conversationMap.Store(conversation.ConversationID, *conversation)
	}
}
func (c *Cache) GetAllConversations() (conversations []*model_struct.LocalConversation) {
	c.conversationMap.Range(func(key, value interface{}) bool {
		temp := value.(model_struct.LocalConversation)
		conversations = append(conversations, &temp)
		return true
	})
	return conversations
}
func (c *Cache) GetAllHasUnreadMessageConversations() (conversations []*model_struct.LocalConversation) {
	c.conversationMap.Range(func(key, value interface{}) bool {
		temp := value.(model_struct.LocalConversation)
		if temp.UnreadCount > 0 {
			conversations = append(conversations, &temp)
		}
		return true
	})
	return conversations
}

func (c *Cache) GetConversation(conversationID string) model_struct.LocalConversation {
	var result model_struct.LocalConversation
	conversation, ok := c.conversationMap.Load(conversationID)
	if ok {
		result = conversation.(model_struct.LocalConversation)
	}
	return result
}

func (c *Cache) GetUserNameFaceURLAndBackgroundUrl(ctx context.Context, userID string) (faceURL, name, backgroundURL string, err error) {
	//find in cache
	if value, ok := c.userMap.Load(userID); ok {
		info := value.(*BaseInfo)
		return info.FaceURL, info.Nickname, info.BackgroundURL, nil
	}
	// 从本地获取
	friendInfo, err := c.friend.Db().GetFriendInfoByFriendUserID(ctx, userID)
	if err == nil {
		faceURL = friendInfo.FaceURL
		if friendInfo.Remark != "" {
			name = friendInfo.Remark
		} else {
			name = friendInfo.Nickname
		}
		backgroundURL = friendInfo.BackgroundURL
		return faceURL, name, backgroundURL, nil
	}
	//从服务端上传
	svrFriends, err := c.friend.GetFriendByIdsSvr(ctx, []string{userID})
	if err == nil && len(svrFriends) > 0 {
		svrFriend := svrFriends[0]
		faceURL = svrFriend.FaceURL
		if svrFriend.Remark != "" {
			name = svrFriend.Remark
		} else {
			name = svrFriend.Nickname
		}
		backgroundURL = svrFriend.BackgroundUrl
		return faceURL, name, backgroundURL, nil
	}
	//get from server db
	users, err := c.user.FindFullProfile(ctx, userID)
	if err != nil {
		return "", "", "", err
	}
	if len(users) == 0 {
		return "", "", "", sdkerrs.ErrUserIDNotFound.Wrap(userID)
	}
	c.userMap.Store(userID, &BaseInfo{FaceURL: users[0].FaceURL, Nickname: users[0].Nickname, BackgroundURL: ""})
	return users[0].FaceURL, users[0].Nickname, "", nil
}

func (c *Cache) Store(userID string, data *BaseInfo) {
	c.userMap.Store(userID, data)
}
