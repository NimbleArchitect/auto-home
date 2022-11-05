package pluginManager

import (
	"errors"
	"testing"
)

func TestMakeError(t *testing.T) {
	type Result struct {
		dataId         int
		dataErr        error
		expectedString string
	}

	table := []Result{
		{dataId: 123, dataErr: nil, expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true}}`},
		{dataId: 123, dataErr: errors.New("test error"), expectedString: `{"Method":"result","Id":123,"Data":{"Ok":false,"Message":"test error"}}`},
	}

	for i, v := range table {
		out := makeError(v.dataId, v.dataErr)
		strOut := string(out)
		if strOut != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, strOut)
		}
	}
}

func TestMakeResponse(t *testing.T) {
	type Result struct {
		dataId         int
		dataArgs       interface{}
		expectedString string
	}

	argArr := []string{"val1", "val2"}

	argMap1 := make(map[string]interface{})
	argMap1["val1"] = "out1"
	argMap1["val2"] = 2

	table := []Result{
		{dataId: 123, dataArgs: nil,
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":null}}}`},
		{dataId: 123, dataArgs: "trial string",
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":"trial string"}}}`},
		{dataId: 123, dataArgs: 21,
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":21}}}`},
		{dataId: 123, dataArgs: true,
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":true}}}`},
		{dataId: 123, dataArgs: argArr,
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":["val1","val2"]}}}`},
		{dataId: 123, dataArgs: argMap1,
			expectedString: `{"Method":"result","Id":123,"Data":{"Ok":true,"Data":{"0":{"val1":"out1","val2":2}}}}`},
	}

	for i, v := range table {
		out := makeResponse(v.dataId, v.dataArgs)
		strOut := string(out)
		if strOut != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, strOut)
		}
	}
}
