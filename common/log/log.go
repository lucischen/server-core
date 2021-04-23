package log

import (
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/francoispqt/onelog"
)

const (
	// 起始的 caller 層數
	startCallerLevel = 4
	// DEBUG 除錯
	DEBUG string = "debug"
	// INFO 資訊
	INFO string = "info"
	// ERROR 錯誤
	ERROR string = "error"
)

var levelMap = map[string]uint8{
	DEBUG: onelog.DEBUG | onelog.INFO | onelog.ERROR,
	INFO:  onelog.INFO | onelog.ERROR,
	ERROR: onelog.ERROR,
}

// Log history
type Log struct {
	logger      *onelog.Logger
	callerLevel int
	isMock      bool
}

// Setting 初始化設定
type Setting struct {
	Level         string
	CallersLevels int
	Writer        io.Writer
}

// GetLogConfig 創建設定檔案
func GetLogConfig() *Setting {
	return &Setting{Level: ERROR}
}

// NewLog 建立一個 Log
func NewLog(config *Setting) *Log {
	lv := strings.ToLower(config.Level)
	log := &Log{
		callerLevel: config.CallersLevels,
		logger: onelog.New(
			config.Writer,
			levelMap[lv],
		),
	}

	log.logger.Hook(func(e onelog.Entry) {
		e.String("time", time.Now().Format(time.RFC3339))
	})

	return log
}

func (l *Log) SwitchToMockMode() {
	l.isMock = true
}

// Debug Log debug
func (l *Log) Debug(msg string) {
	if l.isMock {
		return
	}
	l.logger.Debug(msg)
}

// Info Log info
func (l *Log) Info(msg string) {
	if l.isMock {
		return
	}
	l.logger.Info(msg)
}

// Error Log error
func (l *Log) Error(msg string) {
	if l.isMock {
		return
	}
	l.logger.Error(msg)
}

// InfoWithFields 添加額外資訊 Log info
func (l *Log) InfoWithFields(msg string, fields func(onelog.Entry)) {
	if l.isMock {
		return
	}
	l.logger.InfoWithFields(msg, fields)
}

// ErrorWithFields 添加額外資訊 Log error
func (l *Log) ErrorWithFields(msg string, fields func(onelog.Entry)) {
	if l.isMock {
		return
	}
	l.logger.ErrorWithFields(msg, fields)
}

func (l *Log) ErrorMsg(title string) LogField {
	if l.isMock {
		return &Mock{}
	}

	return &Error{
		title:       title,
		logger:      l.logger,
		callerLevel: l.callerLevel,
		m:           map[string]string{},
	}
}

func (l *Log) InfoMsg(title string) LogField {
	if l.isMock {
		return &Mock{}
	}

	return &Info{
		title:       title,
		logger:      l.logger,
		callerLevel: l.callerLevel,
		m:           map[string]string{},
	}
}

func (l *Log) DebugMsg(title string) LogField {
	if l.isMock {
		return &Mock{}
	}

	return &Debug{
		title:       title,
		logger:      l.logger,
		callerLevel: l.callerLevel,
		m:           map[string]string{},
	}
}

// GenerateCallerList 生成 caller 清單
func (l *Log) GenerateCallerList() string {
	var callers strings.Builder

	for i := startCallerLevel; ; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok || i == getCallerLevels(l.callerLevel) {
			break
		}

		var caller strings.Builder
		caller.WriteString(file)
		caller.WriteString(":")
		caller.WriteString(strconv.Itoa(line))

		callers.WriteString(caller.String())
		callers.WriteString(" ")
	}

	return callers.String()
}

func getCallerLevels(callerLevel int) int {
	if callerLevel == -1 {
		return callerLevel
	}
	return startCallerLevel + (callerLevel - 1)
}
