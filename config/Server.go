package config

type ServerConfiguration struct {
	Mode            string
	Addr            string
	Environment     string
	LogDuration     int
	ShutdownTimeout int
	BaseURL         string
	ClientURL       string
}
