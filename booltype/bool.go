package booltype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BoolType struct {
	state bool
	kind  *booleanStrings
}

type booleanStrings struct {
	strTrue  string
	strFalse string
}

// list of true/false values as string representations
var val []*booleanStrings = []*booleanStrings{
	{strTrue: "true", strFalse: "false"},
	{strTrue: "open", strFalse: "close"},
	{strTrue: "yes", strFalse: "no"},
	{strTrue: "on", strFalse: "off"},
	{strTrue: "up", strFalse: "down"},
}

func (b *BoolType) Set(v string) {
	rawString := strings.ToLower(v)

	for _, v := range val {
		if v.strTrue == rawString {
			b.state = true
			b.kind = v
			return
		}
		if v.strFalse == rawString {
			b.state = false
			b.kind = v
			return
		}
	}

	fmt.Println("bool.Set: invalid state defaulting to", val[0].strFalse)
	b.kind = val[0]
	b.state = false
}

func (b *BoolType) String() string {
	if b.kind == nil {
		return "false"
	}

	if b.state {
		return b.kind.strTrue
	}
	return b.kind.strFalse
}

func (b *BoolType) SetBool(value bool) {
	if b.kind == nil {
		b.kind = val[0]
	}
	b.state = value
}

func (b *BoolType) Bool() bool {
	return b.state
}

func (b BoolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BoolType) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	b.Set(statusStr)
	return nil
}
