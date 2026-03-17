package finam

type ClientConfig struct {
	Addr string
}

func (cfg *ClientConfig) setDefaults() *ClientConfig {
	if cfg.Addr == "" {
		cfg.Addr = "api.finam.ru:443"
	}
	return cfg
}
