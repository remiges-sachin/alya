package logHarbour

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

var GlobalLogLevel LogLevel
var goRuntime string

func init() {

	GlobalLogLevel = Inf
	buildInfo, _ := debug.ReadBuildInfo()
	goRuntime = buildInfo.GoVersion
	//foo()
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// MyCaller returns the caller of the function that called it :)
func myCallTrace() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2).Function
}

func getCaller() string {
	_, file, line, _ := runtime.Caller(2)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}

type LogLevel int

const (
	Inf LogLevel = 1
	Err LogLevel = 2
	Dbg LogLevel = 0
	Trc LogLevel = -1
)

func shouldPrint(logLevel LogLevel) bool {
	if isInitalized {
		return logLevel >= GlobalLogLevel
	}
	return false
}

func getLogString(ll LogLevel) string {
	switch ll {
	case Inf:
		return "Info"
	case Err:
		return "Error"
	case Dbg:
		return "Debug"
	case Trc:
		return "Trace"
	default:
		return "Info"
	}
}

type logStruct struct {
	LogLevelInf string `json:"pri"` //could be "level" as per convention
	App         string `json:"app"`
	Module      string `json:"module"`
	System      string `json:"system"`
	Caller      string `json:"caller,omitempty"`
	CallTrace   string `json:"callTrace,omitempty"`
	Pid         int    `json:"pid,omitempty"` //assuming pid for golang process can never be 0
	GoVersion   string `json:"goVersion,omitempty"`
	When        string `json:"when"`
	Who         string `json:"who"`
	RemoteIp    string `json:"remoteIp"`
	Op          string `json:"op"`
	What        string `json:"what"`
	Status      int    `json:"status"`
	Msg         string `json:"msg"`
	Params      any    `json:"params,omitempty"`
}

type DataChg struct {
	Field  string `json:"field"`
	OldVal string `json:"oldVal"`
	NewVal string `json:"newVal"`
}

func GetDataChg(field, oldVal, newVal string) DataChg {
	return DataChg{field, oldVal, newVal}
}

func checkAny(customMsgs ...any) any {
	if len(customMsgs) > 0 {
		return customMsgs
	} else {
		var emptyInterface interface{}
		return emptyInterface
	}
}

func GetKV(key string, val string) map[string]string {
	a := map[string]string{key: val}
	return a
}

// writes Log.
// TODO : check for possible sync issues
func LogWrite(ll LogLevel, when, who, remoteIp, op, what string, status int, msg string, customMsgs ...any) {
	if shouldPrint(ll) {
		if when == "" {
			when = time.Now().UTC().String()
		}
		var ls logStruct
		switch ll {
		case Inf, Err:
			ls = logStruct{getLogString(ll), app, module, system, "", "", 0, "", when, who, remoteIp, op, what, status, msg, checkAny(customMsgs...)}
		case Dbg, Trc:
			ls = logStruct{getLogString(ll), app, module, system, getCaller(), myCallTrace(), os.Getpid(), goRuntime, when, who, remoteIp, op, what, status, msg, checkAny(customMsgs...)}
		}
		sendLog(ls)
	}
}

// this func decides where to send logs
// TODO: create io function to write to output(file/port etc)
func sendLog(ls logStruct) {
	json_data, err := json.Marshal(ls)
	if err != nil {
		log.Fatal("err:", err)
	}
	fmt.Println(string(json_data))
}

var app string
var module string
var system string
var isInitalized bool

func LogInit(appName, moduleName, systemName string) {
	if !isInitalized {
		app = appName
		module = moduleName
		system = systemName
		isInitalized = true
	}
}
