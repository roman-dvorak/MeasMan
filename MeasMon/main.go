package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"os"
	"io/ioutil"
	//"regexp"

	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v3"
	"time"
	//"reflect"
	badger "github.com/dgraph-io/badger/v3"
	//"github.com/biter777/processex"
	"os/exec"
	"strings"

	"strconv"
	"github.com/go-ping/ping"
)

func homeHandler(Devices map[string]ConfigDevices) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		m := map[string]interface{}{
			"Devices": Devices,
		}

		//fmt.Printf("DEVICES >> %+v\n\n", Devices)

		// for _, dev := range Devices {
		// 	//fmt.Printf("DEV >> %+v\n", dev.Monitors)
		// 	for _, mon := range dev.Monitors {
		// 		fmt.Printf("MON >> %+v\n", mon)
		// 	}
		// 	//use dev_i
		// }

		tmpl, _ := template.New("").ParseFiles("MeasMon/templates/index.html", "MeasMon/templates/base.html")
		tmpl.ExecuteTemplate(w, "base", m)

    }
}


func monitorPing(config *ConfigFile, locator MonitorLocation) (err error){
	m := config.Devices[locator.DeviceName].Monitors[locator.MonitorName]

	pinger, err := ping.NewPinger("192.168.1.1")

	pinger.Count = 3
	err = pinger.Run()
	stats := pinger.Statistics()
	fmt.Printf("PING >> %+v\n", stats.AvgRtt)

	m.Status = 0

	if err != nil {
		m.Status = 1
	} else {
		if (stats.AvgRtt/time.Millisecond) > 0 && (stats.AvgRtt/time.Millisecond) < 100 {
			m.Status = 2
		}
	}
	config.Devices[locator.DeviceName].Monitors[locator.MonitorName] = m

	return nil
}

func monitorCheckAppRunning(config *ConfigFile, locator MonitorLocation) (err error){
	m := config.Devices[locator.DeviceName].Monitors[locator.MonitorName]
	m.Status += 1


	if( len(m.Parameters) == 0) { // Initialize monitor
		m.Parameters = make(map[string]interface{})
		fmt.Println("Initialization of monitor check_app_running")

		var pids MonitorCheckAppRunning
		pids.PName = m.AppName

		// GET PID names
		out, _ := exec.Command("pidof", m.AppName).Output()
		pids_s := strings.Split(string(out), " ")

		for _, i := range pids_s{
			v, _ := strconv.Atoi(i)
			pids.Pid = append(pids.Pid, v)
		}

		m.Parameters["Init"] = 1
		m.Parameters["Pids"] = pids
		config.Devices[locator.DeviceName].Monitors[locator.MonitorName] = m
		return nil
	}

	pids:= m.Parameters["Pids"].(MonitorCheckAppRunning)

	// Get PID numbers from name
	out, _ := exec.Command("pidof", m.AppName).Output()
	pids_s := strings.Split(string(out), " ")

	pids.Pid = nil
	for _, i := range pids_s{
		v, _ := strconv.Atoi(i)
		pids.Pid = append(pids.Pid, v)
	}


	m.Parameters["Pids"] = pids
	config.Devices[locator.DeviceName].Monitors[locator.MonitorName] = m

	return nil
}

func monitorNoneTask(config *ConfigFile, locator MonitorLocation) (err error){
	fmt.Println("TASK UNKNOWN")

	return nil
}


func get_task(mon *ConfigMonitors) func(config *ConfigFile, locator MonitorLocation) (err error) {
	fmt.Println("Vyhledavam spoustec pro:", mon.MonitorType, mon.AppName)
	switch mon.MonitorType {
		case "check_app_running":
			return monitorCheckAppRunning
		case "ping":
			return monitorPing
	}

	return monitorNoneTask
}


func main() {
	// Load configuration file
	db, err := badger.Open(badger.DefaultOptions("/tmp/MeasMan"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	yfile, err := ioutil.ReadFile("config/conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var data ConfigFile
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal(err)
	}

	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()

	fmt.Println(data)

	for dev_i, dev := range data.Devices {
		fmt.Println("Device ", dev_i, dev)
		for mon_i, mon := range dev.Monitors {
			fmt.Printf("MON >> %+v\n", mon)
			locator := MonitorLocation{dev_i, mon_i}
			fmt.Printf("LOCATOR >> %+v\n", locator)

			var identificator = dev.DeviceName+"__"+mon.AppName
			mon.Identificator = identificator

			fmt.Println("\t Monitor:", mon_i, " > ", mon.AppName, mon.Interval)
			fmt.Println(">>", &mon)
			f := get_task(&mon)
			f(&data, locator)
			_, _ = s.Every(mon.Interval).Seconds().Do(f, &data, locator)

		}
	}

	fmt.Println("Spustit cron..")
	s.StartAsync()

	fs := http.FileServer(http.Dir("./MeasMon/templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler(data.Devices))
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
