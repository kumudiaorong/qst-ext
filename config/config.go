package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppAttr struct {
	// prompt of extension
	Name string `yaml:"prompt"`
	// icon of extension
	Icon string `yaml:"icon"`
	// dir of extension
	Dir string `yaml:"path"`
	// executable file of extension
	Exec string `yaml:"exec"`
}

type FileAttr struct {
	// QstPath is the path of qst
	Apps map[uint32]*AppAttr `yaml:"exts"`
}
type App struct {
	*os.Process
}
type AppStatus struct {
	Attr       FileAttr
	MotifyTime int64
	RunStat    map[uint32]*App //区分本地和远程
}

var (
	// decoder is the decoder of config file
	decoder *yaml.Decoder
	// Config is the config of qst
	Status AppStatus
)

func init() {
	var config FileAttr
	var cfgPath string
	home := os.Getenv("HOME")
	if home == "" {
		log.Println("Can't Find Home Path")
	} else {
		cfgPath = home + "/.config/qst/ext-appsearcher.yaml"
		ifs, err := os.Open(cfgPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(home+"/.config/qst", 0755)
				if err != nil {
					log.Printf("Can't Create Config Dir:%v\n", err)
				}
				ifs, err = os.Create(cfgPath)
				if err != nil {
					log.Printf("Can't Create Config File:%v\n", err)
				} else {
					defer ifs.Close()
				}
			} else {
				log.Printf("Can't Open Config File:%v\n", err)
			}
		} else {
			defer ifs.Close()
		}
		if ifs != nil {
			decoder = yaml.NewDecoder(ifs)
			decoder.Decode(&config)
		}
	}
	Status.RunStat = make(map[uint32]*App)
}
