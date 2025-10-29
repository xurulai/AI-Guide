package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Dify       DifyConfig       `mapstructure:"dify"`
	Redis      RedisConfig      `mapstructure:"redis"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Log        LogConfig        `mapstructure:"log"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	Business   BusinessConfig   `mapstructure:"business"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DifyConfig Dify配置
type DifyConfig struct {
	BaseURL   string              `mapstructure:"base_url"`
	APIKey    string              `mapstructure:"api_key"`
	Timeout   time.Duration       `mapstructure:"timeout"`
	Workflows DifyWorkflowsConfig `mapstructure:"workflows"`
}

// DifyWorkflowsConfig Dify工作流配置
type DifyWorkflowsConfig struct {
	Planner   DifyWorkflowConfig            `mapstructure:"planner"`
	Executors map[string]DifyWorkflowConfig `mapstructure:"executors"`
}

// DifyWorkflowConfig 单个工作流配置
type DifyWorkflowConfig struct {
	AppID   string        `mapstructure:"app_id"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr            string        `mapstructure:"addr"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxRetries      int           `mapstructure:"max_retries"`
	SessionTTL      time.Duration `mapstructure:"session_ttl"`
	CacheTTL        time.Duration `mapstructure:"cache_ttl"`
	RateLimitWindow time.Duration `mapstructure:"rate_limit_window"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Auth      AuthConfig      `mapstructure:"auth"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	Burst             int  `mapstructure:"burst"`
}

// AuthConfig 鉴权配置
type AuthConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	JWTSecret   string        `mapstructure:"jwt_secret"`
	TokenExpire time.Duration `mapstructure:"token_expire"`
}

// BusinessConfig 业务配置
type BusinessConfig struct {
	Session SessionConfig `mapstructure:"session"`
	Product ProductConfig `mapstructure:"product"`
	Retry   RetryConfig   `mapstructure:"retry"`
}

// SessionConfig 会话配置
type SessionConfig struct {
	MaxMessages  int    `mapstructure:"max_messages"`
	DefaultStyle string `mapstructure:"default_style"`
}

// ProductConfig 商品配置
type ProductConfig struct {
	TopK        int  `mapstructure:"top_k"`
	EnableCache bool `mapstructure:"enable_cache"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts  int           `mapstructure:"max_attempts"`
	InitialDelay time.Duration `mapstructure:"initial_delay"`
	MaxDelay     time.Duration `mapstructure:"max_delay"`
	Multiplier   int           `mapstructure:"multiplier"`
}

var globalConfig *Config

// Load 加载配置
func Load(configPath string, env string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// 读取默认配置
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 读取环境特定配置
	if env != "" {
		v.SetConfigName(fmt.Sprintf("config.%s", env))
		if err := v.MergeInConfig(); err != nil {
			// 环境配置可选，不存在时不报错
			fmt.Printf("No environment config found for %s, using default\n", env)
		}
	}

	// 支持环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量读取敏感信息
	if apiKey := v.GetString("DIFY_API_KEY"); apiKey != "" {
		cfg.Dify.APIKey = apiKey
	}
	if jwtSecret := v.GetString("JWT_SECRET"); jwtSecret != "" {
		cfg.Middleware.Auth.JWTSecret = jwtSecret
	}
	if mysqlPass := v.GetString("MYSQL_PASSWORD"); mysqlPass != "" {
		cfg.MySQL.Password = mysqlPass
	}
	if redisPass := v.GetString("REDIS_PASSWORD"); redisPass != "" {
		cfg.Redis.Password = redisPass
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// GetDSN 获取MySQL DSN
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}
