package conf

import (
	"encoding/json"
	"github.com/freeFishGo"
	appConfig "github.com/freeFishGo/config"
	"github.com/freeFishGo/examples/fishgo"
	"github.com/freeFishGo/middlewares/mvc"
	"os"
)

var Build *freeFishGo.ApplicationBuilder

type config struct {
	*appConfig.Config
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
