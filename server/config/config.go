package config

// Config refers
type Config struct {
	Owner       string
	Repo        string
	HTTPListen  string
	AccessToken string
}

// NewConfig creates a
func NewConfig() Config {
	return Config{
		Owner:      "",
		Repo:       "",
		HTTPListen: "",
	}
}
