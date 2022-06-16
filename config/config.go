package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Server struct {
	Port string
}

type File struct {
	UploadPermission      bool
	DownloadPermission    bool
	PublicFileAllowType   []string
	PrivateFileAllowType  []string
	AlbumFileAllowType    []string
	AvatarFileAllowType   []string
	PublicFilePath        string
	PrivateFilePath       string
	TempFilePath          string
	AlbumRawFilePath      string
	AlbumCompressFilePath string
	AvatarFilePath        string
	RobotsTxtPath         string
}

type MySQL struct {
	Name            string
	Args            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime int // seconds
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Admin struct {
	SessionMaxAge     int
	DefaultAvatarName string
}

type Wechat struct {
	AppId     string
	AppSecret string
}

type Turing struct {
	TuringApiUrl string
	TuringKey    string
}

type Config struct {
	Server Server
	MySQL  MySQL
	Redis  Redis
	Admin  Admin
	File   File
	Wechat Wechat
	Turing Turing
}

var config *Config

func GetConfig() *Config {
	return config
}

func LoadConfig() {
	file, err := os.Open("conf.yaml")
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(
		file)
	if err != nil {
		panic(err)
	}

	cfg := Config{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("load config:", cfg)
	config = &cfg
}
