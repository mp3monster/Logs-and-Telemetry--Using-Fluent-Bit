package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func createInsertTestData() []string {
	var testdata []string = make([]string, 5)
	testdata[0] = strconv.Itoa(rand.Intn(1000))
	testdata[1] = "test data " + strconv.Itoa(rand.Intn(10))
	testdata[2] = strconv.Itoa(rand.Intn(100))
	v, _ := time.Now().UTC().MarshalText()
	testdata[3] = fmt.Sprintf(string(v))
	testdata[4] = strconv.Itoa(rand.Intn(10000)) + "." + strconv.Itoa(rand.Intn(1000))
	fmt.Println(testdata)
	return testdata
}

func main() {
	params := SqlParams{}
	params.Host = "192.168.1.135"
	params.Port = "5455"
	params.User = "postgresUser"
	params.Password = "postgresPW"
	params.SequencerCol = "a_key"
	params.TableName = "pluginsrc"
	params.DBName = "postgres"
	params.DBType = PostgresDBType
	params.PK = "a_key"
	params.DeleteAfterQuery = true

	err := validateSqlParams(&params)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	dynamicQuery(&params)

	params.ColsCSV = "a_key, a_string, a_number, a_dtg, a_decimal"
	//var testdata []string = createInsertTestData()

	const noTestDatas = 3
	var testdataSet [][]string = make([][]string, noTestDatas)
	for i := 0; i < noTestDatas; i++ {
		testdataSet[i] = createInsertTestData()
	}
	fmt.Println(testdataSet)
	//execInsert(params, testdataSet)
	//fmt.Println(buildInsertExpr(testdata))
}
