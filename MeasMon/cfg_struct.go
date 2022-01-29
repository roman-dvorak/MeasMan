package main

type ConfigFile struct {
	Name    string          `yaml:"name"`
	Web     int             `yaml:"web"`
	Devices []ConfigDevices `yaml:"devices,flow"`
}

type ConfigDevices struct {
	DeviceName string           `yaml:"device_name"`
	Url        string           `yaml:"url"`
	Executor   string           `yaml:"executor"`
	Protocol   string           `yaml:"protocol"`
	System     string           `yaml:"system"`
	Monitors   []ConfigMonitors `yaml:"monitors,flow"`
}



type ConfigMonitors struct {
	MonitorType string `yaml:"type"`
	AppName     string `yaml:"app_name"`
	Interval    int    `yaml:"interval"`
	Status      int
	Device      ConfigDevices
}

type DevicesStatus struct {
	Status int
	Informations int
}

type TmplHomeData struct {
	Devices []DevicesStatus

}
