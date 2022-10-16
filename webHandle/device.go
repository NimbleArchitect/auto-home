package webHandle

import (
	"errors"
	"log"
	"server/deviceManager"
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

	// fmt.Println(">> current-devices:", h.HomeManager.GetDevices())
	return h.HomeManager.AddHub(n)
}

func (h *Handler) regDeviceList(d jsonDevice, clientId string) error {
	n := deviceManager.NewDevice(h.HomeManager.MaxPropertyHistory)
	n.Description = d.Description
	n.Id = d.Id
	n.ClientId = clientId
	n.Name = d.Name
	n.Help = d.Help

	window := h.HomeManager.DeviceWindow(d.Name)

	// build properties list for each device
	for _, v := range d.Properties {
		if thisType, ok := v["type"]; ok {
			switch strings.ToLower(thisType.(string)) {
			case "switch":
				if prop, err := n.Map2Switch(v); err == nil {
					if _, ok := n.PropertySwitch[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.SetSwitch(prop.Name, prop)
						n.SetSwitchWindow(prop.Name, window[prop.Name])
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}

			case "dial":
				if prop, err := n.Map2Dial(v); err == nil {
					if _, ok := n.PropertyDial[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.SetDial(prop.Name, prop)
						n.SetDialWindow(prop.Name, window[prop.Name])
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}

			case "button":
				if prop, err := n.Map2Button(v); err == nil {
					if _, ok := n.PropertyButton[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.SetButton(prop.Name, prop)
						n.SetButtonWindow(prop.Name, window[prop.Name])
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}

			case "text":
				if prop, err := n.Map2Text(v); err == nil {
					if _, ok := n.PropertyText[prop.Name]; !ok {
						log.Println("adding property", prop.Name)
						n.SetText(prop.Name, prop)
						n.SetTextWindow(prop.Name, window[prop.Name])
					} else {
						return errors.New("duplicate property name detected, peoperty " + prop.Name + " is already in use")
					}
				} else {
					return err
				}

			}
		}
	}

	// for _, v := range d.Uploads {
	// 	n.Uploads = append(n.Uploads, home.Upload{
	// 		Name:  v.Name,
	// 		Alias: v.Alias,
	// 	})
	// }

	err := h.HomeManager.AddDevice(n, clientId)
	return err
}
