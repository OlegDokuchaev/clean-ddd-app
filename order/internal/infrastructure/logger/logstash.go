package logger

import (
	"fmt"
	"net"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

func createLogstashAddress(config *Config) string {
	return fmt.Sprintf("%s:%s", config.LogstashHost, config.LogstashPort)
}

func createLogstashConn(config *Config) (net.Conn, error) {
	address := createLogstashAddress(config)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewLogstash(config *Config) (logrus.Hook, error) {
	conn, err := createLogstashConn(config)
	if err != nil {
		return nil, err
	}

	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": config.ServiceName}))
	return hook, nil
}
