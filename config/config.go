//
//  config.go
//  mercuryx
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package config

import (
    "github.com/BurntSushi/toml"
    "github.com/Hooks-Alerts/mercuryx/logger"
)

const ServiceName    = "mercuryx"
const ServiceVersion = "1.0"

type globalConfig struct {
    Logger  LoggerConfig    `toml:"logger"`
    Server  ServerConfig    `toml:"server"`
    MySql   MySqlConfig     `toml:"my_sql"`
    Apns    ApnsConfig      `toml:"apns"`
    Gcm     GcmConfig       `toml:"gcm"`
    Chrome  GcmConfig       `toml:"chrome"`
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
    Database   string  `toml:""`
}

type ApnsConfig struct {
    PoolSize        int     `toml:"pool_size"`
    CertFile        string  `toml:"cert"`
    SandboxCertFile string  `toml:"sandbox_cert"`
}

type GcmConfig struct {
    PoolSize int `toml:"pool_size"`
}

var Logger  LoggerConfig
var Server  ServerConfig
var MySql   MySqlConfig
var Apns    ApnsConfig
var Gcm     GcmConfig
var Chrome  GcmConfig

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
        Chrome = conf.Chrome
    }
}

func initDefaultSettings() {

    // logger
    Logger.Level   = "DEBUG"
    Logger.Logfile = "mercuryx.log"

    // server
    Server.ListAddr = ""
    Server.Port = 8080

    // storage
    MySql.Host     = "localhost:3306"
    MySql.User     = "root"
    MySql.Password = "1234"
    MySql.Database = ""

    // apns
    Apns.PoolSize = 8
    Apns.CertFile = "cert.p12"
    Apns.SandboxCertFile = "cert.p12"

    // gcm
    Gcm.PoolSize = 8
    Chrome.PoolSize = 8
}
