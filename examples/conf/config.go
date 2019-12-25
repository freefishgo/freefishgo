package conf

import (
	"encoding/json"
	"github.com/freefishgo/freeFishGo"
	"github.com/freefishgo/freeFishGo/examples/fishgo"
	"github.com/freefishgo/freeFishGo/middlewares/mvc"
	"os"
)

var Build *freeFishGo.ApplicationBuilder

type config struct {
	*freeFishGo.Config
	WebConfig *mvc.WebConfig
}

func init() {
	Build = freeFishGo.NewFreeFishApplicationBuilder()
	conf := new(config)
	f, _ := os.Open("conf/app.conf")
	json.NewDecoder(f).Decode(conf)
	Build.Config = conf.Config
	fishgo.Mvc.Config = conf.WebConfig

}
