package config

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
)

type Config interface {
	GetServerConn() string
	GetLogLevel() string
	GetWhiteListDomains() []string
	GetVirtualFSPath() string
	GetOwnerKey() string
	GetEnableHTTPS() bool
	GetMaxUploadSize() int
	GetDiscordLink() string
}

type config struct {
	ServerHost string `env:"SERVER_HOST, required"`
	ServerPort int    `env:"SERVER_PORT, required"`

	LogLevel         string   `env:"LOG_LEVEL, required"`
	WhiteListDomains []string `env:"WHITE_LIST_DOMAINS, required"`
	VirtualFSPath    string   `env:"VIRTUAL_FS_PATH, required"`
	OwnerKey         string   `env:"OWNER_KEY, required"`
	EnableHTTPS      bool     `env:"ENABLE_HTTPS, required"`
	MaxUploadSize    int      `env:"MAX_UPLOAD_SIZE, required"`

	DiscordLink string `env:"DISCORD_LINK, required"`
}

func NewFromENV(ctx context.Context) (Config, error) {
	c := &config{}
	if err := envconfig.Process(ctx, c); err != nil {
		return c, err
	}
	return c, nil
}

func (c *config) GetServerConn() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}

func (c *config) GetLogLevel() string {
	return c.LogLevel
}

func (c *config) GetWhiteListDomains() []string {
	return c.WhiteListDomains
}

func (c *config) GetVirtualFSPath() string {
	return c.VirtualFSPath
}

func (c *config) GetOwnerKey() string {
	return c.OwnerKey
}

func (c *config) GetEnableHTTPS() bool {
	return c.EnableHTTPS
}

func (c *config) GetMaxUploadSize() int {
	return c.MaxUploadSize
}

func (c *config) GetDiscordLink() string {
	return c.DiscordLink
}
