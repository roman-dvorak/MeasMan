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

func homeHandler(devices map[string]DevicesStatus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {


		m := map[string]interface{}{
			"Devices": devices,
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




func main() {
	//itteration := 0

	devices := make(map[string]DevicesStatus)

	devices["A"] = DevicesStatus{1, 0}
	devices["B"] = DevicesStatus{10, 10}

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
			fmt.Println("\t Monitor:", mon_i, " > ", mon.AppName)
			_, _ = s.Every(mon.Interval).Do(task2)
		}
	}

	s.StartAsync()

	fs := http.FileServer(http.Dir("./MeasMon/templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/a", homeHandler(devices))
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
