package log

import (
	"fmt"
	"log"
	"time"
	//"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

//===========================================

type Logger struct {
	Name string
}

func NewLogger(name string) *Logger {
	return &Logger{Name: buildLoggerName(name)}
}

func buildLoggerName(name string) string {
	return fmt.Sprintf("[%s] ", name)
}

var defaultlLogger = NewLogger("")

func SetDefaultLoggerName(name string) {
	defaultlLogger.Name = buildLoggerName(name)
}

//===========================================

func (logger *Logger) Debug(v ...interface{}) {
	log.Print(now("DEBUG"), logger.Name, fmt.Sprint(v...))
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	log.Print(now("DEBUG"), logger.Name, fmt.Sprintf(format, v...))
}

func (logger *Logger) Info(v ...interface{}) {
	log.Print(now("INFO"), logger.Name, fmt.Sprint(v...))
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	log.Print(now("INFO"), logger.Name, fmt.Sprintf(format, v...))
}

func (logger *Logger) Warning(v ...interface{}) {
	log.Print(now("WARNING"), logger.Name, fmt.Sprint(v...))
}

func (logger *Logger) Warningf(format string, v ...interface{}) {
	log.Print(now("WARNING"), logger.Name, " [WARNING] ", fmt.Sprintf(format, v...))
}

func (logger *Logger) Error(v ...interface{}) {
	log.Print(now("ERROR"), logger.Name, fmt.Sprint(v...))
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	log.Print(now("ERROR"), logger.Name, fmt.Sprintf(format, v...))
}

func (logger *Logger) Fatal(v ...interface{}) {
	log.Fatal(now("FATAL"), logger.Name, fmt.Sprint(v...))
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	log.Fatal(now("FATAL"), logger.Name, fmt.Sprintf(format, v...))
}

//======================================

func Debug(v ...interface{}) {
	defaultlLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultlLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultlLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultlLogger.Infof(format, v...)
}

func Warning(v ...interface{}) {
	defaultlLogger.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	defaultlLogger.Warningf(format, v...)
}

func Error(v ...interface{}) {
	defaultlLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultlLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	defaultlLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultlLogger.Fatalf(format, v...)
}

//=====================================

func now(level string) string {
	return fmt.Sprintf("[%s %s]", time.Now().Format("2006-01-02 15:04:05"), level)
}
