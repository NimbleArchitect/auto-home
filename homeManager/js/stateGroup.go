package js

import (
	"server/groupManager"
)

// SetGroup initilises a new jsGroup object
func (r *JavascriptVM) SetGroup(id string, name string, groups []string, devices []string, liveGroup *groupManager.Group) {

	r.groups[id] = jsGroup{
		Id:        id,
		Name:      name,
		groups:    groups,
		devices:   devices,
		liveGroup: liveGroup,
	}
}
