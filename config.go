package pg

import "github.com/jackc/pgx"

// ConnConfig is a database connection config
type ConnConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

// PoolConfig is a database connection config
type PoolConfig struct {
	ConnConfig
	MaxConnections int
}

func (c ConnConfig) config() pgx.ConnConfig {
	return pgx.ConnConfig{
		Host:     c.Host,
		Port:     c.Port,
		Database: c.Database,
		User:     c.User,
		Password: c.Password,
	}
}

func (c PoolConfig) config() pgx.ConnPoolConfig {
	return pgx.ConnPoolConfig{
		ConnConfig:     c.ConnConfig.config(),
		MaxConnections: c.MaxConnections,
	}
}
