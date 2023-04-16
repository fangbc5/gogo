package config

import (
	"gogo/constant"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func LoadLocalConfig(path string, config interface{}) {
	vLocal := viper.New()
	vLocal.SetDefault("Consul.Host", "localhost")
	vLocal.SetDefault("Consul.Port", "8500")
	vLocal.SetDefault("Consul.Config", "common")
	vLocal.AddConfigPath("../" + path)
	vLocal.SetConfigName("app")
	vLocal.SetConfigType("json")
	if err := vLocal.ReadInConfig(); err != nil {
		log.Fatalln("load local config fail...")
		return
	}

	if err := vLocal.Unmarshal(config); err != nil {
		log.Fatalln("local config serialization fail...")
		return
	}
	vLocal.WatchConfig()
	vLocal.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config file change" + in.Name)
		refreshLocalConfig(vLocal, config)
	})
}

func refreshLocalConfig(v *viper.Viper, config interface{}) {
	err := v.Unmarshal(config)
	if err != nil {
		log.Println("refresh local config serialization fail...")
		return
	}
}

// 从consul读取配置
func LoadRemoteConfig(paths string, config interface{}) {
	pathArr := strings.Split(paths, ",")
	for _, path := range pathArr {
		vRemote := viper.New()
		vRemote.AddRemoteProvider("consul", "http://localhost:8500", constant.Namespace+"/"+path)
		vRemote.SetConfigType("json")
		if err := vRemote.ReadRemoteConfig(); err != nil {
			log.Fatalln("load remote config fail...")
			return
		}
		if err := vRemote.Unmarshal(config); err != nil {
			log.Fatalln("remote config serialization fail...")
			return
		}
		// go func() {
		// 	for {
		// 		time.Sleep(time.Second * 5) // delay after each request

		// 		// currently, only tested with etcd support
		// 		err := vRemote.WatchRemoteConfig()
		// 		if err != nil {
		// 			log.Printf("unable to read remote config: %v", err)
		// 			continue
		// 		}

		// 		// unmarshal new config into our runtime config struct. you can also use channel
		// 		// to implement a signal to notify the system of the changes
		// 		vRemote.Unmarshal(config)
		// 	}
		// }()
		// vRemote.WatchRemoteConfig()
		// vRemote.OnConfigChange(func(in fsnotify.Event) {
		// 	log.Println("config file change" + in.Name)
		// 	refreshRemoteConfig(vRemote, config)
		// })
	}
}

// func refreshRemoteConfig(v *viper.Viper, config interface{}) {
// 	err := v.Unmarshal(config)
// 	if err != nil {
// 		log.Println("refresh remote config serialization fail...")
// 		return
// 	}
// }
