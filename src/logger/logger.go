package logger

import (
	"encoding/json"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "Unkown"
	}
}

type Log struct {
	curLevel    Level
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

var WinLog *Log
var ClickLog *Log
var DisplayLog *Log
var BidLog *Log
var GiveupLog *Log

func NewLog(path string) *Log {
	f, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil {
		panic(e)
	}
	return &Log{
		curLevel:    DEBUG,
		debugLogger: log.New(f, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		infoLogger:  log.New(f, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(f, "[ERROR] ", log.LstdFlags),
		fatalLogger: log.New(f, "[FATAL] ", log.LstdFlags),
	}
}

func (log *Log) SetLevel(l Level) {
	log.curLevel = l
}

func (log *Log) GetLevel() Level {
	return log.curLevel
}

func (log *Log) Log(l Level, v ...interface{}) {
	if l >= log.curLevel {
		switch l {
		case DEBUG:
			log.debugLogger.Println(v)
		case INFO:
			log.infoLogger.Println(v)
		case ERROR:
			log.errorLogger.Println(v)
		case FATAL:
			log.fatalLogger.Panic(v)
		default:
			panic("Unkown log level")
		}
	}
}

func (log *Log) JsonLog(l Level, m map[string]interface{}) bool {
	if b, err := json.Marshal(m); err != nil {
		return false
	} else {
		log.Log(l, string(b))
		return true
	}
}
