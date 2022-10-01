package webHandle

import "encoding/json"

type Generic struct {
	Method *string
	Data   *json.RawMessage
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
	Name        string
	Description string
	Id          string
	Help        string
	Devices     []jsonDevice
}

type jsonDevice struct {
	Id          string
	Name        string
	Description string
	Help        string
	// Groups      []jsonGroup
	Properties []map[string]interface{}
	Uploads    []jsonUpload
	// Actions     []jsonAction
}

type jsonGroup struct {
	Name string
}

type jsonUpload struct {
	Name  string
	Alias []string
}

type jsonAction struct {
	Name string
}

type eventData struct {
	Id         string
	Properties []map[string]interface{}
	Timestamp  string
}

type JsonEvent struct {
	Method string
	Data   eventData
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
