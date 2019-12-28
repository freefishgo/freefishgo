package conf

import (
	"encoding/json"
	"github.com/freefishgo/freeFishGo/middlewares/mvc"
	"github.com/freefishgo/freefishgo"
	"os"
)

//var Build *freeFishGo.ApplicationBuilder

type config struct {
	*freeFishGo.Config
	WebConfig *mvc.MvcWebConfig
}

func init() {
	conf := new(config)
	f, _ := os.Open("conf/app.conf")
	json.NewDecoder(f).Decode(conf)
	freeFishGo.DefaultConfig = conf.Config
	mvc.DefaultMvcWebConfig = conf.WebConfig

}
