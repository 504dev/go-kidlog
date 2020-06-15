package logr_go_client

import (
	"github.com/504dev/logr/types"
	"net"
	"os"
	"time"
)

var hostname, _ = os.Hostname()
var pid = os.Getpid()
var commit = readCommit()
var tag = readTag()

type Config struct {
	Udp        string
	DashId     int
	PublicKey  string
	PrivateKey string
	Hostname   string
	Version    string
}

func (c *Config) NewLogger(logname string) (*Logger, error) {
	conn, err := net.Dial("udp", c.Udp)
	res := &Logger{
		Config:  c,
		Logname: logname,
		Prefix:  "{time} {level} ",
		Body:    "[{version}, pid={pid}, {initiator}] {message}",
		Conn:    conn,
		Counter: &Counter{
			Config:  c,
			Conn:    conn,
			Logname: logname,
			Tmp:     make(map[string]*types.Count),
		},
	}
	res.Counter.run(10 * time.Second)
	return res, err
}

func (c *Config) GetHostname() string {
	if c.Hostname != "" {
		return c.Hostname
	}
	return hostname
}

func (c *Config) GetPid() int {
	return pid
}

func (c *Config) GetVersion() string {
	if c.Version != "" {
		return c.Version
	} else if tag != "" {
		return tag
	} else if len(commit) >= 6 {
		return commit[:6]
	} else {
		return ""
	}
}
