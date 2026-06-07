package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = log.LstdFlags
)

const (
	LevelDebug = (iota + 1) * 100
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelAlert
)

const (
	namePrefix = "LEVEL"
	levelDepth = 4
)

type Logger struct {
	level  int
	writer io.Writer
	logger *log.Logger
}

func MappingLevelWithName(level string) int {
	ret := LevelInfo

	switch strings.ToUpper(level) {
	case "ALL":
		fallthrough
	case "DEBUG":
		ret = LevelDebug

	case "WARN":
		fallthrough
	case "WARNING":
		ret = LevelWarning

	case "ERROR":
		ret = LevelError

	case "FATAL":
		ret = LevelFatal
	}

	return ret
}

func mappingLevel(level int) string {

	var result string
	switch level {
	case LevelInfo:
		result = "[INFO]"
	case LevelFatal:
		result = "[FATAL]"
	case LevelWarning:
		result = "[WARN]"
	case LevelError:
		result = "[ERROR]"
	case LevelAlert:
		result = " MB_FATAL "
	default:
		result = "[DEBUG]"
	}
	return result
}

func New(out io.Writer, prefix string, flag, level int) *Logger {

	lvl := os.Getenv("MB_LOG_LEVEL")
	if lvl != "" {
		level, _ = strconv.Atoi(lvl)
	}

	return &Logger{
		level:  level,
		writer: out,
		logger: log.New(out, prefix, flag),
	}
}

func (l *Logger) Flags() int {
	return l.logger.Flags()
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) Prefix() string {
	return l.logger.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Logger) Level() int {
	return l.level
}

// SetLevel is not locked.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) output(level, calldepth int, s string) error {
	if l == std {
		calldepth++
	}
	return l.logger.Output(calldepth, mappingLevel(level)+" "+s)
}

func (l *Logger) Err(level, calldepth int, err error) error {
	if err != nil && level >= l.level {
		return l.output(level, calldepth, err.Error())
	}
	return nil
}

func (l *Logger) Output(level, calldepth int, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprint(a...))
	}
	return nil
}

func (l *Logger) Outputf(level, calldepth int, format string, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprintf(format, a...))
	}
	return nil
}

func (l *Logger) Outputln(level, calldepth int, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprintln(a...))
	}
	return nil
}

func (l *Logger) ErrDebug(err error) {
	l.Err(LevelDebug, levelDepth, err)
}

func (l *Logger) ErrInfo(err error) {
	l.Err(LevelInfo, levelDepth, err)
}

func (l *Logger) ErrWarning(err error) {
	l.Err(LevelWarning, levelDepth, err)
}

func (l *Logger) ErrError(err error) {
	l.Err(LevelError, levelDepth, err)
}

func (l *Logger) ErrFatal(err error) {
	if err != nil {
		l.Err(LevelFatal, levelDepth, err)
		os.Exit(1)
	}
}

func (l *Logger) ErrAlert(err error) {
	if err != nil {
		l.Err(LevelAlert, levelDepth, err)
	}
}

func (l *Logger) Debug(a ...interface{}) {
	l.Output(LevelDebug, levelDepth, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.Output(LevelInfo, levelDepth, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.Output(LevelWarning, levelDepth, a...)
}

func (l *Logger) Error(a ...interface{}) {
	l.Output(LevelError, levelDepth, a...)
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Output(LevelFatal, levelDepth, a...)
	os.Exit(1)
}

func (l *Logger) Alert(a ...interface{}) {
	l.Output(LevelAlert, levelDepth, a...)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Outputf(LevelDebug, levelDepth, format, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Outputf(LevelInfo, levelDepth, format, a...)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Outputf(LevelWarning, levelDepth, format, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Outputf(LevelError, levelDepth, format, a...)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Outputf(LevelFatal, levelDepth, format, a...)
	os.Exit(1)
}

func (l *Logger) Alertf(format string, a ...interface{}) {
	l.Outputf(LevelAlert, levelDepth, format, a...)
}

func (l *Logger) Debugln(a ...interface{}) {
	l.Outputln(LevelDebug, levelDepth, a...)
}

func (l *Logger) Infoln(a ...interface{}) {
	l.Outputln(LevelInfo, levelDepth, a...)
}

func (l *Logger) Warningln(a ...interface{}) {
	l.Outputln(LevelWarning, levelDepth, a...)
}

func (l *Logger) Errorln(a ...interface{}) {
	l.Outputln(LevelError, levelDepth, a...)
}

func (l *Logger) Fatalln(a ...interface{}) {
	l.Outputln(LevelFatal, levelDepth, a...)
	os.Exit(1)
}

func (l *Logger) Alertln(a ...interface{}) {
	l.Outputln(LevelAlert, levelDepth, a...)
}

var std = New(os.Stderr, "", LstdFlags, LevelInfo)

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func GetOutput() io.Writer {
	return std.writer
}

func Flags() int {
	return std.Flags()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func Prefix() string {
	return std.Prefix()
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func Level() int {
	return std.Level()
}

// SetLevel is not locked.
func SetLevel(level int) {
	std.SetLevel(level)
}

func Err(level, calldepth int, err error) error {
	return std.Err(level, calldepth, err)
}

func Output(level, calldepth int, a ...interface{}) error {
	return std.Output(level, calldepth, a...)
}

func Outputf(level, calldepth int, format string, a ...interface{}) error {
	return std.Outputf(level, calldepth, format, a...)
}

func Outputln(level, calldepth int, a ...interface{}) error {
	return std.Outputln(level, calldepth, a...)
}

func ErrDebug(err error) {
	std.ErrDebug(err)
}

func ErrInfo(err error) {
	std.ErrInfo(err)
}

func ErrWarning(err error) {
	std.ErrWarning(err)
}

func ErrError(err error) {
	std.ErrError(err)
}

func ErrFatal(err error) {
	std.ErrFatal(err)
}

func ErrAlert(err error) {
	std.ErrAlert(err)
}

func Debug(a ...interface{}) {
	std.Debug(a...)
}

func Info(a ...interface{}) {
	std.Info(a...)
}

func Warning(a ...interface{}) {
	std.Warning(a...)
}

func Error(a ...interface{}) {
	std.Error(a...)
}

func Fatal(a ...interface{}) {
	std.Fatal(a...)
}

func Alert(a ...interface{}) {
	std.Alert(a...)
}

func Debugf(format string, a ...interface{}) {
	std.Debugf(format, a...)
}

func Infof(format string, a ...interface{}) {
	std.Infof(format, a...)
}

func Warningf(format string, a ...interface{}) {
	std.Warningf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	std.Errorf(format, a...)
}

func Fatalf(format string, a ...interface{}) {
	std.Fatalf(format, a...)
}

func Alertf(format string, a ...interface{}) {
	std.Alertf(format, a...)
}

func Debugln(a ...interface{}) {
	std.Debugln(a...)
}

func Infoln(a ...interface{}) {
	std.Infoln(a...)
}

func Warningln(a ...interface{}) {
	std.Warningln(a...)
}

func Errorln(a ...interface{}) {
	std.Errorln(a...)
}

func Fatalln(a ...interface{}) {
	std.Fatalln(a...)
}

func Alertln(a ...interface{}) {
	std.Alertln(a...)
}
