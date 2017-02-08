//
//  config.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package config

import (
    "github.com/BurntSushi/toml"
    "github.com/ortuman/mercury/logger"
)

const ServiceName    = "mercury"
const ServiceVersion = "1.0"

type globalConfig struct {
    Logger  LoggerConfig    `toml:"logger"`
    Server  ServerConfig    `toml:"server"`
    MySql   MySqlConfig     `toml:"my_sql"`
    Apns    ApnsConfig      `toml:"apns"`
    Gcm     GcmConfig       `toml:"gcm"`
}

type LoggerConfig struct {
    Level   string  `toml:"level"`
    Logfile string  `toml:"logpath"`
}

type ServerConfig struct {
    ListAddr      string  `toml:"bind_address"`
    Port          int     `toml:"port"`
}

type MySqlConfig struct {
    Host       string  `toml:"host"`
    User       string  `toml:"user"`
    Password   string  `toml:"pass"`
}

type ApnsConfig struct {
    PoolSize        uint32  `toml:"pool_size"`
    CertFile        string  `toml:"cert"`
    SandboxCertFile string  `toml:"sandbox_cert"`
}

type GcmConfig struct {
    PoolSize uint32 `toml:"pool_size"`
}

var Logger  LoggerConfig
var Server  ServerConfig
var MySql   MySqlConfig
var Apns    ApnsConfig
var Gcm     GcmConfig

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
        MySql  = conf.MySql
        Apns   = conf.Apns
        Gcm    = conf.Gcm
    }
}

func initDefaultSettings() {

    // logger
    Logger.Level   = "DEBUG"
    Logger.Logfile = "mercury.log"

    // server
    Server.ListAddr = ""
    Server.Port = 8080

    // storage
    MySql.Host     = "localhost:3306"
    MySql.User     = "root"
    MySql.Password = "1234"

    // apns
    Apns.PoolSize = 16
    Apns.CertFile = "cert.p12"
    Apns.SandboxCertFile = "cert.p12"

    // gcm
    Gcm.PoolSize = 16
}
