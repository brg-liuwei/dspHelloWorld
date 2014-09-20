package manager

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	redis "github.com/gosexy/redis"
	"logger"
)

var ErrNoCmd error = errors.New("No Cmd")

type Commander struct {
	host string
	port uint
	rkey string // key of redis queue
	cli  *redis.Client
	err  error
}

func NewCommander(host string, port uint, rkey string) *Commander {
	c := redis.New()
	err := c.Connect(host, port)
	return &Commander{host: host,
		port: port,
		rkey: rkey,
		cli:  c,
		err:  err}
}

func (c *Commander) Repair() bool {
	if c.err == nil {
		return true
	}
	if c.err = c.cli.Connect(c.host, c.port); c.err != nil {
		return false
	}
	return true
}

func (c *Commander) GetCmd() (cmd string, err error) {
	if !c.Repair() {
		return "", c.err
	}
	ss, e := c.cli.BLPop(600, c.rkey)
	switch e {
	case nil:
		cmd = ss[0]
	case redis.ErrNilReply:
		err = ErrNoCmd
	default:
		err = e
		c.Repair()
	}
	return
}

func CommanderRoutine(host string, port uint, rkey string, ch chan<- string) {
	c := NewCommander(host, port, rkey)
	for {
		cmd, err := c.GetCmd()
		switch err {
		case ErrNoCmd:
			time.Sleep(5 * time.Second)
		case nil:
			ch <- cmd
		default:
			logger.Log(logger.ERROR, "redis connect err:", err)
			c.Repair()
			time.Sleep(1 * time.Second)
		}
	}
}

type CmdType int

const (
	NilCmd   CmdType = iota // 0
	AddOrder                // 1
	AddAd                   // 2
	ModOrder                // 3
	ModAd                   // 4
	DelOrder                // 5
	DelAd                   // 6
)

type Command struct {
	Ctype    CmdType
	Cversion string
	Data     interface{}
}

func NewCommand() *Command {
	return &Command{Ctype: NilCmd,
		Cversion: "NoVersion"}
}

func (c *Command) Parse(jsonCmd string) bool {
	type Fmt struct {
		Oper_type string
		Fmt_ver   string
		Data      interface{}
	}
	var cmd Fmt
	dec := json.NewDecoder(strings.NewReader(jsonCmd))
	err := dec.Decode(&cmd)
	if err != nil {
		logger.Log(logger.ERROR, "decode cmd err: ", jsonCmd, err)
		return false
	}
	switch cmd.Oper_type {
	case "1":
		c.Ctype = AddOrder
	case "2":
		c.Ctype = AddAd
	case "3":
		c.Ctype = ModOrder
	case "4":
		c.Ctype = ModAd
	case "5":
		c.Ctype = DelOrder
	case "6":
		c.Ctype = DelAd
	default:
		logger.Log(logger.ERROR, "cmd type error: ", jsonCmd, cmd.Oper_type)
		return false
	}
	c.Cversion = cmd.Fmt_ver
	c.Data = cmd.Data
	return true
}

func (c *Command) Execute()
