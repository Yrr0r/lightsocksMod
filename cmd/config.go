package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"flag"

	"../core"
)

var (
	// 配置文件路径
	configPath string
)

var listen string
var passwd string
var remote string

//Config temporary structure
type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func init() {

	flag.StringVar(&configPath, "conf", "NO", "配置文件路径")

	flag.StringVar(&listen, "port", "EMPTY", "服务器端口")
	flag.StringVar(&passwd, "pass", "EMPTY", "密码")
	flag.StringVar(&remote, "remote", "EMPTY", "remote")

	var confgen = flag.String("gconf", "NO", "用此命令生成配置文件然后手动填入参数。")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()

	if *confgen != "NO" {
		configPath = *confgen
		content := Config{
			ListenAddr: ":7448",
			// 密码随机生成
			Password: core.RandPassword().String(),
		}
		configJson, _ := json.MarshalIndent(content, "", "	")
		err := ioutil.WriteFile(configPath, configJson, 0644)
		if err != nil {
			fmt.Errorf("保存配置到文件出错: ", configPath, err)
		}
		log.Printf("保存配置到文件成功", configPath, "\n")
		os.Exit(0)
	}
	if configPath == "NO" {
		file, err := os.Getwd()
		configPath = path.Join(file, "conf.json")
		fmt.Println(configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func (config *Config) ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("读取配置:", configPath, "\n")
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", file)
		}
	}
	if listen != "EMPTY" {
		config.ListenAddr = ":" + listen
		fmt.Println("端口：命令行参数先于配置文件")
	}
	if passwd != "EMPTY" {
		config.Password = passwd
		fmt.Println("密码：命令行参数先于配置文件")
	}
	if remote != "EMPTY" {
		config.RemoteAddr = remote
		fmt.Println("Remote：命令行参数先于配置文件")
	}

}
