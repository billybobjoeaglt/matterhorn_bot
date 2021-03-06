package commands

import (
	"net/url"

	"gopkg.in/telegram-bot-api.v4"
)

type LmgtfyHandler struct {
}

var lmgtfyHandlerInfo = CommandInfo{
	Command:     "lmgtfy",
	Args:        "(.+)",
	Permission:  3,
	Description: "let me google that for you",
	LongDesc:    "",
	Usage:       "/lmgtfy [input]",
	Examples: []string{
		"/lmgtfy hello world",
	},
	ResType: "message",
}

func (h *LmgtfyHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "http://lmgtfy.com/?q="+url.QueryEscape(args[0]))
	bot.Send(msg)
}

func (h *LmgtfyHandler) Info() *CommandInfo {
	return &lmgtfyHandlerInfo
}

func (h *LmgtfyHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *LmgtfyHandler) Setup(setupFields map[string]interface{}) {

}
