package booltype

import "testing"

func TestSet(t *testing.T) {
	type Result struct {
		data           string
		expectedString string
		expectedBool   bool
	}

	table := []Result{
		{data: "OFF", expectedString: "off", expectedBool: false},
		{data: "ON", expectedString: "on", expectedBool: true},

		{data: "UP", expectedString: "up", expectedBool: true},
		{data: "DOWN", expectedString: "down", expectedBool: false},

		{data: "OPEN", expectedString: "open", expectedBool: true},
		{data: "CLOSE", expectedString: "close", expectedBool: false},

		{data: "YES", expectedString: "yes", expectedBool: true},
		{data: "NO", expectedString: "no", expectedBool: false},

		{data: "TRUE", expectedString: "true", expectedBool: true},
		{data: "FALSE", expectedString: "false", expectedBool: false},

		{data: "invalid", expectedString: "false", expectedBool: false},
	}

	for i, v := range table {
		var b BoolType

		b.Set(v.data)
		if b.String() != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, b.String())
		}
		if b.Bool() != v.expectedBool {
			t.Fatalf("check %d expected \"%t\", recieved \"%t\"", i, v.expectedBool, b.Bool())
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
		{data: true, expectedString: "true", expectedBool: true},
		{data: false, expectedString: "false", expectedBool: false},
	}

	for i, v := range table {
		var b BoolType

		b.SetBool(v.data)
		if b.String() != v.expectedString {
			t.Fatalf("check %d expected \"%s\", recieved \"%s\"", i, v.expectedString, b.String())
		}
		if b.Bool() != v.expectedBool {
			t.Fatalf("check %d expected \"%t\", recieved \"%t\"", i, v.expectedBool, b.Bool())
		}

	}
}
