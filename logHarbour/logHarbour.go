package logHarbour

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync/atomic"
	"time"
)

// log level to print. Will be used while printing of log messages to a defined output
var globalLogLevel int64

func SetGlobalLogLevel(ll LogLevel) {
	atomic.StoreInt64(&globalLogLevel, int64(ll))
}

// go runtime version
var goRuntime string

type logIdentity struct {
	App    string `json:"app"`
	Module string `json:"module"`
	System string `json:"system"`
}

var identity logIdentity

// checks if system is initialized
var isInitalized bool

// initializes logHarbour with app, module and system names
func LogInit(appName, moduleName, systemName string) {
	//will allow initialization only once
	if !isInitalized {
		identity = logIdentity{appName, moduleName, systemName}
		isInitalized = true
	}
}

func init() {
	atomic.StoreInt64(&globalLogLevel, int64(Inf))
	buildInfo, _ := debug.ReadBuildInfo()
	goRuntime = buildInfo.GoVersion
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

// getCallTrace returns the caller of the function that called it :)
func getCallTrace() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2).Function
}

// getCaller() returns the file and line no from where the func is called
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

// struct to write debug info to log messages
type debugInfo struct {
	Caller    string `json:"caller,omitempty"`
	CallTrace string `json:"callTrace,omitempty"`
	Pid       int    `json:"pid,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
}

// log level
type LogLevel int64

// log level constants
const (
	Inf LogLevel = 1
	Err LogLevel = 2
	Dbg LogLevel = 0
	Trc LogLevel = -1
)

// func checks log level and returns a bool confirming whether the statement should actually be logged or not.
func shouldPrint(logLevel LogLevel) bool {
	if isInitalized {
		return int64(logLevel) >= globalLogLevel
	}
	return false
}

// func returns a string based on log level
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

// struct for logging messages
type logMsg struct {
	LogLevelInf string `json:"pri"` //could be "level" as per convention
	logIdentity
	debugInfo
	When     string `json:"when"`
	Who      string `json:"who"`
	RemoteIp string `json:"remoteIp"`
	Op       string `json:"op"`
	What     string `json:"what"`
	Status   int    `json:"status"`
	Msg      string `json:"msg"`
	Params   any    `json:"params,omitempty"`
}

// struct for managing data change objects
type dataChg struct {
	Field  string `json:"field"`
	OldVal string `json:"oldVal"`
	NewVal string `json:"newVal"`
}

// func returns data change object
func GetDataChg(field, oldVal, newVal string) dataChg {
	return dataChg{field, oldVal, newVal}
}

// creates an empty interface if nothing is passed to the variadic function. This is imp so that the json parser skips the "params" tag while creating the message
func checkAny(customMsgs ...any) any {
	if len(customMsgs) > 0 {
		return customMsgs
	} else {
		var emptyInterface interface{}
		return emptyInterface
	}
}

// func returns a key:value map for printing a json
func GetKV(key string, val string) map[string]string {
	a := map[string]string{key: val}
	return a
}

// writes Log.
// TODO : check for possible sync issues
func LogWrite(ll LogLevel, when, who, remoteIp, op, what string, status int, msg string, customMsgs ...any) {
	if !isInitalized {
		log.Fatalf("logHarbour not initialized. callTrace[%s]. caller[%s]\n", getCallTrace(), getCaller())
	}
	if shouldPrint(ll) {
		//if no time is passed, it'll print time of printing the log
		if when == "" {
			when = time.Now().Format("2006-01-02T15:04:05Z")
		}
		var lm logMsg
		switch ll {
		case Inf, Err:
			lm = logMsg{getLogString(ll), identity, debugInfo{}, when, who, remoteIp, op, what, status, msg, checkAny(customMsgs...)}
		case Dbg, Trc:
			lm = logMsg{getLogString(ll), identity, debugInfo{getCaller(), getCallTrace(), os.Getpid(), goRuntime}, when, who, remoteIp, op, what, status, msg, checkAny(customMsgs...)}
		}
		sendLog(lm)
	}
}

// this func decides where to send logs
// TODO: create io function to write to output(file/port etc)
// TODO: check for possible sync issues
func sendLog(ls logMsg) {
	json_data, err := json.Marshal(ls)
	if err != nil {
		log.Fatalf("Error while encoding to json msg[%v]: \nError: %v \n", ls, err)
	}
	fmt.Println(string(json_data))
}
