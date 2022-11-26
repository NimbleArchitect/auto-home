package home

type Hub struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ClientId    string   `json:"clientid"`
	Help        string   `json:"help"`
	Devices     []string `json:"devices"`
}

// type group struct {
// 	Id          string
// 	Name        string
// 	Description string
// 	Devices     []string
// 	Groups      []string
// 	Users       []string
// }

type Action struct {
	Name     string
	Location string
	// jsvm     *js.JavascriptVM
}
