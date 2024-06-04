package main

// this file represents the common logic between the two plugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/ziutek/mymysql"
)

const Plugin_InstanceId = "plugin_instance_id"
const Plugin_Host = "db_host"
const Plugin_Port = "db_port"
const Plugin_User = "db_user"
const Plugin_Password = "db_password"
const Plugin_Ordering = "ordering_col"
const Plugin_TableName = "table_name"
const Plugin_DBName = "db_name"
const Plugin_Type = "db_type"
const Plugin_PK = "pk"
const Plugin_Limit = "limit"
const Plugin_Delete = "delete"
const Plugin_WhereExpr = "where_expression"
const Plugin_ColsCSV = "query_cols"
const Plugin_QueryFrequency = "query_frequency"
const Plugin_LatestSequencerId = "LstSeqId"

// https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go
// https://go101.org/article/struct.html
type SqlParams struct {
	PluginName       string `json:"pgnname,omitempty"` // the name of the plugin that the configuration applies to
	InstanceName     string `json:"instNme,omitempty"` // to help differentiate the activities of different occurrences of the plugin we can provide an identifying name
	Host             string `json:"host,omitempty"`    // the server running the DB
	Port             string `json:"port,omitempty"`    // The port to use to connect to the DB
	User             string `json:"usr,omitempty"`     // uasername to connect to the DB with
	Password         string `json:"pw,omitempty"`      // the password to use when connecting to the DB
	DBName           string `json:"dbnme,omitempty"`   // the database name
	ColsCSV          string `json:"cols,omitempty"`    // comma separated list pf the columns we want put or get for the named table
	SequencerCol     string `json:"seqr,omitempty"`    // the column which determines correct record sequence - so that we get the records in the right order
	TableName        string `json:"tbl,omitempty"`     // name of the table in the database
	WhereExpr        string `json:"where,omitempty"`   // any additional statements to make yup a where statement, no need for the word 'where'
	DeleteAfterQuery bool   `json:"del,omitempty"`     // defines whether any records read should then be deleted once retrieved
	PK               string `json:"pk,omitempty"`      // The primary key of the table - necessary to drive the deletion
	DBType           string `json:"dbtype,omitempty"`  // The database type mysql, postgres
	QueryFrequency   int    `json:"freq,omitempty"`    // the number of seconds until the next query assuming all existing records have been retrieved

	//the following attributes are for operational caching purposes and aren't reflected in the configuration
	LatestSequencerId string `json:"seqrId,omitempty"`
}

const PostgresDBType = "postgres"
const mysqlDBType = "mysql"
const ParamCGFPostfix = "-cfg"
const InsertTimeout = time.Second * 1

// return a string representation of the data type
func printType(description string, data interface{}) {
	xType := fmt.Sprintf("%T", data)
	fmt.Println(description, ":", xType)
}

// Provide a simple function to create our config data, then if we need to ibitialise any values we can
// incorporate it in a single place
func NewSqlParams() *SqlParams {
	params := SqlParams{}
	// if you want to asset default values - then follow the details
	// https://pkg.go.dev/gopkg.in/mcuadros/go-defaults.v1#section-readme
	return &params
}

// each element is represented with this definiot - with the column name being the string key
// by sharing the data with the column name and the data in its native type rather than a string helps
// with subsequent perocessing which may want to know what the base daat type is
type recordValType map[string]interface{}

// to ensure we get a proper converstion between non string data types and strings
// we can use this function to ensure they're handled ok
func typeToStr(data interface{}, quoteStrings bool) string {

	if data == nil {
		log.Printf("typeToStr - defensive check - data is nil")
		return ""
	}

	switch (data).(type) {
	case int:
		return strconv.Itoa(data.(int))
	case uint8:
	case []uint8:
		if quoteStrings {
			return "'" + fmt.Sprintf("%s", data) + "'"
		}
		return fmt.Sprintf("%s", data)

	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case int32:
		return strconv.FormatInt(data.(int64), 10)
	case uint64:
		return strconv.FormatUint(data.(uint64), 10)
	case bool:
		return strconv.FormatBool(data.(bool))
	case float64:
		return strconv.FormatFloat(data.(float64), 'E', -1, 64)
	case float32:
		return strconv.FormatFloat(data.(float64), 'E', -1, 32)
	case string:
		if quoteStrings {
			return "'" + data.(string) + "'"
		}
		return data.(string)
	case uint:
		return strconv.FormatUint(data.(uint64), 10)
	default:
		printType("data type is", data)
		return fmt.Sprintf("%v", data)
	}

	log.Println("Dropped out of typeToStr")
	return ""
}

// Put the internal configuration values into a printable format
func SprintfParams(params *SqlParams, pluginName string) string {
	if params == nil {
		log.Printf("[%s] SprintfParams called with no params struct", pluginName)
		return ""
	}
	var paramStr string = paramsToJSON(params)
	paramStr = fmt.Sprintf("[%s]\"Connection\":{%s},\nQuery:%s\n", paramStr, buildConnectionStr(params), buildQueryExpr(params, false))
	return paramStr
}

// Build a JSON representation of our configuration and other values we'd like to hold in our context
func paramsToJSON(params *SqlParams) string {
	json, err := json.Marshal(*params)
	if err != nil {
		log.Printf("[%s] paramsToJSON error - %s", params.PluginName, err)
	}

	return fmt.Sprintf(string(json))
}

// Convert a JSON representation of our context data back to the relevant data structure
func JSONToParams(paramsJSON string, pluginName string) *SqlParams {
	var params SqlParams
	err := json.Unmarshal([]byte(paramsJSON), &params)
	if err != nil {
		log.Printf("[%s] JSONToParams error - %s", pluginName, err)
	}
	return &params
}

// if we're using environmental variables as a caching mechanism, then we need to be able to clear those values out
func clearEnvParams(pluginName string) {
	blankParams := NewSqlParams()
	log.Printf("[%s] Flushing environment params", pluginName)
	paramsToEnv(blankParams, pluginName)
}

// store the parameter values as envcironment veriables - this is part of a work around for the missing context on the input side of the plugin
// NOTE: if any additional elements added into the SqlParams struct then they need to be factored into this and its opposite function
// We do note perform any value validation here - as we assume this is done during the initialization phase, and the env vars are NOT tampered with
func paramsToEnv(params *SqlParams, pluginName string) error {
	var err error = nil // use this if at some point we need to communicate an error
	os.Setenv(pluginName+"_"+Plugin_Host, params.Host)
	os.Setenv(pluginName+"_"+Plugin_Port, params.Port)
	os.Setenv(pluginName+"_"+Plugin_User, params.User)
	os.Setenv(pluginName+"_"+Plugin_Password, params.Password)
	os.Setenv(pluginName+"_"+Plugin_Ordering, params.SequencerCol)
	os.Setenv(pluginName+"_"+Plugin_TableName, params.TableName)
	os.Setenv(pluginName+"_"+Plugin_DBName, params.DBName)
	os.Setenv(pluginName+"_"+Plugin_Type, params.DBType)
	os.Setenv(pluginName+"_"+Plugin_PK, params.PK)
	os.Setenv(pluginName+"_"+Plugin_QueryFrequency, strconv.Itoa(params.QueryFrequency))
	os.Setenv(pluginName+"_"+Plugin_Delete, strconv.FormatBool(params.DeleteAfterQuery))
	os.Setenv(pluginName+"_"+Plugin_ColsCSV, (params.ColsCSV))
	os.Setenv(pluginName+"_"+Plugin_WhereExpr, (params.WhereExpr))
	os.Setenv(pluginName+"_"+Plugin_LatestSequencerId, (params.LatestSequencerId))
	return err
}

// This takes our params construct and caches the values as environment variables. Each value has to be a separate value as
// trying to use a JSON construct has the potential to be too long, resulting in string truncation and
// the loss of information and the unmarshalling failing
func envToParams(pluginName string) *SqlParams {
	params := NewSqlParams()
	params.Host = os.Getenv((pluginName + "_" + Plugin_Host))
	params.Port = os.Getenv((pluginName + "_" + Plugin_Port))
	params.User = os.Getenv((pluginName + "_" + Plugin_User))
	params.Password = os.Getenv((pluginName + "_" + Plugin_Password))
	params.SequencerCol = os.Getenv((pluginName + "_" + Plugin_Ordering))
	params.TableName = os.Getenv((pluginName + "_" + Plugin_TableName))
	params.DBName = os.Getenv((pluginName + "_" + Plugin_DBName))
	params.DBType = os.Getenv((pluginName + "_" + Plugin_Type))
	params.PK = os.Getenv((pluginName + "_" + Plugin_PK))
	params.ColsCSV = os.Getenv((pluginName + "_" + Plugin_ColsCSV))
	params.WhereExpr = os.Getenv((pluginName + "_" + Plugin_WhereExpr))
	params.LatestSequencerId = os.Getenv((pluginName + "_" + Plugin_LatestSequencerId))

	freqStr := strings.TrimSpace(os.Getenv(pluginName + "_" + Plugin_QueryFrequency))
	if len(freqStr) > 0 {
		var err error = nil
		params.QueryFrequency, err = strconv.Atoi(freqStr)
		if err != nil {
			log.Printf("[%s] envToParams fault, and defaulting queryFrequency - received %s, error is %s\n", pluginName, freqStr, err)
		}
	} else {
		params.QueryFrequency = 1
	}

	delStr := strings.TrimSpace(os.Getenv(pluginName + "_" + Plugin_Delete))
	if len(delStr) > 0 {
		params.DeleteAfterQuery = strings.Contains(strings.ToLower(delStr), "true")
	}

	return params
}

// creates the correct conection string based on DB type
func buildConnectionStr(params *SqlParams) string {
	var connectStr string = ""

	switch params.DBType {
	case PostgresDBType:
		connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", params.Host, params.Port, params.User, params.Password, params.DBName)
	case mysqlDBType:
		connectStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", params.User, params.Password, params.Host, params.Port, params.DBName)
	default:
		log.Println("[%s] Unknown db type >", params.DBType, "<", params.PluginName)
		connectStr = ""
	}

	return connectStr
}

// Once the values have been retrieved and loaded into the SqlParams struct we need to verify whether the mandatory
// elements have meaningful values and return an error if not. This will help the user when testing configurations
// TODO: extend so we can validate based on in/out/filter
func validateSqlParams(params *SqlParams) error {

	if len(params.PluginName) == 0 {
		log.Printf("validateSqlParams - defaulting plugin name\n")
		params.PluginName = "gdb"
	}

	// make sure the host has been set
	params.Host = strings.TrimSpace(params.Host)
	if len(params.Host) == 0 {
		return errors.New("No " + Plugin_Host + " defined for " + params.PluginName)
	}

	// remove any white space and confirm there is a value
	params.Port = strings.TrimSpace(params.Port)
	if len(params.Port) == 0 {
		return errors.New("No " + Plugin_Port + " defined for " + params.PluginName)
	}

	// test port is numeric
	_, err := strconv.Atoi(params.Port)
	if err != nil {
		return errors.New(Plugin_Port + " is not numeric for " + params.PluginName)
	}

	// remove any white space and confirm there is a value
	params.User = strings.TrimSpace(params.User)
	if len(params.User) == 0 {
		return errors.New("No " + Plugin_User + " defined for " + params.PluginName)
	}

	// remove any white space and confirm there is a value
	params.DBName = strings.TrimSpace(params.DBName)
	if len(params.DBName) == 0 {
		return errors.New("No " + Plugin_DBName + " name defined for " + params.PluginName)
	}

	// if a comma separated list of columns is not provided then set the columns to be a wildcard
	params.ColsCSV = strings.TrimSpace(params.ColsCSV)
	if len(params.ColsCSV) == 0 {
		params.ColsCSV = "*"
		log.Printf("[%s]Defaulting query columns to %s\n", params.PluginName, params.ColsCSV)
	}

	params.DBType = strings.TrimSpace(params.DBType)
	if len(params.DBType) == 0 {
		return errors.New("No " + Plugin_Type + " defined for " + params.PluginName)
	} else {
		switch params.DBType {
		case PostgresDBType:
			// nothing to do we know about this DB type
		case mysqlDBType:
			// nothing to do we know about this DB type
		default:
			params.DBType = ""
			return errors.New("Unknown " + Plugin_Type + " defined  " + params.DBType + " for " + params.PluginName)
		}
	}

	// test query interval is numeric and default if not set
	if params.QueryFrequency <= 0 {
		params.QueryFrequency = 1
	}

	return nil
}

// This builds the SQL expression. Uses standard ANSI SQL, but could be customized for optimization
// based on other DBs if so desired
// The predefined SQL is capitalized so it will stand out when reviewing
// The countStmt allows us to generate a record count statement that can be execute rather than
// retrieving the actual record elements. This is necessary to defend against MyDQL driver
// throwing a fatal error when it tries to handle a rowset with now values.
func buildQueryExpr(params *SqlParams, countStmt bool) string {
	var sqlStmt string = "SELECT " + params.ColsCSV + " FROM " + params.TableName
	var whereStmt string = ""

	if countStmt {
		sqlStmt = "SELECT COUNT(*) FROM " + params.TableName
	}

	if len(params.WhereExpr) > 0 {
		whereStmt = " WHERE " + params.WhereExpr
	}

	if len(params.LatestSequencerId) > 0 && !params.DeleteAfterQuery {
		exprStr := " AND "
		if len(params.WhereExpr) == 0 {
			exprStr = " WHERE "
		}
		whereStmt = exprStr + params.SequencerCol + " > " + params.LatestSequencerId

	}
	sqlStmt = sqlStmt + whereStmt

	if !countStmt {
		if len(params.SequencerCol) > 0 {
			sqlStmt = sqlStmt + " ORDER BY " + params.SequencerCol
		}
		sqlStmt = sqlStmt + " LIMIT 1" //+ strconv.Itoa(params.Limit)
	}
	log.Printf("[%s]%s Query constructed:%s", params.PluginName, params.InstanceName, sqlStmt)

	return sqlStmt
}

// Create the delete SQL statement for removing data values
func buildDeleteExpr(params *SqlParams, recdKey string) string {
	var sqlStmt = "DELETE FROM " + params.TableName + " WHERE " + params.PK + " = " + recdKey

	return sqlStmt
}

// create a transaction with delete statements using the retrieved pk (primary key)
func execDelete(params *SqlParams, keyList []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cancel()

	db, err := sql.Open(params.DBType, buildConnectionStr(params))
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	for rowValIdx := 0; rowValIdx < len(keyList); rowValIdx++ {
		sqlStmt := buildDeleteExpr(params, typeToStr(keyList[rowValIdx], true))
		if err == nil {
			_, err := tx.ExecContext(ctx, sqlStmt)
			fmt.Println(sqlStmt)
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			log.Println(err)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type RowDefinition map[interface{}]interface{}
type ManyRowDefinition []RowDefinition

// build up the SQL statement, as we don't know whether we're popukating the entire DB row
// we need to use the column names.
// We assume that the values provided are in the same order as the defined params to target
// TODO: address the possibility that the values provided mismatch the column names
// TODO: handle the possibility that the params is empty or a wildcard
func buildInsertExpr(params *SqlParams, values RowDefinition) (string, error) {
	if values == nil || len(values) == 0 {
		return "", errors.New("No data values provided")
	}

	var valLen int = len(values)
	var orderedColNames []string = make([]string, valLen)
	var colnames string = params.ColsCSV

	if colnames == "*" {
		colnames = ""
		var ctr int = 0
		for key, _ := range values {
			orderedColNames[ctr] = key.(string)
			ctr++
			if len(colnames) == 0 {
				colnames = key.(string)
			} else {
				colnames = colnames + "," + key.(string)
			}
		}
	} else {
		fmt.Println(" ")
		splitStr := strings.Split(params.ColsCSV, ",")
		if len(splitStr) != len(orderedColNames) {
			return "", errors.New("Number of data elements to no cols mismatched")
		}
	}

	var sqlStmt = "INSERT INTO " + params.TableName + " (" + colnames + ")"

	var valueStr string = ""
	//for _, val := range values {
	for valIdx := 0; valIdx < valLen; valIdx++ {
		if valIdx == 0 {
			valueStr = typeToStr(values[orderedColNames[valIdx]], true)
		} else {
			// if it isnt the last element we need a comma to separate the values
			valueStr = valueStr + "," + typeToStr(values[orderedColNames[valIdx]], true)
		}
	}
	sqlStmt = sqlStmt + " VALUES (" + valueStr + ")"
	return sqlStmt, nil
}

// get the SQL generated and execute the statement
// as we want to potentially execute multiple rows - let's wrap it
// inside a transaction with a time out
// the data needs to be presented as an array of string arrays
func execInsert(params *SqlParams, value RowDefinition) error {
	ctx, cancel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cancel()

	db, err := sql.Open(params.DBType, buildConnectionStr(params))
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStmt, err := buildInsertExpr(params, value)
	if err == nil {
		_, err := tx.ExecContext(ctx, sqlStmt)
		fmt.Printf("[%s]%s insert expression: %s", params.PluginName, params.InstanceName, sqlStmt)
		if err != nil {
			log.Println("[%s]%s Error with insert %v", params.PluginName, params.InstanceName, err)
			return err
		}
	} else {
		log.Println("[%s]%s SQL context error %v", params.PluginName, params.InstanceName, err)
	}
	//}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// without resorting to a full query validate that the conection details will work.
func testConnectionOk(params *SqlParams) bool {
	db, err := sql.Open(params.DBType, buildConnectionStr(params))

	if db == nil {
		log.Printf("[%s]%s connection test failed for\n%s\n%v", params.PluginName, params.InstanceName, SprintfParams(params, params.PluginName), err)
		return false
	}

	if err = db.Ping(); err != nil {
		db.Close()
		log.Printf("[%s]%s connection test ping failed for\n%s\n%v", params.PluginName, params.InstanceName, SprintfParams(params, params.PluginName), err)
		return false
	}
	return true
}

// check the DB has records worth retrieving before we actual pull them back - see above for more detail on
// the nature of this issue.
// TODO: Determine whether this is just a MySQL driver issue of also true of Postgres etc
func checkForData(db *sql.DB, sqlExpr string) (bool, error) {
	var count int = 0
	//log.Printf("checkForData stmt %s", sqlExpr)

	// Query for a value based on a single row.
	if err := db.QueryRow(sqlExpr).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("Unexpected no rows response on execDataCheck")
		}
		return false, fmt.Errorf("Unexpected err in execDataCheck %v", err)
	}
	return count > 0, nil
}

// Executes the SQL statement and dynamically resolves the number of columns that maybe retrieved
// based on https://kylewbanks.com/blog/query-result-to-map-in-golang
// func execQuery(sqlExpr string, sequencerCol string, db *sql.DB) (map[string]interface{}, string, error) {
func execQuery(sqlExpr string, sequencerCol string, pk string, db *sql.DB) ([]interface{}, []interface{}, string, error) {

	dbRows, err := db.Query(sqlExpr)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("execQuery - no data")
			return nil, nil, "", nil

		}
		log.Printf("execQuery - err from db query call: %s", err)
		return nil, nil, "", err
	}
	defer dbRows.Close()

	colNames, err := dbRows.Columns()
	if err != nil || colNames == nil {
		log.Printf("execQuery - error during retrieval of columns: %s", err)
		return nil, nil, "", err
	}

	var myData []interface{} = nil
	var myKeys []interface{} = nil
	var lastSequenceValue *interface{} = nil

	if dbRows == nil {
		log.Printf("execQuery empty result set from query")
		return nil, nil, "", nil
	}

	// for MySQL we will get a SQL error if there are no rows retrieved
	dbErr := dbRows.Err()
	if dbErr != nil {
		log.Printf("execQuery - DbErr:%v", dbErr)
		return nil, nil, "", nil

	}

	for dbRows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		//myMap := make(map[string]interface{})
		myMap := make(recordValType, 1)

		columns := make([]interface{}, len(colNames))
		columnPointers := make([]interface{}, len(colNames))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		err := dbRows.Scan(columnPointers...)
		if err != nil {
			log.Printf("execQuery error scanning:%s", err)
			return nil, nil, "", err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		for i, colName := range colNames {
			val := columnPointers[i].(*interface{})
			*val = typeToStr(*val, false)
			if colName == sequencerCol {
				lastSequenceValue = val
			}

			// if value is the identified primary then add the value to the myKeys array
			if colName == pk {
				if myKeys == nil {
					myKeys = make([]interface{}, 1)
					myKeys[0] = *val
				} else {
					myKeys = append(myKeys, *val)
				}
			}
			myMap[colName] = *val
			//printType(colName, *val)

		}

		log.Printf("execQuery row being sent = %v", myMap)

		if myData == nil {
			myData = make([]interface{}, 1)

			myData[0] = myMap
		} else {
			myData = append(myData, myMap)
		}

		rowErr := dbRows.Err()
		if rowErr != nil {
			log.Printf("execQuery - row error %s", rowErr)
			return myData, myKeys, typeToStr(*lastSequenceValue, false), nil
		}
	}

	return myData, myKeys, typeToStr(*lastSequenceValue, false), nil

}

// builds the relevant connections and executes the query
// it then translates the resultant structure to a JSON output
func dynamicQuery(params *SqlParams) ([]interface{}, string) {
	db, err := sql.Open(params.DBType, buildConnectionStr(params))
	if err != nil {
		log.Printf("gdb - dynamicQuery - received an error during open, about to panic\n%v", err)
		panic(err.Error())
	}
	defer db.Close()

	hasData, err := checkForData(db, buildQueryExpr(params, true))
	if err == nil && hasData {

		result, keyList, lastSeqId, err := execQuery(buildQueryExpr(params, false), params.SequencerCol, params.PK, db)
		if err != nil {
			log.Printf("dynamicQuery - received an error from execQuery about to panic")
			panic(err)
		}

		//fmt.Println("Result=", result)
		log.Printf("KeyList=%s Last Sequence Id=%s,  delete is %t\n", keyList, lastSeqId, params.DeleteAfterQuery)

		if keyList != nil && params.DeleteAfterQuery {
			execDelete(params, keyList)
		}
		return result, lastSeqId
	} else {
		if err != nil {
			log.Printf("dynamicQuery - check for data error %v", err)
		}
		return nil, params.LatestSequencerId
	}
}
