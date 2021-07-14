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

	/*//
	//AlogTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")
	//c.viper.SetConfigType("env")
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))		// .env notOK, env Ok
	//if nil == c.viper.MergeInConfig() {
	//}

	c.viper.SetConfigType("env")
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))		// .env notOK, env Ok
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))		// .env notOK, env Ok
	if nil == c.viper.MergeInConfig() {
	}

	//AlogTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")
	c.viper.SetConfigFile(".env")
	c.viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))		// .env OK, env notOk
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if nil == c.viper.MergeInConfig() {
		// dir, err := os.Getwd()
		// fmt.Printf("working.dir:%s\n", dir)

		// var j []byte
		// if j, err = ioutil.ReadFile(".env"); err == nil {
		// 		fmt.Printf("dump .env:\n%s\n", hex.Dump(j))
		// }
	}

	//AlogTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")
	// fmt.Printf(".env config\n")
	// fmt.Printf("config.all: %v\n", c.viper.AllSettings())
	// fmt.Printf("env[heatgo.log.level]: %v\n", c.viper.Get("heatgo.log.level"))
	// fmt.Printf("env[heatgo.log.json]: %v\n", c.viper.Get("heatgo.log.json"))
	// fmt.Printf("env[heatgo.log.out]: %v\n", c.viper.Get("heatgo.log.out"))
	// fmt.Println()
	//*/

	//AlogTrace().Str(elf.LogObject, c.String()).Msgf("begin to load config:%s", path)
	if "" != path {
		c.viper.SetConfigFile(path)
		c.viper.SetConfigType("json")
		err = c.viper.ReadInConfig()
		if nil != err {
			return err
		}

		// fmt.Printf("json config\n")
		// fmt.Printf("config.all: %v\n", c.viper.AllSettings())
		// fmt.Printf("env[heatgo.log.level]: %v\n", c.viper.Get("heatgo.log.level"))
		// fmt.Printf("env[heatgo.log.json]: %v\n", c.viper.Get("heatgo.log.json"))
		// fmt.Printf("env[heatgo.log.out]: %v\n", c.viper.Get("heatgo.log.out"))
		// fmt.Printf("env[heatgo.log.color]: %v\n", c.viper.Get("heatgo.log.color"))
		// fmt.Println()
	}

	// bind env
	err = c.viper.BindEnv("heatgo.log.level", "HEATGO_LOG_LEVEL")
	err = c.viper.BindEnv("heatgo.log.out", "HEATGO_LOG_OUT")
	err = c.viper.BindEnv("heatgo.log.json", "HEATGO_LOG_JSON")
	err = c.viper.BindEnv("heatgo.log.color", "HEATGO_LOG_COLOR")
	err = c.viper.BindEnv("heatgo.log.shortCaller", "HEATGO_LOG_SHORTCALLER")

	err = c.viper.BindEnv("heatgo.app.interval", "HEATGO_APP_INTERVAL")
	err = c.viper.BindEnv("heatgo.app.floors", "HEATGO_APP_FLOORS")
	err = c.viper.BindEnv("heatgo.app.fetchLimit", "HEATGO_APP_FETCHLIMIT")

	//c.viper.SetDefault("heatgo.log.level", "debug")
	//c.viper.SetDefault("heatgo.log.out", "stdout,server.log")
	//c.viper.SetDefault("heatgo.log.json", false)
	//c.viper.SetDefault("heatgo.log.color", true)
	//c.viper.SetDefault("heatgo.log.shortCaller", false)

	// ////////////////////////////////////////////////////////////////
	// // env override 에 이슈가 없는 코드 블럭
	// //AlogTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")

	// viper expects SYSENV_* as system environment variables 
	c.viper.SetEnvPrefix(viper.GetString("sysenv"))
	c.viper.SetConfigType("env")
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // .env notOK, env Ok
	if nil == c.viper.MergeInConfig() {
	}
	// ////////////////////////////////////////////////////////////////

	// fmt.Printf("env config\n")
	// fmt.Printf("config.all: %v\n", c.viper.AllSettings())
	// fmt.Printf("env[heatgo.log.level]: %v\n", c.viper.Get("heatgo.log.level"))
	// fmt.Printf("env[heatgo.log.json]: %v\n", c.viper.Get("heatgo.log.json"))
	// fmt.Printf("env[heatgo.log.out]: %v\n", c.viper.Get("heatgo.log.out"))
	// fmt.Printf("env[heatgo.log.color]: %v\n", c.viper.Get("heatgo.log.color"))
	// fmt.Println()

	////////////////////////////////////////////////////////////////
	// // .env override 에 이슈가 없는 코드 블럭
	//AlogTrace().Str(elf.LogObject, c.String()).Msg("begin to load env var")
	c.viper.SetConfigFile(".env")
	//c.viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))		// .env OK, env notOk
	if nil == c.viper.MergeInConfig() {
	}
	////////////////////////////////////////////////////////////////

	// fmt.Printf(".env config\n")
	// fmt.Printf("config.all: %v\n", c.viper.AllSettings())
	// fmt.Printf("env[heatgo.log.level]: %v\n", c.viper.Get("heatgo.log.level"))
	// fmt.Printf("env[heatgo.log.json]: %v\n", c.viper.Get("heatgo.log.json"))
	// fmt.Printf("env[heatgo.log.out]: %v\n", c.viper.Get("heatgo.log.out"))
	// fmt.Printf("env[heatgo.log.color]: %v\n", c.viper.Get("heatgo.log.color"))
	// fmt.Println()

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

// GetString get string value of key
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetBool get bool value of key
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// SetBool set bool value of key
func (c *Config) SetBool(key string, value bool) {
	c.viper.Set(key, value)
}

// GetInt get int64 value of key
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetUint get uint value of key
func (c *Config) GetUint(key string) uint {
	return c.viper.GetUint(key)
}

// GetInt32 get int32 value of key
func (c *Config) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetUint32 get Uint32 value of key
func (c *Config) GetUint32(key string) uint32 {
	return c.viper.GetUint32(key)
}

// GetInt64 get int64 value of key
func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetUint64 get uint64 value of key
func (c *Config) GetUint64(key string) uint64 {
	return c.viper.GetUint64(key)
}

// GetFloat64 get float64 value of key
func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// Set set interface value of key
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}
