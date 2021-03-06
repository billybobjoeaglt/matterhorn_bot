package chat

import (
	"cloud.google.com/go/datastore"
	"context"
	"math/rand"
	"time"
)

var ChatRand *rand.Rand

func init() {
	ChatRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const ChatKeyKind = "Chat"
const KeyNamespace = "matterhorn-bot"

type Chat struct {
	Name     string
	Type     string
	Settings ChatSettings `datastore:",noindex"`
}

type ChatSettings struct {
	NSFW       bool        `datastore:",noindex"`
	AlertTimes []AlertTime `datastore:",noindex"`
	KeyWords   []KeyWord   `datastore:",noindex"`
	QuotesDoc  int         `datastore:",noindex"`
}

type AlertTime struct {
	ID      int    `datastore:",noindex" json:"id"`
	Time    string `datastore:",noindex" json:"time"`
	Message string `datastore:",noindex" json:"msg"`
}

type KeyWord struct {
	ID      int    `datastore:",noindex" json:"id"`
	Key     string `datastore:",noindex" json:"key"`
	Message string `datastore:",noindex" json:"msg"`
}

func MakeAlertTime(time, message string) AlertTime {
	return AlertTime{int(ChatRand.Int31()), time, message}
}

func MakeKeyWord(key, message string) KeyWord {
	return KeyWord{int(ChatRand.Int31()), key, message}
}

func NewKeyFromChatID(chatID int64) *datastore.Key {
	key := datastore.IDKey(ChatKeyKind, chatID, nil)
	key.Namespace = KeyNamespace
	return key
}

type Datastore struct {
	Client *datastore.Client
	Ctx    context.Context
}

func NewDatastore(client *datastore.Client, ctx context.Context) *Datastore {
	return &Datastore{
		Client: client,
		Ctx:    ctx,
	}
}

func (ds *Datastore) GetChat(chatID int64) (chat Chat, exists bool, err error) {
	err = ds.Client.Get(ds.Ctx, NewKeyFromChatID(chatID), &chat)

	exists = err != datastore.ErrNoSuchEntity

	return
}

func (ds *Datastore) InsertChat(chat Chat, chatID int64) error {
	_, err := ds.Client.Put(ds.Ctx, NewKeyFromChatID(chatID), &chat)

	return err
}

func (ds *Datastore) UpdateChat(makeChangesFunc func(Chat) Chat, chatID int64) (Chat, error) {
	tx, err := ds.Client.NewTransaction(ds.Ctx)
	if err != nil {
		return Chat{}, err
	}

	chatKey := NewKeyFromChatID(chatID)

	var oldChat Chat
	if err := tx.Get(chatKey, &oldChat); err != nil {
		return Chat{}, err
	}

	updatedChat := makeChangesFunc(oldChat)

	if _, err := tx.Put(chatKey, &updatedChat); err != nil {
		return Chat{}, err
	}
	if _, err := tx.Commit(); err != nil {
		return Chat{}, err
	}

	return updatedChat, nil
}
