package main

/*
#include <stdlib.h>
*/
import (
	"C"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/input"
)

const PluginName = "in_gdb"

func getParams(plugin unsafe.Pointer) (*SqlParams, error) {
	params := SqlParams{}
	params.PluginName = PluginName
	params.Host = input.FLBPluginConfigKey(plugin, Plugin_Host)
	params.Port = input.FLBPluginConfigKey(plugin, Plugin_Port)
	params.User = input.FLBPluginConfigKey(plugin, Plugin_User)
	params.Password = input.FLBPluginConfigKey(plugin, Plugin_Password)
	params.SequencerCol = input.FLBPluginConfigKey(plugin, Plugin_Ordering)
	params.TableName = input.FLBPluginConfigKey(plugin, Plugin_TableName)
	params.DBName = input.FLBPluginConfigKey(plugin, Plugin_DBName)
	params.DBType = input.FLBPluginConfigKey(plugin, Plugin_Type)
	params.PK = input.FLBPluginConfigKey(plugin, Plugin_PK)
	params.ColsCSV = input.FLBPluginConfigKey(plugin, Plugin_ColsCSV)
	params.WhereExpr = input.FLBPluginConfigKey(plugin, Plugin_WhereExpr)

	freqStr := input.FLBPluginConfigKey(plugin, Plugin_QueryFrequency)
	if len(freqStr) > 0 {
		freq, err := strconv.Atoi(freqStr)
		if err != nil {
			return nil, err
		} else {
			params.QueryFrequency = freq
		}
	}

	params.DeleteAfterQuery = strings.Contains(strings.ToLower(input.FLBPluginConfigKey(plugin, Plugin_Delete)), "true")

	return &params, nil
}

func cacheParams(params *SqlParams) {
	paramsToEnv(params, PluginName)
}

func releaseResources() error {
	clearEnvParams(PluginName)
	return nil
}

func retrieveParams() *SqlParams {
	params := envToParams(PluginName)
	return params
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	log.Printf("[%s] Register called", PluginName)
	return input.FLBPluginRegister(def, PluginName, "Go plugin for reading content from a database")
}

// (fluentbit will call this)
// plugin (context) pointer to fluentbit context (state/ c code)
//
//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	params, err := getParams(plugin)
	if err != nil {
		return input.FLB_ERROR
	}
	if strings.Contains(strings.ToLower(input.FLBPluginConfigKey(plugin, "Log_Level")), "debug") {
		fmt.Printf("[%s] configured with %v\n", params.PluginName, params)
	}

	validateErr := validateSqlParams(params)
	if validateErr == nil {
		cacheParams(params)
		log.Printf(SprintfParams(params, PluginName))

		return input.FLB_OK
	} else {
		fmt.Printf("[%s] - Configuration error - %s \n", params.PluginName, validateErr)
		return input.FLB_ERROR
	}

}

//export FLBPluginInputCallback
func FLBPluginInputCallback(data *unsafe.Pointer, size *C.size_t) int {
	//log.Printf("FLBPluginInputCallback - START --------------")
	now := time.Now()
	params := retrieveParams()

	flbTime := input.FLBTime{now}

	dataSet, sequenceId := dynamicQuery(params)

	if dataSet != nil {
		if len(sequenceId) > 0 {
			params.LatestSequencerId = sequenceId
			// as we're using the last key - we need to update our cache
			cacheParams(params)
		}
		dataCtr := len(dataSet)

		log.Printf("[%s] InputCallback no records: %v\n", PluginName, dataCtr)
		if dataCtr > 0 {
			var entry []interface{}
			for dataLine := 0; dataLine < dataCtr; dataLine++ {
				recd := []interface{}{flbTime, dataSet[dataLine]}
				entry = []interface{}{flbTime, recd}
				log.Printf("[%s] InputCallback - retrieved data %v\n", PluginName, entry)
			}

			enc := input.NewEncoder()
			packed, err := enc.Encode(entry)
			if err != nil {
				log.Printf("[%s] error: %s,\n Can't convert to msgpack: %v\n", PluginName, err, entry)
				return input.FLB_ERROR
			}

			length := len(packed)
			*data = C.CBytes(packed)
			*size = C.size_t(length)
		} else {
			length := 0
			*data = nil
			*size = C.size_t(length)
		}
	} else {
		log.Printf("[%s] InputCallback -- no data found", PluginName)

		// For emitting interval adjustment.
		time.Sleep(time.Second * time.Duration(params.QueryFrequency))
	}

	//log.Printf("FLBPluginInputCallback - END ==========")
	return input.FLB_OK
}

//export FLBPluginInputCleanupCallback
func FLBPluginInputCleanupCallback(data unsafe.Pointer) int {
	return input.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	err := releaseResources()
	if err != nil {
		log.Printf("%s had an error during cleanup, error is %s\n", PluginName, err)
		return input.FLB_ERROR
	}
	return input.FLB_OK
}
