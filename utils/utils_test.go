package utils

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
	TestStructList map[string]TestStruct
}

type TestStruct struct {
	Field1 string
	Field2 int
	Field3 []string
	Field4 float64
	Field5 *string
	field6 string
}

func (suite *UtilsTestSuite) SetUpUtilsTestSuite() {

	// initialize the map
	suite.TestStructList = make(map[string]TestStruct)

	tStr := "tStr"

	ts1 := TestStruct{"44", 44, []string{"t1", "t2"}, 12.33, &tStr, "unexported"}
	ts2 := TestStruct{"", 44, []string{"t1", "t2"}, 12.33, &tStr, "unexported"}

	// fill the map
	suite.TestStructList["ts1"] = ts1
	suite.TestStructList["ts2"] = ts2

}

func (suite *UtilsTestSuite) TestCheckForNulls() {

	// tests the normal case
	suite.Equal(nil, CheckForNulls(suite.TestStructList["ts1"]))

	// tests the case of an object containing a field which is empty
	suite.Equal(errors.New("utils.TestStruct object contains an empty value for field: Field1"), CheckForNulls(suite.TestStructList["ts2"]))
}

func (suite *UtilsTestSuite) TestGetFieldByName() {

	// tests the normal case
	val1, err1 := GetFieldValueByName(suite.TestStructList["ts1"], "Field1")

	// tests the case of a missing field
	val2, err2 := GetFieldValueByName(suite.TestStructList["ts1"], "Field10")

	// tests the case of an empty field
	val3, err3 := GetFieldValueByName(suite.TestStructList["ts2"], "Field1")

	// tests the case of an unexported field
	val4, err4 := GetFieldValueByName(suite.TestStructList["ts2"], "field6")

	suite.Equal("44", val1.(string))
	suite.Equal(nil, val2)
	suite.Equal(nil, val3)
	suite.Equal(nil, val4)

	suite.Equal(nil, err1)
	suite.Equal("Field: Field10 has not been declared.", err2.Error())
	suite.Equal("utils.TestStruct object contains an empty value for field: Field1", err3.Error())
	suite.Equal("you are trying to access an unexported field", err4.Error())

}

func (suite *UtilsTestSuite) TestStructToMap(){

	//tests the normal case with unexported field
	tStr := "tStr"
	expMap := map[string]interface{}{"Field1": "44", "Field2": 44, "Field3": []string{"t1", "t2"}, "Field4" : 12.33, "Field5": &tStr}
	suite.Equal(expMap, StructToMap(suite.TestStructList["ts1"]))

	//tests the case of nil input
	suite.Nil(StructToMap(nil))

}

func (suite *UtilsTestSuite) TestIsCapitalized(){

	suite.Equal(true, IsCapitalized("Str1"))
	suite.Equal(false, IsCapitalized("str1"))
	suite.Equal(false, IsCapitalized(""))
}

func TestUtilsTestSuite(t *testing.T) {
	utilsTestSuite := new(UtilsTestSuite)
	utilsTestSuite.SetUpUtilsTestSuite()
	suite.Run(t, utilsTestSuite)

}
