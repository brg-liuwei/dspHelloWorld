package logger

import (
	"log"
	"os"
)

var debugLogger *log.Logger
var infoLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
	FATAL
)

var curLevel Level = DEBUG

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

func SetLevel(l Level) {
	curLevel = l
}

func GetLevel() Level {
	return curLevel
}

func Init(path string) bool {
	f, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil {
		panic(e)
	}
	debugLogger = log.New(f, "[DEBUG]", log.LstdFlags|log.Lshortfile)
	infoLogger = log.New(f, "[INFO]", log.LstdFlags)
	errorLogger = log.New(f, "[ERROR]", log.LstdFlags)
	fatalLogger = log.New(f, "[FATAL]", log.LstdFlags)
	return true
}

func Log(l Level, v ...interface{}) {
	if l >= curLevel {
		switch l {
		case DEBUG:
			go debugLogger.Println(v)
		case INFO:
			go infoLogger.Println(v)
		case ERROR:
			go errorLogger.Println(v)
		case FATAL:
			go fatalLogger.Panic(v)
		default:
			panic("Unkown log level")
		}
	}
}
