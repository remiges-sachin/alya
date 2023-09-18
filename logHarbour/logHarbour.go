package logHarbour

import (
	"context"
	"fmt"
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
	PRI_DEBUG    = "DEBUG"
	PRI_INFO     = "INFO"
	PRI_WARN     = "WARN"
	PRI_ERROR    = "ERROR"
	PRI_CRITICAL = "CRIT"
	PRI_SECURITY = "SEC"
	// Log level Constants from a custom logging package.
	LevelDebug2   = slog.Level(-8)
	LevelDebug1   = slog.Level(-7)
	LevelDebug    = slog.LevelDebug
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

// to change global logging level.
//
// TODO : Identify way to change the log level at runtime
func ChangeGlobalLogLevel(level slog.Level) {
	programLevel.Set(level)
	fmt.Printf("[%s]:GLOBAL log level set to: [%s]\n", time.Now().UTC(), getLogLevelString(level))
}

// struct to manage 3 types of logger handles
type LogHandles struct {
	ActivityLogger   *slog.Logger
	DataChangeLogger *slog.Logger
	DebugLogger      *slog.Logger
}

// struct for managing data change objects
type DataChg struct {
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
}

// initializes logHarbour with app, module and system names.
// Note that LogHarbour can only be initialized once.
func LogInit(appName, moduleName, systemName string) LogHandles {
	//will allow initialization only once
	if !isInitalized {
		identity = appIdentifier{appName, moduleName, systemName}
		isInitalized = true
	}
	return getLogger()
}

// manageAttributes is a function that manages the attributes of a slog.Attr object.
//
// It takes a slog.Attr object as a parameter and returns a slog.Attr object.
// If the Key of the parameter is equal to slog.TimeKey, it returns an empty slog.Attr object.
// If the Key of the parameter is equal to slog.LevelKey, it handles custom level values and returns the modified slog.Attr object.
// Otherwise, it returns the original slog.Attr object.
func manageAttributes(a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
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
	case level <= LevelDebug:
		levelString = slog.StringValue(PRI_DEBUG)
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

// func returns 3 logHandles for ActivityLog, DatachangeLog and DebugLog
func getLogger() LogHandles {
	lg := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		return manageAttributes(a)
	}})).With("app", identity.App).With("module", identity.Module).With("system", identity.System)

	return LogHandles{lg.With("handle", ACTIVITY_LOGGER),
		lg.With("handle", DATACHANGE_LOGGER),
		lg.With("handle", DEBUG_LOGGER).With("pid", os.Getpid()).With("goVersion", goRuntime)}
}

// func returns data change object
func GetDataChg(field, oldVal, newVal string) DataChg {
	return DataChg{field, oldVal, newVal}
}

// LogValue implements slog.LogValuer.
// It returns a group containing the fields of
// the Name, so that they appear together in the log output.
func (dc DataChg) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("field", dc.Field),
		slog.String("oldVal", dc.OldVal),
		slog.String("newVal", dc.NewVal))
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

// func checks customMsgs of type Any to see if there is any data present in it. If no data is present, it returns an empty set of attribute.
func checkCustomMsg(customMsgs ...any) slog.Attr {
	if len(customMsgs) > 0 {
		return slog.Any("params", customMsgs)
	} else {
		return slog.Attr{}
	}
}

// func writes log to specified source using slog
func LogWrite(lgger *slog.Logger, ll slog.Level, spanId, correlationId, when, who, remoteIp, op, what string, status int, msg string, customMsgs ...any) {
	if !isInitalized {
		log.Fatalf("logHarbour not initialized. source[%s]. caller[%s]\n", getCallTrace(), getCaller())
	}

	if ll <= LevelDebug {
		// In case of level of type Debug, additional information is passed to loggers
		lgger.Log(ctx, ll, msg, "source", getCaller(), "callTrace", getCallTrace(), "spanId", spanId, "correlationId", correlationId, "when", when, "who", who, "remoteIp", remoteIp, "op", op, "what", what, "status", status, checkCustomMsg(customMsgs...))
	} else {
		lgger.Log(ctx, ll, msg, "spanId", spanId, "correlationId", correlationId, "when", when, "who", who, "remoteIp", remoteIp, "op", op, "what", what, "status", status, checkCustomMsg(customMsgs...))
	}
}
