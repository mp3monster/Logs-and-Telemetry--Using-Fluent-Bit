package main

import (
	"C"
	"log"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"

	"strings"
)
import "errors"

const PluginName = "out_gdb"

func releaseResources() error {
	clearEnvParams(PluginName)
	return nil
}

// retrieves the configuration values etc
func getParams(plugin unsafe.Pointer) (*MyParams, error) {
	if plugin == nil {
		return nil, errors.New("Plugin not provided")
	}
	params := MyParams{}
	params.PluginName = PluginName
	params.Host = output.FLBPluginConfigKey(plugin, Plugin_Host)
	//TODO: implement the full logic

	//params.DeleteAfterQuery = strings.Contains(strings.ToLower(output.FLBPluginConfigKey(plugin, Plugin_Delete)), "true")

	return &params, nil
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	log.Printf("[%s] Register called", PluginName)

	registered := output.FLBPluginRegister(def, PluginName, "Your description here")
	log.Printf("[%s] Registration result =%v\n", PluginName, registered == output.FLB_OK)
	return output.FLB_OK
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	params, err := getParams(plugin)
	if err != nil {
		log.Printf("[%s]Init for %s errored with: %s", PluginName, PluginName, err)
		return output.FLB_ERROR
	}
	if strings.Contains(strings.ToLower(output.FLBPluginConfigKey(plugin, "Log_Level")), "debug") {
		log.Printf("[%s] configured with %v\n", params.PluginName, SprintfParams(params, PluginName))
	}

	//TODO: validateMyParams needs implementing unless you leverage the common element of the complete code
	validateErr := validateMyParams(params)
	if validateErr != nil {
		log.Printf("[%s] Configuration error -%s\n", params.PluginName, validateErr)
		return output.FLB_ERROR
	}

	//paramsToEnv(params, PluginName)
	paramsJSON := paramsToJSON(params)
	log.Printf("Adding to context params==>%s", paramsJSON)
	output.FLBPluginSetContext(plugin, &paramsJSON)

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.Printf("[%s] Flush called for unknown instance\n", PluginName)
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Type assert context back into the original type for the Go variable
	var params *MyParams
	myContext := output.FLBPluginGetContext(ctx)
	if myContext != nil {
		strContext := myContext.(*string)
		params = JSONToParams(*strContext, PluginName)
		log.Printf("[%s] Flush called with context", params.PluginName)
	} else {
		log.Printf("[%s] Flush called with no context\n", PluginName)
		if params == nil {
			log.Printf("[%s] FLBPluginFlushCtx no params\n", PluginName)
			return output.FLB_ERROR
		}
	}

	// without our params we cant do anthing - bail
	if params == nil {
		log.Printf("[%s] FLBPluginFlushCtx - invoked without params", PluginName)
		return output.FLB_ERROR
	}

	dec := output.NewDecoder(data, int(length))

	log.Printf("[%s] FLBPluginFlushCtx about to process:\n", PluginName)

	count := 0
	for { // for as long as there is a data value to insert
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		var timestamp time.Time
		switch t := ts.(type) {
		case output.FLBTime:
			timestamp = ts.(output.FLBTime).Time
		case uint64:
			timestamp = time.Unix(int64(t), 0)
		default:
			log.Println("[%s]time provided invalid, defaulting to now - received %T (%v). Timestamp is %d", PluginName, t, t, timestamp)
			timestamp = time.Now()
		}

		// Print record keys and values
		//log.Printf("[%s] record received:%v", PluginName, record)
		count++
		//TODO: implement logic to store data value
		if insertErr != nil {
			log.Printf("[%s] Error during insert, returning fail\n%v", PluginName, insertErr)
			return output.FLB_ERROR
		}
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	log.Print("[out_gdb] Exit called for unknown instance")
	return output.FLB_OK
}

//export FLBPluginExitCtx
func FLBPluginExitCtx(ctx unsafe.Pointer) int {
	// Type assert context back into the original type for the Go variable
	context := output.FLBPluginGetContext(ctx)
	if context != nil {
		log.Printf("[%s] Exit called with context: %v", PluginName, context)
	} else {
		// if we don't have a context to work with we have to pass our configs around using the
		log.Printf("[%s] Exit called NIL context", PluginName)
		releaseResources()
	}
	return output.FLB_OK
}

//export FLBPluginUnregister
func FLBPluginUnregister(def unsafe.Pointer) {
	log.Print("[out_gdb] Unregister called")
	output.FLBPluginUnregister(def)
}
