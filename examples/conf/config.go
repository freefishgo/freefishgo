package conf

import (
	"encoding/json"
	"github.com/freefishgo/freeFish"
	"github.com/freefishgo/freeFish/examples/fishgo"
	"github.com/freefishgo/freeFish/middlewares/mvc"
	"os"
)

var Build *freeFish.ApplicationBuilder

type config struct {
	*freeFish.Config
	WebConfig *mvc.WebConfig
}

func init() {
	Build = freeFish.NewFreeFishApplicationBuilder()
	conf := new(config)
	f, _ := os.Open("conf/app.conf")
	json.NewDecoder(f).Decode(conf)
	Build.Config = conf.Config
	fishgo.Mvc.Config = conf.WebConfig

}
