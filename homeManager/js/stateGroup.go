package js

// SetGroup initilises a new jsGroup object
func (r *JavascriptVM) SetGroup(id string, name string, groups []string, devices []string) {

	r.groups[id] = jsGroup{
		Id:      id,
		Name:    name,
		groups:  groups,
		devices: devices,
	}
}
