package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func mvWhenDone(filepath string, dest string, fileName string, user string, wait int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * time.Duration(wait))
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size0 := int64(0)
	newSize := fi.Size()

	for newSize != size0 {
		time.Sleep(time.Second * time.Duration(wait))
		fi, err = file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		size0 = newSize
		newSize = fi.Size()
	}
	// Move file option 1
	//os.Rename(event.Name, fmt.Sprintf("%s/%s", dest, fileName))

	// Move file option 2
	cmd := exec.Command("mv", filepath, fmt.Sprintf("%s/%s", dest, fileName))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("chown", user, fmt.Sprintf("%s/%s", dest, fileName))
	_ = cmd.Run()
}

func main() {
	// Parse config flag
	flag.Parse()

	// Read configuration file
	var config = readConfig()

	// Create log file
	f, err := os.OpenFile(config.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	// Set log file as output default. Comment line for print to stdout
	log.SetOutput(f)

	// Set directory listener
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Prepare channels to stop
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	done := make(chan bool)

	// Run process in other goroutine
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Printf("Event %v -- file %s", event.Op, event.Name)
				if event.Op == fsnotify.Create {
					log.Printf("file created; %s", event.Name)
					// -------- A NEW FILE WAS CREATED -------------
					// Prepare dir path of event and name of the file
					dir := event.Name[:strings.LastIndex(event.Name, "/")]
					fileName := event.Name[strings.LastIndex(event.Name, "/")+1:]

					// Check if path is in the json file
					if dest, ok := dirMap[dir]; ok {
						go mvWhenDone(event.Name, dest.OutPath, fileName, dest.User, config.Wait)
					} else {
						log.Printf("Not found destination match for dir: %s", dir)
					}
				}

			case err = <-watcher.Errors:
				log.Println("error:", err)

			case signalType := <-ch:
				signal.Stop(ch)
				log.Println("Exit command received. Exiting...")
				log.Println("Received signal type : ", signalType)

				done <- true
			}
		}
	}()

	// Load dirs from json
	loadDirMap(config.Jsonpaths)

	// Add dirs to watch
	for key := range dirMap {
		// First walk
		files, _ := ioutil.ReadDir(key)
		for _, f := range files {
			if !f.IsDir() {
				log.Printf("file found; %s", f.Name())
				// -------- A NEW FILE WAS CREATED -------------
				// Prepare dir path of event and name of the file
				dir := key
				fileName := f.Name()

				// Check if path is in the json file
				if dest, ok := dirMap[dir]; ok {
					go mvWhenDone(fmt.Sprintf("%s/%s", dir, fileName), dest.OutPath, fileName, dest.User, config.Wait)
				} else {
					log.Printf("Not found destination match for dir: %s", dir)
				}
			}
		}

		// add to watcher
		err = watcher.Add(key)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print(dirMap)

	<-done
	os.Exit(0)
}
