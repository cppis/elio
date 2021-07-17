package elio

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config config
type Config struct {
	viper *viper.Viper
}

// String object to string
func (c *Config) String() string {
	return fmt.Sprintf("Config::%p", c)
}

// NewConfig new config
func NewConfig() (c *Config) {
	c = new(Config)
	if nil != c {
		c.viper = viper.New()
	}

	return c
}

/*// viper config precedence order
Viper uses the following precedence order. Each item takes precedence over the item below it:
  * explicit call to Set
  * flag
  * env
  * config
  * key/value store
  * default
//*/

// Load load
func (c *Config) Load(path string) (err error) {
	c.viper.AddConfigPath("./")
	c.viper.AutomaticEnv()

	AppTrace().Str(LogObject, c.String()).Msgf("begin to load config:%s", path)
	if "" != path {
		c.viper.SetConfigFile(path)
		c.viper.SetConfigType("yaml")
		c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // .env notOK, env Ok
		err = c.viper.ReadInConfig()
		if nil != err {
			return err
		}
	}

	// bind env
	err = c.viper.BindEnv("elio.log.level", "ELIO_LOG_LEVEL")
	err = c.viper.BindEnv("elio.log.out", "ELIO_LOG_OUT")
	err = c.viper.BindEnv("elio.log.json", "ELIO_LOG_JSON")
	err = c.viper.BindEnv("elio.log.color", "ELIO_LOG_COLOR")
	err = c.viper.BindEnv("elio.log.shortCaller", "ELIO_LOG_SHORTCALLER")
	err = c.viper.BindEnv("elio.app.intervalMs", "ELIO_APP_INTERVALMS")
	err = c.viper.BindEnv("elio.app.fetchLimit", "ELIO_APP_FETCHLIMIT")

	// viper expects SYSENV_* as system environment variables
	c.viper.SetEnvPrefix(viper.GetString("sysenv"))
	c.viper.SetConfigType("env")
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // .env notOK, env Ok
	if nil == c.viper.MergeInConfig() {
	}

	////////////////////////////////////////////////////////////////
	// // .env override 에 이슈가 없는 코드 블럭
	//AppTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")
	c.viper.SetConfigFile(".env")
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))		// .env OK, env notOk
	if nil == c.viper.MergeInConfig() {
	}
	////////////////////////////////////////////////////////////////

	return err
}

// Exists check key exists
func (c *Config) Exists(path string) bool {
	return c.viper.IsSet(path)
}

// Get get
func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

// GetOrDefault get or default
func (c *Config) GetOrDefault(key string, d interface{}) (interface{}, bool) {
	if c.Exists(key) {
		return c.viper.Get(key), true
	}
	return d, false
}

// Set set interface value of key
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

// GetString get string value of key
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetStringOrDefault get string or default value of key
func (c *Config) GetStringOrDefault(key string, d string) (string, bool) {
	if c.Exists(key) {
		return c.viper.GetString(key), true
	}
	return d, false
}

// // GetString get string value of key
// func (c *Config) GetString(key string) string {
// 	return c.viper.GetString(key)
// }

// GetBool get bool value of key
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetBoolOrDefault get bool or default value or of key
func (c *Config) GetBoolOrDefault(key string, d bool) (bool, bool) {
	if c.Exists(key) {
		return c.viper.GetBool(key), true
	}
	return d, false
}

// SetBool set bool value of key
func (c *Config) SetBool(key string, value bool) {
	c.viper.Set(key, value)
}

// GetInt get int value of key
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetIntOrDefault get int or default value of key
func (c *Config) GetIntOrDefault(key string, d int) (int, bool) {
	if c.Exists(key) {
		return c.viper.GetInt(key), true
	}
	return d, false
}

// GetUint get uint value of key
func (c *Config) GetUint(key string) uint {
	return c.viper.GetUint(key)
}

// GetUintOrDefault get uint or default value of key
func (c *Config) GetUintOrDefault(key string, d uint) (uint, bool) {
	if c.Exists(key) {
		return c.viper.GetUint(key), true
	}
	return d, false
}

// GetInt32 get int32 value of key
func (c *Config) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetInt32OrDefault get int32 or default value of key
func (c *Config) GetInt32OrDefault(key string, d int32) (int32, bool) {
	if c.Exists(key) {
		return c.viper.GetInt32(key), true
	}
	return d, false
}

// GetUint32 get Uint32 value of key
func (c *Config) GetUint32(key string) uint32 {
	return c.viper.GetUint32(key)
}

// GetUint32OrDefault get uint32 or default value of key
func (c *Config) GetUint32OrDefault(key string, d uint32) (uint32, bool) {
	if c.Exists(key) {
		return c.viper.GetUint32(key), true
	}
	return d, false
}

// GetInt64 get int64 value of key
func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetInt64OrDefault get int64 or default value of key
func (c *Config) GetInt64OrDefault(key string, d int64) (int64, bool) {
	if c.Exists(key) {
		return c.viper.GetInt64(key), true
	}
	return d, false
}

// GetUint64 get uint64 value of key
func (c *Config) GetUint64(key string) uint64 {
	return c.viper.GetUint64(key)
}

// GetUint64OrDefault get uint64 or default value of key
func (c *Config) GetUint64OrDefault(key string, d uint64) (uint64, bool) {
	if c.Exists(key) {
		return c.viper.GetUint64(key), true
	}
	return d, false
}

// GetFloat64 get float64 value of key
func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// GetFloat64OrDefault get float64 or default value of key
func (c *Config) GetFloat64OrDefault(key string, d float64) (float64, bool) {
	if c.Exists(key) {
		return c.viper.GetFloat64(key), true
	}
	return d, false
}
