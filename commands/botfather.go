package commands

import "gopkg.in/telegram-bot-api.v4"

type BotFatherHandler struct {
}

var botFatherHandlerInfo = CommandInfo{
	Command:     "botfather",
	Args:        "",
	Permission:  3,
	Description: "gets botfather list",
	LongDesc:    "",
	Usage:       "/botfather]",
	Examples: []string{
		"/botfather",
	},
	ResType: "message",
}

func (responder BotFatherHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msgStr string
	for _, cmd := range *CommandList {
		msgStr += cmd.Info().Command + " - " + cmd.Info().Description + "\n"
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgStr)
	bot.Send(msg)
}

func (responder BotFatherHandler) Info() *CommandInfo {
	return &botFatherHandlerInfo
}