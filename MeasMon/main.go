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
)

func homeHandler(configMonitor map[string]ConfigMonitors) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {


		m := map[string]interface{}{
			"Devices": configMonitor,
		}


		tmpl, _ := template.New("").ParseFiles("MeasMon/templates/index.html", "MeasMon/templates/base.html")
		tmpl.ExecuteTemplate(w, "base", m)

    }
}

func task(itteration *int) {
	*itteration += 1
	fmt.Print("TASK PRINTOUT ... ", *itteration, "\n")
}

func task2() {
	fmt.Print("TASK AAAAAAa ... \n")
}


func monitorCheckAppRunning(mon map[string]ConfigMonitors, name string) (err error){
	fmt.Println("TASK, check app running", mon[name])
	i := mon[name]
	i.Status += 1
	mon[name] = i

	return nil
}

func monitorNoneTask(mon map[string]ConfigMonitors, name string) (err error){
	fmt.Println("TASK UNKNOWN")

	return nil
}

// func get_task(mon *ConfigMonitors) func(mon *ConfigMonitors) {
// 	fmt.Println("Vyhledavam spoustec pro:", mon.MonitorType)
// 	switch mon.MonitorType {
// 		case "check_app_running":
// 			return monitorCheckAppRunning
// 	}
//
// 	return monitorNoneTask
// }

func get_task(mon *ConfigMonitors) func(mon map[string]ConfigMonitors, name string) (err error) {
	fmt.Println("Vyhledavam spoustec pro:", mon.MonitorType, mon.AppName)
	switch mon.MonitorType {
		case "check_app_running":
			return monitorCheckAppRunning
	}

	return monitorNoneTask
}


func main() {
	//itteration := 0

	devices := make(map[string]ConfigMonitors)
	//
	// devices["A"] = DevicesStatus{1, 0}
	// devices["B"] = DevicesStatus{10, 10}

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

	for dev_i, dev := range data.Devices {
		fmt.Println("Device ", dev_i)
		for mon_i, mon := range dev.Monitors {
			//mon.Device = dev
			devices[mon.AppName] =  mon

			fmt.Println("\t Monitor:", mon_i, " > ", mon.AppName, mon.Interval)
			//get_task(&mon)
			fmt.Println(">>", &mon)
			f := get_task(&mon)
			_, _ = s.Every(mon.Interval).Seconds().Do(f, devices, mon.AppName)

		}
	}

	fmt.Println("Spustit cron..")
	s.StartAsync()

	fs := http.FileServer(http.Dir("./MeasMon/templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/a", homeHandler(devices))
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
