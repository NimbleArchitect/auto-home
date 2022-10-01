package homeClient

type Hub struct {
	data hubData
}

type hubData struct {
	Name        string
	Description string
	Id          string
	Devices     []Device
}

func NewHub(name string, deviceid string) Hub {

	return Hub{
		data: hubData{
			Name: name,
			Id:   deviceid,
		},
	}
}

func (h *Hub) AddDevice(device Device) error {

	h.data.Devices = append(h.data.Devices, device)

	return nil
}
