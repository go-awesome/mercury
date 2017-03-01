//
//  config/config.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ortuman/mercury/logger"
)

const ServiceName = "mercury"
const ServiceVersion = "1.0"

type globalConfig struct {
	Logger  LoggerConfig  `toml:"logger"`
	Server  ServerConfig  `toml:"server"`
	Redis   RedisConfig   `toml:"redis"`
	Apns    ApnsConfig    `toml:"apns"`
	Gcm     GcmConfig     `toml:"gcm"`
	Safari  SafariConfig  `toml:"safari"`
	Chrome  WebPushConfig `toml:"chrome"`
	Firefox WebPushConfig `toml:"firefox"`
}

type LoggerConfig struct {
	Level   string `toml:"level"`
	Logfile string `toml:"logpath"`
}

type ServerConfig struct {
	BindAddress          string `toml:"bind_address"`
	Port                 int    `toml:"port"`
	UnregisteredCallback string `toml:"unregistered_callback"`
}

type RedisConfig struct {
	Host string `toml:"host"`
}

type ApnsConfig struct {
	MaxConn         uint32 `toml:"max_conn"`
	CertFile        string `toml:"cert"`
	SandboxCertFile string `toml:"sandbox_cert"`
}

type GcmConfig struct {
	MaxConn uint32 `toml:"max_conn"`
	ApiKey  string `toml:"api_key"`
}

type SafariConfig struct {
	MaxConn  uint32 `toml:"max_conn"`
	CertFile string `toml:"cert"`
}

type WebPushConfig struct {
	MaxConn    uint32 `toml:"max_conn"`
	Subject    string `toml:"sub"`
	PublicKey  string `toml:"public_key"`
	PrivateKey string `toml:"private_key"`
}

var Logger LoggerConfig
var Server ServerConfig
var Redis RedisConfig
var Apns ApnsConfig
var Gcm GcmConfig
var Safari SafariConfig
var Chrome WebPushConfig
var Firefox WebPushConfig

func init() {
	initDefaultSettings()
}

func Load(cfgFile string) {
	var conf globalConfig
	if _, err := toml.DecodeFile(cfgFile, &conf); err != nil {
		logger.Warnf("main: couldn't load config file '%s': %v", cfgFile, err)
	} else {
		Logger = conf.Logger
		Server = conf.Server
		Redis = conf.Redis
		Apns = conf.Apns
		Gcm = conf.Gcm
		Safari = conf.Safari
		Chrome = conf.Chrome
		Firefox = conf.Firefox
	}
}

func initDefaultSettings() {

	// logger
	Logger.Level = "DEBUG"
	Logger.Logfile = "mercury.log"

	// server
	Server.BindAddress = ""
	Server.Port = 8080

	// storage
	Redis.Host = "localhost:6379"

	// apns
	Apns.MaxConn = 16
	Apns.CertFile = "cert.p12"
	Apns.SandboxCertFile = "cert.p12"

	// gcm
	Gcm.MaxConn = 16

	// safari
	Safari.MaxConn = 16

	// chrome
	Chrome.MaxConn = 16

	// firefox
	Firefox.MaxConn = 16
}
