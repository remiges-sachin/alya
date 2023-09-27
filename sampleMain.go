package main

import (
	"fmt"
	logHarbour "go-framework/logHarbour"
	"time"
)

var timeFormat string = "2006-01-02T15:04:05Z"

func main() {
	loggers := logHarbour.LogInit("app1", "module1", "system1")

	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelInfo, "spanid1", "correlationid1", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id1", 1, "This is an activity logger info message", "somekey", "somevalue", "key2", "value2")
	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelError, "spanid2", "correlationid2", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id2", 1, "This is an activity logger error message", "somekey", "somevalue", logHarbour.DataChg("amt", "100", "200"), logHarbour.DataChg("qty", "1", "2"))
	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelError, "spanid3", "correlationid3", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id3", 1, "This is an activity logger error message")

	fmt.Println("----")
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelInfo, "spanid4", "correlationid4", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id4", 1, "This is an DataChange logger info message", "somekey", "somevalue")
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelInfo, "spanid5", "correlationid5", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id5", 1, "This is an DataChange logger info message", "somekey", "somevalue", logHarbour.DataChg("amt", "100", "200"), logHarbour.DataChg("qty", "1", "2"))
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelError, "spanid6", "correlationid6", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id6", 1, "This is an DataChange logger error message", logHarbour.DataChg("amt", "100", "200"), logHarbour.DataChg("qty", "1", "2"))
	fmt.Println("----")
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug0, "spanid7", "correlationid7", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id7", 1, "This is an DEBUG0 logger message 1", "key3", "value3", logHarbour.DataChg("amt", "100", "200"), logHarbour.DataChg("qty", "1", "2"))
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug1, "spanid8", "correlationid8", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id8", 1, "This is an DEBUG1 logger message 2", "key3", "value3", logHarbour.DataChg("qty", "1", "2"))
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug2, "spanid9", "correlationid9", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id9", 1, "This is an DEBUG2 logger message 3")
	fmt.Println("---------------------------")
	fmt.Println("---------------------------")
	firstFunc(loggers)
}

func firstFunc(loggers logHarbour.LogHandles) {
	//logHarbour.GlobalLogLevel = logHarbour.Dbg
	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelInfo, "spanid12", "correlationid12", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id12", 1, "Activity logger info Message", logHarbour.DataChg("qty", "1", "2"), logHarbour.DataChg("amt", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelError, "spanid13", "correlationid13", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id13", 1, "Datachangelogger error Message", logHarbour.DataChg("qty", "1", "2"), logHarbour.DataChg("amt", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug1, "spanid14", "correlationid14", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "valueBeingUpdated", "id14", 1, "debug DEBUG1 Message", logHarbour.DataChg("qty", "1", "2"), "key1", "value1")
}
