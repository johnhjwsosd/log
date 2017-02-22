package logger

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/mgo.v2-unstable/bson"
)

//Logger ...
type Logger struct {
	msName string
	ls     LogStorage
}

//LogStorage ...
type LogStorage interface {
	Write(content *LogContent) (err error)
	Read(conditions bson.M) ([]LogContent, error)
}

//LogContent ...
type LogContent struct {
	Level   string
	MsName  string
	Time    string
	Content bson.M
}

var newLogStorageFuncs = make(map[string]func(stHost string, stPort int, stName, appName string) (ls LogStorage, err error))

//NewLogger ...
func NewLogger(stType, stHost string, stPort int, stName, appName, msName string) (logger *Logger, err error) {
	if newLSFunc, e := newLogStorageFuncs[stType]; e {
		fmt.Printf("\n[Log save at %s -> %s:%d log->%s]\n", stType, stHost, stPort, appName)
		ls, errNewLS := newLSFunc(stHost, stPort, stName, appName)
		if errNewLS != nil {
			fmt.Printf("Error:\n%s\n", errNewLS.Error())
			// err = errNewLS
			// return
		}
		fmt.Println()
		logger = &Logger{
			msName: msName,
			ls:     ls}
		return
	}
	err = errors.New("Log storage type not found")
	return
}

//LogStorageRegister ...
func LogStorageRegister(stType string, newFunc func(string, int, string, string) (LogStorage, error)) error {
	if _, e := newLogStorageFuncs[stType]; e {
		return errors.New("Log storage type already registered")
	}
	newLogStorageFuncs[stType] = newFunc
	return nil
}

//Trace ...
func (l *Logger) Trace(content interface{}) {
	now := time.Now()
	fmt.Printf("[TRACE] %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
}

//Debug ...
func (l *Logger) Debug(content interface{}) {
	now := time.Now()
	fmt.Printf("[DEBUG] %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
}

//Info ...
func (l *Logger) Info(content interface{}) {
	now := time.Now()
	fmt.Printf("[INFO]  %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
	if l.ls != nil {
		go l.ls.Write(&LogContent{"INFO", l.msName, now.Format("2006-01-02 15:04:05"), toBsonM(content)})
	}
}

//Warn ...
func (l *Logger) Warn(content interface{}) {
	now := time.Now()
	fmt.Printf(
		"[WARN]  %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
	if l.ls != nil {
		go l.ls.Write(&LogContent{"WARN", l.msName, now.Format("2006-01-02 15:04:05"), toBsonM(content)})
	}
}

//Error ...
func (l *Logger) Error(content interface{}) {
	now := time.Now()
	fmt.Printf("[ERROR] %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
	if l.ls != nil {
		go l.ls.Write(&LogContent{"ERROR", l.msName, now.Format("2006-01-02 15:04:05"), toBsonM(content)})
	}
}

//Fatal ...
func (l *Logger) Fatal(content interface{}) {
	now := time.Now()
	fmt.Printf(
		"[FATAL] %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, content))
	if l.ls != nil {
		go l.ls.Write(&LogContent{"FATAL", l.msName, now.Format("2006-01-02 15:04:05"), toBsonM(content)})
	}
}

func (l *Logger) Read(where interface{}) ([]LogContent, error) {
	now := time.Now()
	fmt.Printf(
		"[READ] %s -> %s\n",
		now.Format("2006-01-02 15:04:05"),
		formatPrint(0, where))
	if l.ls != nil {
		return l.ls.Read(toBsonM(where))
	}
	return nil, errors.New("缺少日志对象")
}

func toBsonM(content interface{}) (bm bson.M) {
	bm = bson.M{}
	switch content.(type) {
	case int:
		bm["content"] = content
	case string:
		bm["content"] = content
	case error:
		bm["content"] = content.(error).Error()
	default:
		bc, _ := bson.Marshal(content)
		bson.Unmarshal(bc, &bm)
	}
	if len(bm) == 0 {
		fmt.Printf("Log nothing ! content type is [%T]\n", content)
	}
	return
}

func formatPrint(tier int, obj interface{}) (fmt string) {
	switch obj.(type) {
	case error:
		return "\"" + obj.(error).Error() + "\""
	case string:
		return "\"" + obj.(string) + "\""
	case int:
		return strconv.Itoa(obj.(int))
	case int64:
		return strconv.FormatInt(obj.(int64), 10)
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case map[string]interface{}:
		fmt = "\n"
		for k, v := range obj.(map[string]interface{}) {
			for i := 0; i < tier+1; i++ {
				fmt += "    "
			}
			fmt += k + " = " + formatPrint(tier+1, v) + "\n"
		}
	case []interface{}: //数组
		fmt = "\n"
		for i, v := range obj.([]interface{}) {
			for i := 0; i < tier+1; i++ {
				fmt += "    "
			}
			fmt += "[" + strconv.Itoa(i) + "] = " + formatPrint(tier+1, v) + "\n"
		}
	default:
		if obj == nil {
			fmt = "nil"
			return
		}
		fmt = "[" + reflect.TypeOf(obj).String() + " object]"
	}
	return
}
