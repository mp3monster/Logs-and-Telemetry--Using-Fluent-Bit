package main

import (
	"C"
	"log"
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

// built the params struct by retrieving from the plugin congiuration the values, including
// translating data to the correct types
func getParams(plugin unsafe.Pointer) (*SqlParams, error) {
	if plugin == nil {
		return nil, errors.New("Plugin not provided")
	}
	params := SqlParams{}
	params.PluginName = PluginName

	params.InstanceName = output.FLBPluginConfigKey(plugin, Plugin_InstanceId)
	params.Host = output.FLBPluginConfigKey(plugin, Plugin_Host)
	params.Port = output.FLBPluginConfigKey(plugin, Plugin_Port)
	params.User = output.FLBPluginConfigKey(plugin, Plugin_User)
	params.Password = output.FLBPluginConfigKey(plugin, Plugin_Password)
	params.SequencerCol = output.FLBPluginConfigKey(plugin, Plugin_Ordering)
	params.TableName = output.FLBPluginConfigKey(plugin, Plugin_TableName)
	params.DBName = output.FLBPluginConfigKey(plugin, Plugin_DBName)
	params.DBType = output.FLBPluginConfigKey(plugin, Plugin_Type)
	params.PK = output.FLBPluginConfigKey(plugin, Plugin_PK)
	params.ColsCSV = output.FLBPluginConfigKey(plugin, Plugin_ColsCSV)
	params.WhereExpr = output.FLBPluginConfigKey(plugin, Plugin_WhereExpr)

	params.DeleteAfterQuery = strings.Contains(strings.ToLower(output.FLBPluginConfigKey(plugin, Plugin_Delete)), "true")

	return &params, nil
}

// The function that provides the details of the plugin to Fluent Bit to allow the association
// of the plugin name to the config file, and allow the CLI help to show the plugin role.
//
//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	log.Printf("[%s] Register called", PluginName)

	registered := output.FLBPluginRegister(def, PluginName, "Go plugin for writing content from a database")
	log.Printf("[%s] Registration result =%v\n", PluginName, registered == output.FLB_OK)
	return output.FLB_OK
}

// Called after the the registration, this callback is triggered so that the configuration values can be retrieved
// and checked - if we do not have all the necessary values, or they're incorrectly formed we need to report an
// error back to the core of Fluent Bit.
// For this plugin we check the parameters and prove we can ping the database
//
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

	validateErr := validateSqlParams(params)
	if validateErr != nil {
		log.Printf("[%s] %s Configuration error -%s\n", params.PluginName, params.InstanceName, validateErr)
		return output.FLB_ERROR
	}

	connectWorks := testConnectionOk(params)
	log.Printf("[%s] %s Init connection test successful %t\n", params.PluginName, params.InstanceName, connectWorks)
	if !connectWorks {
		return output.FLB_ERROR
	}

	//paramsToEnv(params, PluginName)
	paramsJSON := paramsToJSON(params)
	log.Printf("Adding to context params==>%s", paramsJSON)
	output.FLBPluginSetContext(plugin, &paramsJSON)

	return output.FLB_OK
}

// The plugin has flush with and without contexts. We want to use the conte
// , so just return ok. But also allow us to see it being invoked with a log message
//
//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.Printf("[%s] Flush called for unknown instance\n", PluginName)
	return output.FLB_OK
}

// This is the context based flush, This callback is responsible for the outyput to the destination of
// the supplied data. To ensure we're sending the data in the correct direction we make use of
// the context data we have stored.
//
//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Type assert context back into the original type for the Go variable
	//var params *SqlParams
	params := NewSqlParams()
	myContext := output.FLBPluginGetContext(ctx)
	if myContext != nil {
		strContext := myContext.(*string)
		params = JSONToParams(*strContext, PluginName)
		//log.Printf("[%s] Flush called for context: %s", params.PluginName, *strContext)
		log.Printf("[%s]%s Flush called with context", params.PluginName, params.InstanceName)
	} else {
		log.Printf("[%s] Flush called with no context\n", PluginName)
		params = envToParams(PluginName)
		if params == nil {
			log.Printf("[%s] FLBPluginFlushCtx no params\n", PluginName)
			return output.FLB_ERROR
		}
	}

	// without our params we cant do anthing - bail
	if params == nil {
		log.Printf("[%s] FLBPluginFlushCtx - invoked without params", params.PluginName)
		return output.FLB_ERROR
	}

	dec := output.NewDecoder(data, int(length))

	count := 0
	for { // for as long as there is a data value to insert
		ret, ts, record := output.GetRecord(dec)

		if ret != 0 {
			break
		}
		log.Printf("[%s]%s FLBPluginFlushCtx about to process:%v with timestamp %v", PluginName, params.InstanceName, ret, ts)

		/*		var timestamp time.Time
				switch t := ts.(type) {
				case output.FLBTime:
					timestamp = ts.(output.FLBTime).Time
				case uint64:
					timestamp = time.Unix(int64(t), 0)
				case string:
					timestamp = ts.(string)
				default:
					log.Println("[%s]%s time provided invalid, defaulting to now - received (%T). Timestamp is %v", params.PluginName, params.InstanceName, ts, ts)
					timestamp = time.Now()
				}
		*/
		// Print record keys and values
		//log.Printf("[%s] record received:%v", PluginName, record)
		count++
		insertErr := execInsert(params, record)
		if insertErr != nil {
			log.Printf("[%s]%s Error during insert, returning fail\n%v", params.PluginName, params.InstanceName, insertErr)
			return output.FLB_ERROR
		}
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	log.Print("[%s] Exit called for unknown instance", PluginName)
	return output.FLB_OK
}

// When the plugin needs to be stopped then this method is called. For example a hot reload
//
//export FLBPluginExitCtx
func FLBPluginExitCtx(ctx unsafe.Pointer) int {
	// Type assert context back into the original type for the Go variable
	params := NewSqlParams()
	context := output.FLBPluginGetContext(ctx)
	if context != nil {
		strContext := context.(*string)
		params = JSONToParams(*strContext, PluginName)
		//log.Printf("[%s] Flush called for context: %s", params.PluginName, *strContext)
		log.Printf("[%s]%s Flush called with context", params.PluginName, params.InstanceName)
	} else {
		params = envToParams(PluginName)
		log.Printf("[%s]FLBPluginExitCtx fell back to env\n", PluginName)
		if params == nil {
			log.Printf("[%s] FLBPluginExitCtx no params\n", PluginName)
			return output.FLB_ERROR
		}
	}
	releaseResources()
	return output.FLB_OK
}

// This is the final callback used- and we need to release any retained resources.
// In our case - there is nothing to do
//
//export FLBPluginUnregister
func FLBPluginUnregister(def unsafe.Pointer) {
	log.Print("[out_gdb] Unregister called")
	output.FLBPluginUnregister(def)
}
