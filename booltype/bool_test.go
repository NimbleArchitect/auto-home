package booltype

import "testing"

func TestSet(t *testing.T) {
	type Result struct {
		data           string
		expectedString string
		expectedBool   bool
		expectedType   uint
	}

	table := []Result{
		{data: "OFF", expectedType: ONOFF, expectedString: "off", expectedBool: false},
		{data: "ON", expectedType: ONOFF, expectedString: "on", expectedBool: true},

		{data: "UP", expectedType: UPDOWN, expectedString: "up", expectedBool: true},
		{data: "DOWN", expectedType: UPDOWN, expectedString: "down", expectedBool: false},

		{data: "OPEN", expectedType: OPENCLOSE, expectedString: "open", expectedBool: true},
		{data: "CLOSE", expectedType: OPENCLOSE, expectedString: "close", expectedBool: false},

		{data: "YES", expectedType: YESNO, expectedString: "yes", expectedBool: true},
		{data: "NO", expectedType: YESNO, expectedString: "no", expectedBool: false},

		{data: "TRUE", expectedType: TRUEFALSE, expectedString: "true", expectedBool: true},
		{data: "FALSE", expectedType: TRUEFALSE, expectedString: "false", expectedBool: false},

		{data: "invalid", expectedType: TRUEFALSE, expectedString: "false", expectedBool: false},
	}

	for i, v := range table {
		var b BoolType

		b.Set(v.data)
		if b.String() != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, b.String())
		}
		if b.GetBool() != v.expectedBool {
			t.Fatalf("check %d expected \"%t\", recieved \"%t\"", i, v.expectedBool, b.GetBool())
		}
		if b.kind != v.expectedType {
			t.Fatalf("check %d expected \"%d\", recieved \"%d\"", i, v.expectedType, b.kind)
		}
	}
}

func TestSetBool(t *testing.T) {
	type Result struct {
		data           bool
		expectedString string
		expectedBool   bool
		expectedType   uint
	}

	table := []Result{
		{data: true, expectedType: TRUEFALSE, expectedString: "true", expectedBool: true},
		{data: false, expectedType: TRUEFALSE, expectedString: "false", expectedBool: false},
	}

	for i, v := range table {
		var b BoolType

		b.SetBool(v.data)
		if b.String() != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, b.String())
		}
		if b.GetBool() != v.expectedBool {
			t.Fatalf("check %d expected \"%t\", recieved \"%t\"", i, v.expectedBool, b.GetBool())
		}
		if b.kind != v.expectedType {
			t.Fatalf("check %d expected \"%d\", recieved \"%d\"", i, v.expectedType, b.kind)
		}
	}
}
