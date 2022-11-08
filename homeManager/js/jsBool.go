package js

// jsBool returns boolean true/false by default also supports the toString method which
// converts to the used string eg open/close, up/down etc
type jsBool struct {
	s string
	b bool
}

func (j *jsBool) ToString() string {
	return j.s
}

func (j *jsBool) ValueOf() bool {
	return j.b
}
