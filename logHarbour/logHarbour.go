// logHarbour is the logging framework from remiges.tech
package logHarbour

import (
	"context"
	kafkaUtil "go-framework/kafkaUtil"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	//Logger types
	ACTIVITY_LOGGER   = "A"
	DEBUG_LOGGER      = "D"
	DATACHANGE_LOGGER = "C"
	//Log level priority
	PRI_DEBUG2   = "DEBUG2"
	PRI_DEBUG1   = "DEBUG1"
	PRI_DEBUG0   = "DEBUG0"
	PRI_INFO     = "INFO"
	PRI_WARN     = "WARN"
	PRI_ERROR    = "ERROR"
	PRI_CRITICAL = "CRIT"
	PRI_SECURITY = "SEC"
	// Log level Constants from a custom logging package.
	LevelDebug2   = slog.Level(-8)
	LevelDebug1   = slog.Level(-7)
	LevelDebug0   = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarning  = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelCritical = slog.Level(12)
	LevelSec      = slog.Level(16)
)

// application Indetifier
type appIdentifier struct {
	App    string `json:"app"`
	Module string `json:"module"`
	System string `json:"system"`
}

var identity appIdentifier

// checks if system is initialized
var isInitalized bool

// log level, used for printing of log entries
var programLevel = new(slog.LevelVar) // Info by default

// struct to manage 3 types of logger handles
type LogHandles struct {
	ActivityLogger   *slog.Logger
	DataChangeLogger *slog.Logger
	DebugLogger      *slog.Logger
}

// struct for managing data change objects
type dataChgObj struct {
	Field  string `json:"field"`
	OldVal string `json:"oldVal"`
	NewVal string `json:"newVal"`
}

// logHarbour Context
var ctx context.Context

// go runtime version
var goRuntime string

func init() {
	ctx = context.Background()
	buildInfo, _ := debug.ReadBuildInfo()
	goRuntime = buildInfo.GoVersion
	//Setting default log level of slog to lowest i.e. debug2 as printing of logs to slog will be managed by logHarbour
	programLevel.Set(LevelDebug2)
}

// initializes logHarbour with app, module and system names.
// Note that LogHarbour can only be initialized once.
func LogInit(appName, moduleName, systemName string) LogHandles {
	if len(appName) <= 0 || len(moduleName) <= 0 || len(systemName) <= 0 {
		log.Fatalf("Invalid app name[%v], module name[%v] or system name[%v]", appName, moduleName, systemName)
	}
	//will allow initialization only once
	if !isInitalized {
		identity = appIdentifier{appName, moduleName, systemName}
		isInitalized = true
		kafkaUtil.KafkaInit(appName + ":" + moduleName + ":" + systemName) //TODO: discuss this
	}
	return getLogger()
}

// func reads file name from config/env variable
func getRigelLogFileName() string {
	//TODO: read these parameters from config file
	filename := "logfile"
	filepath := "."
	fileSuffix := "yyyymmdd"
	fileExtn := ".txt"
	suffix := time.Now().Format(getLogFileFormat(fileSuffix))
	return filepath + "/" + filename + "_" + suffix + fileExtn
}

func getLogFileFormat(s string) string {
	switch s {
	case "yyyymmdd":
		return "20060102"
	case "ddmmyyyy":
		return "01022006"
	case "mmddyyyy":
		return "02012006"
	default:
		return "20060102"
	}
}

// manageAttributes is a function that manages the attributes of a slog.Attr object.
//
// It takes a slog.Attr object as a parameter and returns a slog.Attr object.
// If the Key of the parameter is equal to slog.TimeKey, it returns an empty slog.Attr object.
// If the Key of the parameter is equal to slog.LevelKey, it handles custom level values and returns the modified slog.Attr object.
// Otherwise, it returns the original slog.Attr object.
func manageAttributes(a slog.Attr) slog.Attr {
	/*if a.Key == slog.TimeKey {
		return slog.Attr{}
	}*/
	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Handle custom level values.
		level := a.Value.Any().(slog.Level)
		a.Value = getLogLevelString(level)
	}
	return a
}

// func returns string for log level passed
func getLogLevelString(level slog.Level) (levelString slog.Value) {
	switch {
	case level <= LevelDebug2:
		levelString = slog.StringValue(PRI_DEBUG2)
	case level <= LevelDebug1:
		levelString = slog.StringValue(PRI_DEBUG1)
	case level <= LevelDebug0:
		levelString = slog.StringValue(PRI_DEBUG0)
	case level <= slog.LevelInfo:
		levelString = slog.StringValue(PRI_INFO)
	case level <= LevelWarning:
		levelString = slog.StringValue(PRI_WARN)
	case level <= LevelError:
		levelString = slog.StringValue(PRI_ERROR)
	case level <= LevelCritical:
		levelString = slog.StringValue(PRI_CRITICAL)
	default:
		levelString = slog.StringValue(PRI_SECURITY)
	}
	return
}

var loggerSet LogHandles
var loggerSetInitialized = false

// func returns 3 logHandles for ActivityLog, DatachangeLog and DebugLog
func getLogger() LogHandles {
	if !loggerSetInitialized {
		//create log file
		logFile, err := os.OpenFile(getRigelLogFileName(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		//create Kafka Writer
		kafkaWriter := kafkaUtil.KafkaWriter{}
		//create multiwriter for logger to write to log file, stdout and kafka
		mw := io.MultiWriter(os.Stdout, logFile, kafkaWriter)

		lg := slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: programLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return manageAttributes(a)
		}})).With("app", identity.App).With("module", identity.Module).With("system", identity.System)

		//create child loggers
		loggerSet = LogHandles{
			ActivityLogger:   lg.With("handle", ACTIVITY_LOGGER),
			DataChangeLogger: lg.With("handle", DATACHANGE_LOGGER),
			//Here for debugLogger we can write pid and goruntime while creating handler.
			//However for calltrace() methods will be used while logging to capture correct call trace and source
			DebugLogger: lg.With("handle", DEBUG_LOGGER).With(slog.Int("pid", os.Getpid())).With(slog.String("runtime", goRuntime))}
	}
	return loggerSet
}

// func returns data change object
//
// field: field name whose value is changing | oldVal: old value before data change | newVal: new value after data change
func DataChg(field, oldVal, newVal string) dataChgObj {
	return dataChgObj{field, oldVal, newVal}
}

// getFrame returns the runtime.Frame at the specified index.
// The purpose of the function is to retrieve a specific frame from the call stack.
//
// It takes an integer parameter, skipFrames, which represents the number of frames to skip in the call stack.
// It returns a runtime.Frame.
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

// getCaller returns the name of the file and line number of the calling function.
// It uses the runtime.Caller function to retrieve information about the calling function two levels up the call stack.
//
// No parameters are required.
// It returns a string representing the file name and line number in the format "filename:linenumber".
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

func isTypeDataChg(t interface{}) bool {
	switch t.(type) {
	case dataChgObj:
		return true
	default:
		return false
	}
}

// func checks customMsgs of type Any to see if there is any data present in it. If no data is present, it returns an empty set of attribute.
func checkCustomMsg(lgger *slog.Logger, customMsgs ...any) slog.Attr {
	if len(customMsgs) > 0 {
		if lgger.Handler() == loggerSet.DataChangeLogger.Handler() {
			dco := make([]any, 0)
			for _, val := range customMsgs {
				if isTypeDataChg(val) {
					dco = append(dco, val)
				}
			}
			if len(dco) > 0 {
				return slog.Any("params", dco)
			} else {
				return slog.Attr{}
			}

		} else {
			otherAttribs := make([]any, 0)
			for i := 0; i < len(customMsgs); i++ {
				if !isTypeDataChg(customMsgs[i]) {
					otherAttribs = append(otherAttribs, customMsgs[i])
				}
			}
			if len(otherAttribs) > 0 {
				return slog.Group("params", otherAttribs...)
			} else {
				return slog.Attr{}
			}

		}
	} else {
		return slog.Attr{}
	}
}

// func writes log to specified source using slog
func LogWrite(lgger *slog.Logger, ll slog.Level, spanId, correlationId, when, who, remoteIp, op, whatClass, whatInstanceId string, status int, msg string, customMsgs ...any) {
	if !isInitalized {
		log.Fatalf("logHarbour not initialized. source[%s]. caller[%s]\n", getCallTrace(), getCaller())
	}

	if ll >= getRigelLogLevel() {
		//as a part of optimization, here we are using "slog.String()" func calls as recommended
		//TODO: check for field "when" cannot be a future time.
		//But currently we are taking it as a string, we'll be forced to take it as a time object. Need to discuss
		if ll <= LevelDebug0 {
			// In case of level of type Debug, additional information is passed to loggers
			lgger.LogAttrs(ctx, ll, msg, slog.String("source", getCaller()), slog.String("callTrace", getCallTrace()), slog.String("spanId", spanId), slog.String("correlationId", correlationId), slog.String("when", when), slog.String("who", who), slog.String("remoteIp", remoteIp), slog.String("op", op), slog.String("whatClass", whatClass), slog.String("whatInstanceId", whatInstanceId), slog.Int("status", status), checkCustomMsg(lgger, customMsgs...))
		} else {
			lgger.LogAttrs(ctx, ll, msg, slog.String("spanId", spanId), slog.String("correlationId", correlationId), slog.String("when", when), slog.String("who", who), slog.String("remoteIp", remoteIp), slog.String("op", op), slog.String("whatClass", whatClass), slog.String("whatInstanceId", whatInstanceId), slog.Int("status", status), checkCustomMsg(lgger, customMsgs...))
		}
	}
}

// TODO: STUB func to get log level from Rigel
func getRigelLogLevel() slog.Level {
	//TODO : change this with corresponding log level call from Rigel
	//Else, by default, it returns the log level of slog
	return programLevel.Level()
}
