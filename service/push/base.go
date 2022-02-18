package push

import (
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
)

var pushMap map[string]Push

func init() {
	initPush()
}

// endpoint represents a kind of subscription
type endpoint struct {
	Name  string `config:"name"`
	Type  string `config:"type"`
	URL   string `config:"url"`
	Token string `config:"token"`
}

// initPush bind endpoints with config file
func initPush() {
	conf := config.NewWithOptions("push", func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
		opt.ParseEnv = true
	})
	conf.AddDriver(toml.Driver)
	err := conf.LoadFiles("config/push.toml")
	if err != nil {
		panic(err)
	}

	// Load config file
	size := conf.Get("global.size").(int64)
	endpoints := make([]endpoint, size)
	conf.BindStruct("endpoints", &endpoints)

	// Load token or key here
	if !conf.Get("global.action").(bool) {
		return
	}
	conf.LoadOSEnv([]string{TurboEnv, PushDeerEnv}, false)
	turbo := conf.Get(TurboEnv).(string)
	turboList, tId := strings.Split(turbo, ","), 0
	pushDeer := conf.Get(PushDeerEnv).(string)
	pushDeerList, pId := strings.Split(pushDeer, ","), 0
	for _, endpoint := range endpoints {
		switch endpoint.Type {
		case TurboName:
			endpoint.Token = turboList[tId]
			tId++
			addPush(endpoint.Name, &TurboPush{
				endpoint: endpoint,
			})
		case PushDeerName:
			endpoint.Token = pushDeerList[pId]
			pId++
			addPush(endpoint.Name, &PushDeerPush{
				endpoint: endpoint,
			})
		}
	}
}

// Data represents data needed for push
type Data struct {
	Title   string
	Content string
}

// Push contain all info needed for push action
type Push interface {
	Submit(data Data) error
}

func addPush(name string, push Push) {
	if pushMap == nil {
		pushMap = make(map[string]Push)
	}
	pushMap[name] = push
}

func NewPush(name string) Push {
	if pull, ok := pushMap[name]; ok {
		return pull
	} else {
		return RawPush{}
	}
}

type RawPush struct {
	endpoint
}

func (RawPush) Submit(data Data) error {
	return nil
}
