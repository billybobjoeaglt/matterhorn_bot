package commands

import (
	"errors"
	"os/exec"

	"gopkg.in/telegram-bot-api.v4"
)

type FortuneHandler struct {
}

var fortuneHandlerInfo = CommandInfo{
	Command:     "fortune",
	Args:        "",
	Permission:  3,
	Description: "reads a unix fortune",
	LongDesc:    "",
	Usage:       "/fortune",
	Examples: []string{
		"/fortune",
	},
	ResType: "message",
}

func (responder FortuneHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, fortune := GetFortune()
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, fortune)
	}
	bot.Send(msg)
}

func (responder FortuneHandler) Info() *CommandInfo {
	return &fortuneHandlerInfo
}

func GetFortune() (error, string) {
	cfOut, err := exec.Command("command", "-v", "fortune").Output()
	if err != nil {
		return err, ""
	}
	if len(cfOut) == 0 {
		return errors.New("Fortune command is missing. Check README.md for install instructions."), ""
	}

	fOut, err := exec.Command("fortune", "-a", "fortunes", "riddles").Output()
	if err != nil {
		return err, ""
	}
	return nil, string(fOut)
}