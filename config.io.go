package elio

import "fmt"

const (
	// defaultReadBufferLen default read buffer len
	defaultReadBufferLen int = 8 * 1024
	// defaultInRange default in range
	defaultInRange int = 0
	// defaultInUrl default In url
	defaultInUrl string = ""
	// defaultInIoModel default IO model
	defaultInIoModel string = "auto"
	// defaultInCount default in count
	defaultInCount int = 1
	// defaultInWaitCount default in wait count
	defaultInWaitCount int = 512
	// defaultInNoDelay default in no delay
	defaultInNoDelay bool = false
	// defaultInRecvTimeo default in recv timeout
	defaultInRecvTimeo int = -1
	// defaultInRcvBuff default in rcv buff (KB)
	defaultInRcvBuff int = 0
	// defaultInSndBuff default in snd buff (KB)
	defaultInSndBuff int = 0
	// defaultInReusePort default in reuse port
	defaultInReusePort bool = false
)

// ConfigIo config service
type ConfigIo struct {
	ReadBufferLen int
	InURL         string
	InModel       string
	InCount       int
	InWaitCount   int
	InNoDelay     bool
	InRecvTimeo   int
	InRcvBuff     int
	InSndBuff     int
	InReusePort   bool
}

// ProvideConfigIo returns service config
func ProvideConfigIo(name string, config *Config) ConfigIo {
	c := ConfigIo{}

	c.ReadBufferLen, _ = config.GetIntOrDefault(fmt.Sprintf("%s.in.bufferlen", name), defaultReadBufferLen)
	c.InURL, _ = config.GetStringOrDefault(fmt.Sprintf("%s.in.url", name), defaultInUrl)
	c.InModel, _ = config.GetStringOrDefault(fmt.Sprintf("%s.in.iomodel", name), defaultInIoModel)
	c.InCount, _ = config.GetIntOrDefault(fmt.Sprintf("%s.in.count", name), defaultInCount)
	c.InWaitCount, _ = config.GetIntOrDefault(fmt.Sprintf("%s.in.waitcount", name), defaultInWaitCount)

	return c
}
