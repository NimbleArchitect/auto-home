package js

const (
	FLAG_CONTINUEPROCESSING = iota
	FLAG_GROUPPROCESSING
	FLAG_STOPPROCESSING
	FLAG_PREVENTUPDATE
)

type jsFlag struct {
	flag int
}

func (f *jsFlag) Has(flag int) bool {
	return f.flag == flag
}

func (f *jsFlag) Not(flag int) bool {
	return f.flag != flag
}

// set the should we process flag
func (f *jsFlag) Set(flag int64) {
	f.flag = int(flag)
}
