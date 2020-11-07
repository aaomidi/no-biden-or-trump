package telegram

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

type Telegram struct {
	token string
	bot   *tb.Bot
	log   *log.Entry
}

func New(token string) Telegram {
	return Telegram{
		token: token,
		log:   log.WithField("source", "telegram"),
	}
}

func (t *Telegram) Create() error {
	bot, err := tb.NewBot(tb.Settings{
		Token:  t.token,
		Poller: &tb.LongPoller{Timeout: 15 * time.Second},
	})
	t.bot = bot

	if err != nil {
		return errors.Wrap(err, "creation failed")
	}

	return nil
}

func (t *Telegram) Start() {
	t.bot.Handle(tb.OnText, t.onText)

	t.bot.Start()
}

func (t *Telegram) onText(m *tb.Message) {
	msg := strings.ToLower(m.Text)

	if !strings.Contains(msg, "trump") && !strings.Contains(msg, "biden") {
		return
	}
	user, err := t.bot.ChatMemberOf(m.Chat, m.Sender)
	if err != nil {
		t.log.WithError(err).Infoln("call to chatMemberOf failed")
		return
	}

	isAdmin := user.Role == tb.Creator || user.Role == tb.Administrator

	if isAdmin {
		return
	}

	err = t.bot.Delete(m)

	if err != nil {
		t.log.WithError(err).Infoln("call to delete failed")
		return
	}
}

func (t *Telegram) Stop() {
	t.bot.Stop()
}
