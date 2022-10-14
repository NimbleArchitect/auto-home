package deviceManager

const (
	BUTTON = iota
	DIAL
	SWITCH
	TEXT

	RO
	RW
	WO
)

// type Hub struct {
// 	Id          string
// 	Name        string
// 	Description string
// 	ClientId    string
// 	Help        string
// 	Devices     []string
// }

type group struct {
	Id          string
	Name        string
	Description string
	Devices     []string
	Groups      []string
	Users       []string
}

type Upload struct {
	Name  string
	Alias []string
}

//	type timeoutWindow struct {
//		Name  string
//		Prop  string
//		Value int64
//	}

// type SwitchProperty struct {
// 	Name        string
// 	Description string
// 	Value       booltype.BoolType
// 	Previous    booltype.BoolType
// 	Mode        uint
// }

// type ButtonProperty struct {
// 	Name        string
// 	Description string
// 	Value       booltype.BoolType
// 	Previous    bool
// 	Mode        uint
// }

// type TextProperty struct {
// 	Name        string
// 	Description string
// 	Value       string
// 	Previous    string
// 	Mode        uint
// }
