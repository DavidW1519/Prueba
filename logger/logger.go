/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
)

var alreadyInBracketsRegex = regexp.MustCompile(`\[.*?]+`)

var (
	debugLog     = &CliLoggerFormatter{showType: "debug"}
	infoLog      = &CliLoggerFormatter{showType: "info"}
	warnLog      = &CliLoggerFormatter{showType: "warn"}
	errorLog     = &CliLoggerFormatter{showType: "error"}
	fatalLog     = &CliLoggerFormatter{showType: "fatal"}
	successLog   = &CliLoggerFormatter{showType: "success"}
	separatorLog = &SeparatorFormatter{}
)

type DefaultLogger struct {
	log    *logrus.Logger
	config *core.LoggerConfig
}

func NewDefaultLogger(log *logrus.Logger) (core.Logger, errors.Error) {
	defaultLogger := &DefaultLogger{
		log:    log,
		config: &core.LoggerConfig{},
	}
	return defaultLogger, nil
}

func (l *DefaultLogger) IsLevelEnabled(level core.LogLevel) bool {
	if l.log == nil {
		return false
	}
	return l.log.IsLevelEnabled(logrus.Level(level))
}

func (l *DefaultLogger) Log(level core.LogLevel, format string, a ...interface{}) {
	if l.IsLevelEnabled(level) {
		msg := fmt.Sprintf(format, a...)
		if l.config.Prefix != "" {
			msg = fmt.Sprintf("%s %s", l.config.Prefix, msg)
		}
		l.log.Log(logrus.Level(level), msg)
	}
}

func (l *DefaultLogger) Printf(format string, a ...interface{}) {
	//l.log.SetFormatter()
	l.Log(core.LOG_INFO, format, a...)
}

func (l *DefaultLogger) Debug(format string, a ...interface{}) {
	l.log.SetFormatter(debugLog)
	l.Log(core.LOG_DEBUG, format, a...)
}

func (l *DefaultLogger) Info(format string, a ...interface{}) {
	l.log.SetFormatter(infoLog)
	l.Log(core.LOG_INFO, format, a...)
}

func (l *DefaultLogger) Warn(err error, format string, a ...interface{}) {
	l.log.SetFormatter(warnLog)
	l.Log(core.LOG_WARN, formatMessage(err, format, a...))
}

func (l *DefaultLogger) Error(err error, format string, a ...interface{}) {
	l.log.SetFormatter(errorLog)
	l.Log(core.LOG_ERROR, formatMessage(err, format, a...))
}

func (l *DefaultLogger) SetStream(config *core.LoggerStreamConfig) {
	if config.Path != "" {
		l.config.Path = config.Path
	}
	if config.Writer != nil {
		l.log.SetOutput(config.Writer)
	}
}

func (l *DefaultLogger) GetConfig() *core.LoggerConfig {
	return &core.LoggerConfig{
		Path:   l.config.Path,
		Prefix: l.config.Prefix,
	}
}

func (l *DefaultLogger) Nested(newPrefix string) core.Logger {
	newTotalPrefix := newPrefix
	if newPrefix != "" {
		newTotalPrefix = l.createPrefix(newPrefix)
	}
	newLogger, err := l.getLogger(newTotalPrefix)
	if err != nil {
		l.Error(err, "error getting a new logger")
		return l
	}
	return newLogger
}

func (l *DefaultLogger) getLogger(prefix string) (core.Logger, errors.Error) {
	newLogrus := logrus.New()
	newLogrus.SetLevel(l.log.Level)
	newLogrus.SetFormatter(l.log.Formatter)
	newLogrus.SetOutput(l.log.Out)
	newLogger := &DefaultLogger{
		log: newLogrus,
		config: &core.LoggerConfig{
			Path:   l.config.Path,
			Prefix: prefix,
		},
	}
	return newLogger, nil
}

func (l *DefaultLogger) createPrefix(newPrefix string) string {
	newPrefix = strings.TrimSpace(newPrefix)
	alreadyInBrackets := alreadyInBracketsRegex.MatchString(newPrefix)
	if alreadyInBrackets {
		return fmt.Sprintf("%s %s", l.config.Prefix, newPrefix)
	}
	return fmt.Sprintf("%s [%s]", l.config.Prefix, newPrefix)
}

func formatMessage(err error, msg string, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)
	if err == nil {
		return msg
	}
	formattedErr := strings.ReplaceAll(err.Error(), "\n", "\n\t")
	if msg == "" {
		return formattedErr
	}
	return fmt.Sprintf("%s\n\tcaused by: %s", msg, formattedErr)
}

var _ core.Logger = (*DefaultLogger)(nil)
