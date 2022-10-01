package webHandle

import (
	"errors"
	"fmt"
	"log"
	home "server/homeManager"
	"strings"
)

func (h *Handler) regHubList(d jsonHub, clientId string) error {
	n := home.Hub{}

	n.Description = d.Description
	n.Id = d.Id
	n.ClientId = clientId
	n.Name = d.Name
	n.Help = d.Help

	for _, v := range d.Devices {
		if err := h.regDeviceList(v, clientId); err == nil {
			n.Devices = append(n.Devices, v.Id)
		} else {
			return err
		}
	}

	fmt.Println("current-devices:", h.HomeManager.GetDevices())
	return h.HomeManager.AddHub(n)
}

func (h *Handler) regDeviceList(d jsonDevice, clientId string) error {
	n := home.Device{}

	n.Description = d.Description
	n.Id = d.Id
	n.ClientId = clientId
	n.Name = d.Name
	n.Help = d.Help
	// build properties list for each device
	for _, v := range d.Properties {
		if thisType, ok := v["type"]; ok {
			switch strings.ToLower(thisType.(string)) {
			case "switch":
				if prop, err := home.ReadPropertySwitch(v); err == nil {
					if n.PropertySwitch == nil {
						n.PropertySwitch = make(map[string]home.SwitchProperty)
					}
					if _, ok := n.PropertySwitch[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.PropertySwitch[prop.Name] = prop
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}

			case "dial":
				if prop, err := home.ReadPropertyDial(v); err == nil {
					if n.PropertyDial == nil {
						n.PropertyDial = make(map[string]home.DialProperty)
					}
					if _, ok := n.PropertyDial[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.PropertyDial[prop.Name] = prop
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}
			}
		}
	}

	for _, v := range d.Uploads {
		n.Uploads = append(n.Uploads, home.Upload{
			Name:  v.Name,
			Alias: v.Alias,
		})
	}

	err := h.HomeManager.AddDevice(n)
	return err
}
