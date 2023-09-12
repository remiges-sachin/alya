package main

import (
	"fmt"
	logHarbour "go-framework/logHarbour"
	"time"
)

var timeFormat string = "2006-01-02T15:04:05Z"

func main() {
	loggers := logHarbour.LogInit("app1", "module1", "system1")

	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelInfo, "spanid1", "correlationid1", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info message", "somekey", "somevalue")
	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelError, "spanid2", "correlationid2", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error message")

	fmt.Println("----")
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelInfo, "spanid3", "correlationid3", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an DataChange logger info message", "somekey", "somevalue", logHarbour.DataChg{Field: "amt", OldVal: "100", NewVal: "200"})
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelError, "spanid4", "correlationid4", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an DataChange logger error message", logHarbour.DataChg{Field: "amt", OldVal: "100", NewVal: "200"})
	fmt.Println("----")
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug, "spanid5", "correlationid5", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an DEBUG logger message 1", logHarbour.DataChg{Field: "amt", OldVal: "100", NewVal: "200"}, logHarbour.GetDataChg("qty", "1", "2"))
	logHarbour.ChangeGlobalLogLevel(logHarbour.LevelDebug2)
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug1, "spanid6", "correlationid6", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an DEBUG logger message 2", logHarbour.DataChg{Field: "amt", OldVal: "100", NewVal: "200"}, logHarbour.GetDataChg("qty", "1", "2"))
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug2, "spanid7", "correlationid7", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an DEBUG logger message 3")
	fmt.Println("---------------------------")
	fmt.Println("---------------------------")
	firstFunc(loggers)
	//main2()
}

func firstFunc(loggers logHarbour.LogHandles) {
	//logHarbour.GlobalLogLevel = logHarbour.Dbg
	logHarbour.LogWrite(loggers.ActivityLogger, logHarbour.LevelInfo, "spanid12", "correlationid12", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info firstFunc", logHarbour.GetDataChg("qty", "100", "200"), logHarbour.GetDataChg("qty", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(loggers.DataChangeLogger, logHarbour.LevelError, "spanid13", "correlationid13", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error firstFunc", "reqIdCustom", "123123123", "otherField", "otherfieldvalue", logHarbour.GetDataChg("qty", "1", "2"), logHarbour.GetDataChg("amt", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(loggers.DebugLogger, logHarbour.LevelDebug1, "spanid14", "correlationid14", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger Debug firstFunc", logHarbour.GetDataChg("qty", "100", "200"))
	fmt.Println()
	//secondFunc()
}

/*
func secondFunc() {
	logHarbour.LogWrite(logHarbour.Inf, "spanid15", "correlationid15", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info Second Func")
	logHarbour.LogWrite(logHarbour.Err, "spanid16", "correlationid16", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error Second Func")
	logHarbour.LogWrite(logHarbour.Dbg, "spanid17", "correlationid17", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Second Func")
	logHarbour.SetGlobalLogLevel(logHarbour.Dbg)
	logHarbour.LogWrite(logHarbour.Dbg, "spanid18", "correlationid18", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Second Func")
	fmt.Println()
	thirdFunc()
}

func thirdFunc() {
	logHarbour.SetGlobalLogLevel(logHarbour.Err)
	logHarbour.LogWrite(logHarbour.Inf, "spanid19", "correlationid19", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info Third Func")
	logHarbour.LogWrite(logHarbour.Err, "spanid20", "correlationid20", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error Third Func")
	logHarbour.LogWrite(logHarbour.Dbg, "spanid21", "correlationid21", time.Now().Format(timeFormat), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Third Func")
}
*/
