package js

import "server/groupManager"

type jsGroup struct {
	Id        string
	Name      string
	groups    []string
	devices   []string
	liveGroup *groupManager.Group
}

func (g *jsGroup) HasDevice(name string) interface{} {

	return nil
}

func (g *jsGroup) GetDevice(name string) interface{} {

	return nil
}

func (g *jsGroup) GetGroup(name string) interface{} {

	return nil
}

func (g *jsGroup) GetGroupByPath(name string) interface{} {

	return nil
}

func (g *jsGroup) GetDeviceByPath(name string) interface{} {

	return nil
}

func (g *jsGroup) GetDeviceInGroup(name string) interface{} {

	return nil
}
