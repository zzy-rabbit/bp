package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/tool/mail/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"gopkg.in/gomail.v2"
)

func (s *service) Send(ctx context.Context, message api.Message) xerror.IError {
	m := gomail.NewMessage()

	m.SetHeader("From", message.From)
	m.SetHeader("To", message.To)
	m.SetHeader("Subject", message.Subject)

	m.SetBody("text/html", message.Body)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	err := d.DialAndSend(m)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "send mail fail %v", err)
		return xerror.Extend(xerror.ErrInternalError, "send mail fail %v", err)
	}
	return nil
}
