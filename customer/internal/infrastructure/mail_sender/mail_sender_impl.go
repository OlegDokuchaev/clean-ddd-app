package mail_sender

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	customerAplication "customer/internal/application/customer"

	"github.com/wneessen/go-mail"
)

type MailSenderImpl struct {
	cfg    *Config
	client *mail.Client
}

func New(cfg *Config, client *mail.Client) *MailSenderImpl {
	return &MailSenderImpl{
		cfg:    cfg,
		client: client,
	}
}

func (m *MailSenderImpl) SendOtp(ctx context.Context, toEmail string, code string) error {
	subj := "Your OTP Code"
	body := fmt.Sprintf("Your one-time code is: %s\nIt will expire shortly.", code)
	return m.sendMail(ctx, toEmail, subj, body)
}

// SendPasswordResetLink отправляет письмо со ссылкой сброса
func (m *MailSenderImpl) SendPasswordResetLink(ctx context.Context, toEmail string, token string) error {
	subj := "Password Reset Request"
	resetURL := fmt.Sprintf("token=%s", token)
	bodyTemplate := `Dear user,

We received a request to reset your password. Please click the link below to reset it:

{{ .ResetURL }}

If you did not request this, please ignore this email.
`
	t, err := template.New("reset").Parse(bodyTemplate)
	if err != nil {
		return ErrMailSendFailed
	}
	var sb strings.Builder
	err = t.Execute(&sb, map[string]any{"ResetURL": resetURL})
	if err != nil {
		return ErrMailSendFailed
	}
	return m.sendMail(ctx, toEmail, subj, sb.String())
}

func (m *MailSenderImpl) sendMail(ctx context.Context, to string, subj string, body string) error {
	msg := mail.NewMsg()
	if err := msg.From(fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.FromAddr)); err != nil {
		return ErrMailSendFailed
	}
	if err := msg.To(to); err != nil {
		return ErrMailSendFailed
	}
	msg.Subject(subj)
	msg.SetBodyString(mail.TypeTextHTML, body)

	if err := m.client.DialAndSendWithContext(ctx, msg); err != nil {
		return ErrMailSendFailed
	}
	return nil
}

var _ customerAplication.MailSender = (*MailSenderImpl)(nil)
