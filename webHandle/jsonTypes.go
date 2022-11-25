package webHandle

import "encoding/json"

type Generic struct {
	Method *string          `json:"methos"`
	Data   *json.RawMessage `json:"data"`
}

type jsonPlugin struct {
	Id          string
	Name        string
	Description string
	Help        string
	Kind        string
}

type jsonConnector struct {
	Client string
	Token  string
}

type jsonHub struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Id          string       `json:"id"`
	Help        string       `json:"help"`
	Devices     []jsonDevice `json:"devices"`
}

type jsonDevice struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Help        string `json:"help"`
	// Groups      []jsonGroup
	Properties []map[string]interface{} `json:"properties"`
	Uploads    []jsonUpload             `json:"uploads"`
	// Actions     []jsonAction
}

type jsonGroup struct {
	Name string
}

type jsonUpload struct {
	Name  string   `json:"name"`
	Alias []string `json:"alias"`
}

type jsonAction struct {
	Name string
}

type eventData struct {
	Id         string                   `json:"id"`
	Properties []map[string]interface{} `json:"properties"`
	Timestamp  string                   `json:"timestamp"`
}

type JsonEvent struct {
	Method string    `json:"method"`
	Data   eventData `json:"data"`
}

// func (g *Generic) GetMethod() string {
// 	if g.Method != nil {
// 		return *g.Method
// 	}
// 	return *g.M
// }

// func (g *Generic) GetData() interface{} {
// 	if g.Data != nil {
// 		return *g.Data
// 	}
// 	return *g.D
// }
