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
	logHarbour.GlobalLogLevel = logHarbour.Dbg
	logHarbour.LogWrite(logHarbour.Inf, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger info firstFunc", logHarbour.GetDataChg("qty", "100", "200"), logHarbour.GetDataChg("qty", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger error firstFunc", logHarbour.GetKV("reqIdCustom", "123123123"), logHarbour.GetKV("otherField", "otherfieldvalue"), logHarbour.GetDataChg("qty", "1", "2"), logHarbour.GetDataChg("amt", "100", "200"))
	fmt.Println()
	logHarbour.LogWrite(logHarbour.Dbg, time.Now().UTC().String(), "bhavya", "127.0.0.1", "newLog", "NA", 1, "This is an activity logger Debug firstFunc", logHarbour.GetDataChg("qty", "100", "200"))
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
{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-07 03:34:29.122988013 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info firstFunc","params":[{"field":"qty","oldVal":"100","newVal":"200"},{"field":"qty","oldVal":"100","newVal":"200"}]}

{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:20","callTrace":"main.firstFunc","pid":242396,"goVersion":"go1.20.7","when":"2023-09-07 03:34:29.123061786 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error firstFunc","params":[{"reqIdCustom":"123123123"},{"otherField":"otherfieldvalue"},{"field":"qty","oldVal":"1","newVal":"2"},{"field":"amt","oldVal":"100","newVal":"200"}]}

{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:22","callTrace":"main.firstFunc","pid":242396,"goVersion":"go1.20.7","when":"2023-09-07 03:34:29.123084494 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger Debug firstFunc","params":[{"field":"qty","oldVal":"100","newVal":"200"},{"field":"qty","oldVal":"100","newVal":"200"}]}

{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-07 03:34:29.12309635 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info Second Func"}
{"pri":"Error","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-07 03:34:29.123099658 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error Second Func"}
{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:30","callTrace":"main.secondFunc","pid":242396,"goVersion":"go1.20.7","when":"2023-09-07 03:34:29.123102405 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger debug Second Func"}
{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:32","callTrace":"main.secondFunc","pid":242396,"goVersion":"go1.20.7","when":"2023-09-07 03:34:29.123110203 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger debug Second Func"}

{"pri":"Info","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-07 03:34:29.123118103 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger info Third Func"}
{"pri":"Error","app":"sampleApp","module":"moduleName","system":"systemName","when":"2023-09-07 03:34:29.123120918 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger error Third Func"}
{"pri":"Debug","app":"sampleApp","module":"moduleName","system":"systemName","caller":"sampleMain.go:41","callTrace":"main.thirdFunc","pid":242396,"goVersion":"go1.20.7","when":"2023-09-07 03:34:29.123123754 +0000 UTC","who":"bhavya","remoteIp":"127.0.0.1","op":"newLog","what":"NA","status":1,"msg":"This is an activity logger debug Third Func"}
*/
