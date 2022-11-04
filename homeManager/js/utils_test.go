package js

import (
	"testing"
)

func TestMapToJsSwitch(t *testing.T) {

	type Result struct {
		data          map[string]interface{}
		shouldFail    bool
		expectedName  string
		expectedValue string
		expectedState bool
	}

	table := []Result{
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "value1",
			},
			expectedName:  "thisName",
			expectedValue: "false",
			expectedState: false,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "open",
			},
			expectedName:  "thisName",
			expectedValue: "open",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "no",
			},
			expectedName:  "thisName",
			expectedValue: "no",
			expectedState: false,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "true",
			},
			expectedName:  "thisName",
			expectedValue: "true",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": true,
			},
			expectedName:  "thisName",
			expectedValue: "true",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  false,
				"value": true,
			},
			shouldFail: true,
		},
		{
			data: map[string]interface{}{
				"name":  0,
				"value": true,
			},
			shouldFail: true,
		},
	}

	for i, check := range table {
		result, err := MapToJsSwitch(check.data)
		if check.shouldFail {
			if err == nil {
				t.Fatalf("check %d, expected error but recieved nil", i)
			}
		} else {
			if err != nil {
				t.Fatalf("check %d,expected nil but recieved error: %s", i, err.Error())
			}
			if result.Name != check.expectedName {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedName, result.Name)
			}
			if result.Value != check.expectedValue {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedValue, result.Value)
			}
			if result.state != check.expectedState {
				t.Fatalf("expected \"%t\", recieved \"%t\"", check.expectedState, result.state)
			}
		}
	}

}

func TestMapToJsButton(t *testing.T) {

	type Result struct {
		data          map[string]interface{}
		shouldFail    bool
		expectedName  string
		expectedValue string
		expectedState bool
	}

	table := []Result{
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "value1",
			},
			expectedName:  "thisName",
			expectedValue: "false",
			expectedState: false,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "open",
			},
			expectedName:  "thisName",
			expectedValue: "open",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "no",
			},
			expectedName:  "thisName",
			expectedValue: "no",
			expectedState: false,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "true",
			},
			expectedName:  "thisName",
			expectedValue: "true",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": true,
			},
			expectedName:  "thisName",
			expectedValue: "true",
			expectedState: true,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  false,
				"value": true,
			},
			shouldFail: true,
		},
		{
			data: map[string]interface{}{
				"name":  0,
				"value": true,
			},
			shouldFail: true,
		},
	}

	for i, check := range table {
		result, err := MapToJsButton(check.data)
		if check.shouldFail {
			if err == nil {
				t.Fatalf("check %d, expected error but recieved nil", i)
			}
		} else {
			if err != nil {
				t.Fatalf("check %d,expected nil but recieved error: %s", i, err.Error())
			}
			if result.Name != check.expectedName {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedName, result.Name)
			}
			if result.Value != check.expectedValue {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedValue, result.Value)
			}
			if result.state != check.expectedState {
				t.Fatalf("expected \"%t\", recieved \"%t\"", check.expectedState, result.state)
			}
		}
	}

}

func TestMapToJsDial(t *testing.T) {

	type Result struct {
		data          map[string]interface{}
		shouldFail    bool
		expectedName  string
		expectedValue int
	}

	table := []Result{
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": 4,
			},
			expectedName:  "thisName",
			expectedValue: 4,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": 5000,
			},
			expectedName:  "thisName",
			expectedValue: 5000,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": -234,
			},
			expectedName:  "thisName",
			expectedValue: -234,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "1024",
			},
			expectedName:  "thisName",
			expectedValue: 1024,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "0",
			},
			expectedName:  "thisName",
			expectedValue: 0,
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  false,
				"value": true,
			},
			shouldFail: true,
		},
		{
			data: map[string]interface{}{
				"name":  0,
				"value": true,
			},
			shouldFail: true,
		},
	}

	for i, check := range table {
		result, err := MapToJsDial(check.data)
		if check.shouldFail {
			if err == nil {
				t.Fatalf("check %d, expected error but recieved nil", i)
			}
		} else {
			if err != nil {
				t.Fatalf("check %d,expected nil but recieved error: %s", i, err.Error())
			}
			if result.Name != check.expectedName {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedName, result.Name)
			}
			if result.Value != check.expectedValue {
				t.Fatalf("expected \"%d\", recieved \"%d\"", check.expectedValue, result.Value)
			}
		}
	}

}

func TestMapToJsText(t *testing.T) {

	type Result struct {
		data          map[string]interface{}
		shouldFail    bool
		expectedName  string
		expectedValue string
	}

	table := []Result{
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": 4,
			},
			shouldFail: true,
		},
		{
			data: map[string]interface{}{
				"name":  "thisName",
				"value": "0",
			},
			expectedName:  "thisName",
			expectedValue: "0",
			shouldFail:    false,
		},
		{
			data: map[string]interface{}{
				"name":  false,
				"value": true,
			},
			shouldFail: true,
		},
		{
			data: map[string]interface{}{
				"name":  0,
				"value": true,
			},
			shouldFail: true,
		},
	}

	for i, check := range table {
		result, err := MapToJsText(check.data)
		if check.shouldFail {
			if err == nil {
				t.Fatalf("check %d, expected error but recieved nil", i)
			}
		} else {
			if err != nil {
				t.Fatalf("check %d,expected nil but recieved error: %s", i, err.Error())
			}
			if result.Name != check.expectedName {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedName, result.Name)
			}
			if result.Value != check.expectedValue {
				t.Fatalf("expected \"%s\", recieved \"%s\"", check.expectedValue, result.Value)
			}
		}
	}

}
