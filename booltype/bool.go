package booltype

import (
	"encoding/json"
	"log"
	"strings"
)

const (
	strOn    = "on"
	strOff   = "off"
	strUp    = "up"
	strDown  = "down"
	strYes   = "yes"
	strNo    = "no"
	strOpen  = "open"
	strClose = "close"
	strTrue  = "true"
	strFalse = "false"
)

const (
	ONOFF = iota
	UPDOWN
	OPENCLOSE
	TRUEFALSE
	YESNO
)

type BoolType struct {
	state bool
	kind  uint
}

func (b *BoolType) String() string {
	switch b.kind {
	case ONOFF:
		if b.state {
			return strOn
		} else {
			return strOff
		}
	case UPDOWN:
		if b.state {
			return strUp
		} else {
			return strDown
		}
	case OPENCLOSE:
		if b.state {
			return strOpen
		} else {
			return strClose
		}
	case YESNO:
		if b.state {
			return strYes
		} else {
			return strNo
		}
	case TRUEFALSE:
		if b.state {
			return strTrue
		} else {
			return strFalse
		}
	}
	// should never get here
	return "NA"
}

func (s BoolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *BoolType) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	s.Set(statusStr)
	return nil
}

func (b *BoolType) Set(v string) {
	switch strings.ToLower(v) {
	case strOn:
		b.kind = ONOFF
		b.state = true
	case strOff:
		b.kind = ONOFF
		b.state = false

	case strUp:
		b.kind = UPDOWN
		b.state = true
	case strDown:
		b.kind = UPDOWN
		b.state = false

	case strOpen:
		b.kind = OPENCLOSE
		b.state = true
	case strClose:
		b.kind = OPENCLOSE
		b.state = false

	case strYes:
		b.kind = YESNO
		b.state = true
	case strNo:
		b.kind = YESNO
		b.state = false

	case strTrue:
		b.kind = TRUEFALSE
		b.state = true
	case strFalse:
		b.kind = TRUEFALSE
		b.state = false

	default:
		log.Println("bool.Set: invalid state defaulting to", strFalse)
		b.kind = TRUEFALSE
		b.state = false
	}
}

func (b *BoolType) SetBool(value bool) {
	b.kind = TRUEFALSE
	b.state = value
}

func (b *BoolType) GetBool() bool {
	return b.state
}
