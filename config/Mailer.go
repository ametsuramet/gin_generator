package config

type MailerConfiguration struct {
	Server     string
	Port       uint
	Username   string
	Password   string
	UseTLS     bool
	Sender     string
	MaxAttempt uint
}
