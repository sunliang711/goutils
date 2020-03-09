package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitConfigLogger inits a config file configuration and log configuration
// config file at least contains the following:
// log.level		available values: trace, debug, info, warn(default), error, fatal, panic
// log.logfile
// log.showFullTime
// log.reportCaller
func InitConfigLogger() error {
	configFile := pflag.StringP("config", "c", "config.toml", "config file path")
	pflag.Parse()

	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: viper.GetBool("log.showFullTime")})

	loglevel := viper.GetString("log.level")
	log.Infof("log level: %s", loglevel)
	logrus.SetLevel(convertLevel(loglevel))

	var output io.Writer
	logfilePath := viper.GetString("log.logfile")
	if logfilePath != "" {
		handler, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithFields(log.Fields{"logfile": logfilePath, "error": err.Error()}).Fatal("Open logfile error")
		}
		log.Infof("logfile path: %s", logfilePath)
		output = io.MultiWriter(os.Stderr, handler)
	} else {
		output = os.Stderr
	}

	if viper.GetBool("log.reportCaller") {
		log.Info("log: enable report caller")
		log.SetReportCaller(true)
	}

	log.SetOutput(output)
	return nil
}

func convertLevel(l string) log.Level {
	switch l {
	case "trace":
		return log.TraceLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.ErrorLevel
	}
}
