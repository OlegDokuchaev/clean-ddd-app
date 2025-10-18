package mail_sender

import "github.com/wneessen/go-mail"

func NewClient(cfg *Config) (*mail.Client, error) {
	return mail.NewClient(
		cfg.Host,
		mail.WithPort(cfg.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
	)
}
