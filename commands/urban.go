package commands

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"gopkg.in/telegram-bot-api.v4"
)

type UrbanHandler struct {
}

var urbanHandlerInfo = CommandInfo{
	Command:     "urban",
	Args:        `(.+)`,
	Permission:  3,
	Description: "gets urban dictionary of word",
	LongDesc:    "",
	Usage:       "/urban [word]",
	Examples: []string{
		"/urban shrek",
	},
	ResType: "message",
}

func (h *UrbanHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, def := GetUrban(args[0])
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "<b>"+def.Word+"</b>\n———\n"+def.Definition+"\n\n<i>"+def.Example+"</i>")
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = "HTML"
	}
	bot.Send(msg)
}

func (h *UrbanHandler) Info() *CommandInfo {
	return &urbanHandlerInfo
}

func (h *UrbanHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *UrbanHandler) Setup(setupFields map[string]interface{}) {

}

type UrbanDefinition struct {
	Word       string
	Definition string
	Example    string
	Null       bool
}

func GetUrban(term string) (error, UrbanDefinition) {
	resp, err := http.Get("http://api.urbandictionary.com/v0/define?term=" + url.QueryEscape(term))
	if err != nil {
		return err, UrbanDefinition{Null: true}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid Status Code: " + resp.Status), UrbanDefinition{Null: true}
	}

	var jsonRes map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonRes)
	if err != nil {
		return err, UrbanDefinition{Null: true}
	}

	list := jsonRes["list"].([]interface{})
	if len(list) <= 0 {
		return errors.New("No Definition Found"), UrbanDefinition{Null: true}
	}
	def := list[0].(map[string]interface{})

	return nil, UrbanDefinition{
		Word:       def["word"].(string),
		Definition: def["definition"].(string),
		Example:    def["example"].(string),
		Null:       false,
	}

}
