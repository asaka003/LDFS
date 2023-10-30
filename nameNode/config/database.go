package config

type MysqlConfig struct {
	Host         string
	User         string
	Password     string
	DB           string
	Port         string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Password string
	Port     string
	PoolSize int
	DB       int
}

type SqliteConfig struct {
	Dialect string
	DbFile  string
}
