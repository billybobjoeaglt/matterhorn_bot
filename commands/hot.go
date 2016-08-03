package commands

import (
	"net/http"

	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"

	"bytes"
	"io"

	"gopkg.in/telegram-bot-api.v4"
)

type HotHandler struct {
}

var hotHandlerInfo = CommandInfo{
	Command:     "hot",
	Args:        "",
	Permission:  3,
	Description: "gets the hotness score from howhot.io",
	LongDesc:    "",
	Usage:       "/hot",
	Examples: []string{
		"/hot",
	},
	ResType: "message",
}

func (h HotHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig
	var errMsg tgbotapi.MessageConfig

	defer func(bot *tgbotapi.BotAPI) {
		if errMsg.Text != "" {
			bot.Send(errMsg)
		}
	}(bot)

	replyToMsg := message.MessageID

	if message.Photo == nil {
		if message.ReplyToMessage != nil && message.ReplyToMessage.Photo != nil {
			message.Photo = message.ReplyToMessage.Photo
			replyToMsg = message.ReplyToMessage.MessageID
		} else {
			return
		}
	}

	fileurl, err := bot.GetFileDirectURL((*(message.Photo))[len((*(message.Photo)))-1].FileID)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	photoResp, err := http.Get(fileurl)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("browseFile", fileurl)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}
	if _, err = io.Copy(fw, photoResp.Body); err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, "https://howhot.io/main.php", &b)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	repBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	hot := Hotness{}
	err = json.Unmarshal(repBytes, &hot)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	if !hot.Success {
		errMsg = NewErrorMessage(message.Chat.ID, errors.New("Failed because "+hot.Action))
		return
	}

	htMsg := unkMessageToHotMessage(hot.Message)

	msg = tgbotapi.NewMessage(message.Chat.ID, "<b>Hotness:</b>\nSex: "+htMsg.Gender+"\nAge: "+htMsg.Age+"\nHotness: "+htMsg.Hotness)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = replyToMsg

	_, err = bot.Send(msg)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
	}
}

func (h HotHandler) Info() *CommandInfo {
	return &hotHandlerInfo
}

func (h HotHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return true, ""
}

type Hotness struct {
	Success  bool        `json:"success"`
	Category string      `json:"category"`
	Action   string      `json:"action"`
	Label    string      `json:"label"`
	Message  interface{} `json:"message"`
}

type HotMessage struct {
	Vertices []struct {
		X string `json:"X"`
		Y string `json:"Y"`
	} `json:"vertices"`
	Gender    string `json:"gender"`
	Hotness   string `json:"hotness"`
	Age       string `json:"age"`
	ImageData string `json:"image_data"`
}

func unkMessageToHotMessage(msg interface{}) HotMessage {
	msgMap := msg.(map[string]interface{})

	return HotMessage{
		Gender:    msgMap["gender"].(string),
		Hotness:   msgMap["hotness"].(string),
		Age:       msgMap["age"].(string),
		ImageData: msgMap["image_data"].(string),
	}
}