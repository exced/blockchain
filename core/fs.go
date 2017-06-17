package core

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
)

// Save encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Load decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// Check
func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

// SaveOnInterrupt captures SIGTERM and save file before exiting the program
func SaveOnInterrupt(path string, object interface{}) {
	// Capture SIGTERM
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Printf("Saving file at %s", path)
		Save(path, object)
		cleanupDone <- true
	}()
	<-cleanupDone
}
