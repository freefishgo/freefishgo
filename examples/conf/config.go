package conf

import (
	"encoding/json"
	"github.com/freefishgo/freefishgo"
	"github.com/freefishgo/freefishgo/middlewares/mvc"
	"os"
)

//var Build *freeFishGo.ApplicationBuilder

type config struct {
	*freefishgo.Config
	WebConfig *mvc.MvcWebConfig
}

func init() {
	conf := new(config)
	f, err := os.Open("conf/app.conf")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	json.NewDecoder(f).Decode(conf)
	freefishgo.DefaultConfig = conf.Config
	mvc.DefaultMvcWebConfig = conf.WebConfig

}
