package main

type ConfigFile struct {
	Name    string          `yaml:"name"`
	Web     int             `yaml:"web"`
	//Devices []ConfigDevices `yaml:"devices,flow"`
	Devices map[string]ConfigDevices `yaml:"devices"`
}

type MonitorLocation struct {
	DeviceName string
	MonitorName string
}

type ConfigDevices struct {
	DeviceName string           `yaml:"device_name"`
	Url        string           `yaml:"url"`
	Executor   string           `yaml:"executor"`
	Protocol   string           `yaml:"protocol"`
	System     string           `yaml:"system"`
	//MonitorsRaw   []ConfigMonitors `yaml:"monitors,flow"`
	Monitors   map[string]ConfigMonitors `yaml:"monitors"`
	Status     int
}



type ConfigMonitors struct {
	MonitorType string `yaml:"type"`
	AppName     string `yaml:"app_name"`
	Interval    int    `yaml:"interval"`
	Identificator string
	Status      int
	Parameters map[string]interface{}
	//Device      ConfigDevices
}

type MonitorCheckAppRunning struct {
	Pid      []int
	PName    string
	PCount   int

}

type DevicesStatus struct {
	Status int
	Informations int
}

type TmplHomeData struct {
	Devices []DevicesStatus

}
