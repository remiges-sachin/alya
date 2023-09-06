package main

import (
	"fmt"
	"go-framework/logHarbour"
	"time"
)

func main() {
	logHarbour.LogInit("sampleApp", "moduleName", "systemName")
	logHarbour.LogWrite(logHarbour.Trc, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info message")
	fmt.Println()
	firstFunc()
}

func firstFunc() {
	logHarbour.LogWrite(logHarbour.Inf, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info firstFunc")
	logHarbour.LogWrite(logHarbour.Err, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error firstFunc")
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger Debug firstFunc")
	fmt.Println()
	secondFunc()
}

func secondFunc() {
	logHarbour.LogWrite(logHarbour.Inf, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info Second Func")
	logHarbour.LogWrite(logHarbour.Err, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error Second Func")
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Second Func")
	logHarbour.GlobalLogLevel = logHarbour.Dbg
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Second Func")
	fmt.Println()
	thirdFunc()
}

func thirdFunc() {
	logHarbour.GlobalLogLevel = logHarbour.Dbg
	logHarbour.LogWrite(logHarbour.Inf, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info Third Func")
	logHarbour.LogWrite(logHarbour.Err, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error Third Func")
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger debug Third Func")
}

/* output:
{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.646216581 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info firstFunc"}
{"pri":"Error","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.646276406 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error firstFunc"}

{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.64628321 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info Second Func"}
{"pri":"Error","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.646287205 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error Second Func"}
{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:29","callTrace":"main.secondFunc","when":"2023-09-06 16:05:26.646291583 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger debug Second Func"}

{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.646314479 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info Third Func"}
{"pri":"Error","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-06 16:05:26.646318407 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error Third Func"}
{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:38","callTrace":"main.thirdFunc","when":"2023-09-06 16:05:26.646322172 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger debug Third Func"}

*/
