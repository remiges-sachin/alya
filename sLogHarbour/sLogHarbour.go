package sLogHarbour

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Exported constants from a custom logging package.
const (
	LevelDebug2   = slog.Level(-8)
	LevelDebug1   = slog.Level(-7)
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarning  = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelCritical = slog.Level(12)
	LevelSec      = slog.Level(16)
)

var LogHandle *slog.Logger

type LogHandles struct {
	ActivityLogger   *slog.Logger
	DataChangeLogger *slog.Logger
	DebugLogger      *slog.Logger
}

const (
	ACTIVITY_LOGGER   = "A"
	DEBUG_LOGGER      = "D"
	DATACHANGE_LOGGER = "C"
	PRI_DEBUG2        = "DEBUG2"
	PRI_DEBUG1        = "DEBUG1"
	PRI_DEBUG         = "DEBUG"
	PRI_INFO          = "INFO"
	PRI_WARN          = "WARN"
	PRI_ERROR         = "ERROR"
	PRI_CRITICAL      = "CRIT"
	PRI_SECURITY      = "SEC"
)

type appIdentifier struct {
	App    string `json:"app"`
	Module string `json:"module"`
	System string `json:"system"`
}

var identity appIdentifier

// checks if system is initialized
var isInitalized bool

// initializes logHarbour with app, module and system names
func LogInit(appName, moduleName, systemName string) LogHandles {
	//will allow initialization only once
	if !isInitalized {
		identity = appIdentifier{appName, moduleName, systemName}
		isInitalized = true
	}
	return getLogger()
}

func manageAttributes(a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Rename the level key from "level" to "sev".
		a.Key = "pri"

		// Handle custom level values.
		level := a.Value.Any().(slog.Level)
		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		a.Value = getLogLevelString(level)
	}
	return a
}

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

func getLogger() LogHandles {
	return LogHandles{slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		return manageAttributes(a)
	}})).With("handle", ACTIVITY_LOGGER),
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return manageAttributes(a)
		}})).With("handle", DATACHANGE_LOGGER),
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return manageAttributes(a)
		}})).With("handle", DEBUG_LOGGER).With("pid", os.Getpid()).With("source", getCaller(1)).With("callTrace", getCallTrace(1))}
}

var programLevel = new(slog.LevelVar) // Info by default

// to change global zerolog logging level.
func ChangeGlobalLogLevel(level slog.Level) {
	programLevel.Set(level)
	fmt.Printf("[%s]:GLOBAL log level set to: [%s]\n", time.Now().UTC(), getLogLevelString(level))
}

// struct for managing data change objects
type DataChg struct {
	Field  string `json:"field"`
	OldVal string `json:"oldVal"`
	NewVal string `json:"newVal"`
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
func getCallTrace(i int) string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2 + i).Function
}

func getCaller(i int) string {
	_, file, line, _ := runtime.Caller(2 + i)
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
		log.Fatalf("logHarbour not initialized. source[%s]. caller[%s]\n", getCallTrace(0), getCaller(0))
	}
	ctx := context.Background()
	lgger.Log(ctx, ll, msg, "spanId", spanId, "correlationId", correlationId, "when", when, "who", who, "remoteIp", remoteIp, "op", op, "what", what, "status", status, checkCustomMsg(customMsgs...))
}
