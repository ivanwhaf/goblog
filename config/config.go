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
	AlbumRawFilePath      string
	AlbumCompressFilePath string
	AvatarFilePath        string
	RobotsTxtPath         string
}

type DB struct {
	Name string
	Args string
}

type Config struct {
	Server Server
	DB     DB
	Admin  Admin
	File   File
}

type Admin struct {
	SessionMaxAge     int
	DefaultAvatarName string
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
