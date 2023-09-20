// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	gormUtils "gorm.io/gorm/utils"
	"os"
	"time"
)

var logger *Logger

type Logger struct {
	*logrus.Logger
	Pid int
}

func init() {
	logger = loggerInit("", 6)
}
func NewPrivateLog(moduleName string, logLevel uint32) {
	logger = loggerInit(moduleName, logLevel)
}
func IsNil() bool {
	if logger != nil {
		return false
	}
	return true
}

func loggerInit(moduleName string, logLevel uint32) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.Level(logLevel))
	//Close std console output when running on the server
	if moduleName != "" {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err.Error())
		}
		writer := bufio.NewWriter(src)
		logger.SetOutput(writer)
	}

	//Log Console Print Style Setting
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	//File name and line number display hook
	logger.AddHook(newFileHook())
	//Log file segmentation hook when running on the server
	if moduleName != "" {
		hook := NewLfsHook(time.Duration(24)*time.Hour, 3, moduleName)
		logger.AddHook(hook)
	}
	return &Logger{
		logger,
		os.Getpid(),
	}
}
func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	return lfsHook
}
func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, moduleName string) *rotatelogs.RotateLogs {
	if moduleName != "" {
		moduleName = moduleName + "."
	}
	writer, err := rotatelogs.New(
		"../logs/"+moduleName+level+"."+"%Y-%m-%d",
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

// internal method
func argsHandle(OperationID string, fields logrus.Fields, args []interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
	fields["OperationID"] = OperationID
	fields["PID"] = logger.Pid
}

func NewInfo(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}

func NewError(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}
func NewDebug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}
func NewWarn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}

func Info(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}

func Error(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}
func Debug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}
func Warn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}

type SqlLogger struct {
	LogLevel                  gormLogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func NewSqlLogger(logLevel gormLogger.LogLevel, ignoreRecordNotFoundError bool, slowThreshold time.Duration) *SqlLogger {
	return &SqlLogger{
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		SlowThreshold:             slowThreshold,
	}
}

func (l *SqlLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = logLevel
	return &newLogger
}

func (SqlLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logrus.Info(ctx, msg, args)
}

func (SqlLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logrus.Warn(ctx, msg, nil, args)
}

func (SqlLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logrus.Error(ctx, msg, nil, args)
}

func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logrus.Error(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			logrus.Error(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logrus.Warn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			logrus.Warn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	case l.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logrus.Debug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			logrus.Debug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	}
}
