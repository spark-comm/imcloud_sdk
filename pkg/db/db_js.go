//go:build js && wasm
// +build js,wasm

package db

import (
	"context"
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/wasm/exec"
	"open_im_sdk/wasm/indexdb"
)

var ErrType = errors.New("from javascript data type err")

type IndexDB struct {
	*indexdb.LocalUsers
	*indexdb.LocalConversations
	*indexdb.LocalChatLogs
	*indexdb.LocalSuperGroupChatLogs
	*indexdb.LocalSuperGroup
	*indexdb.LocalConversationUnreadMessages
	*indexdb.LocalGroups
	*indexdb.LocalGroupMember
	*indexdb.LocalCacheMessage
	*indexdb.FriendRequest
	*indexdb.Black
	*indexdb.Friend
	*indexdb.LocalGroupRequest
	*indexdb.LocalChatLogReactionExtensions
	*indexdb.NotificationSeqs
	*indexdb.LocalMoments
	loginUserID string
}

func (i IndexDB) GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) DeleteUpload(ctx context.Context, partHash string) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) DeleteExpireUpload(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) Close(ctx context.Context) error {
	return nil
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := exec.Exec(userID, dataDir)
	return err
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string) (*IndexDB, error) {
	i := &IndexDB{
		LocalUsers:                      indexdb.NewLocalUsers(),
		LocalConversations:              indexdb.NewLocalConversations(),
		LocalChatLogs:                   indexdb.NewLocalChatLogs(loginUserID),
		LocalSuperGroupChatLogs:         indexdb.NewLocalSuperGroupChatLogs(),
		LocalSuperGroup:                 indexdb.NewLocalSuperGroup(),
		LocalConversationUnreadMessages: indexdb.NewLocalConversationUnreadMessages(),
		LocalGroups:                     indexdb.NewLocalGroups(),
		LocalGroupMember:                indexdb.NewLocalGroupMember(),
		LocalCacheMessage:               indexdb.NewLocalCacheMessage(),
		FriendRequest:                   indexdb.NewFriendRequest(loginUserID),
		Black:                           indexdb.NewBlack(loginUserID),
		Friend:                          indexdb.NewFriend(loginUserID),
		LocalGroupRequest:               indexdb.NewLocalGroupRequest(),
		LocalChatLogReactionExtensions:  indexdb.NewLocalChatLogReactionExtensions(),
		NotificationSeqs:                indexdb.NewNotificationSeqs(),
		LocalMoments:                    indexdb.NewLocalMoments(),
		loginUserID:                     loginUserID,
	}
	err := i.InitDB(ctx, loginUserID, dbDir)
	if err != nil {
		return nil, err
	}
	return i, nil
}
