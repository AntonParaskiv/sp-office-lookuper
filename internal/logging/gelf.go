package logging

import (
	"fmt"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
)

func (l *Logger) InitGelf() (*Logger, error) {
	gelfAddr := fmt.Sprintf("%s:%d", l.host, l.port)
	hook := graylog.NewAsyncGraylogHook(gelfAddr, map[string]interface{}{})

	hook.Level = logrus.InfoLevel
	if l.level != 0 {
		hook.Level = logrus.Level(l.level)
	}
	l.logger.AddHook(hook)
	l.flusherHooks = append(l.flusherHooks, hook)

	return l, nil
}
