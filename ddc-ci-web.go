package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	defer syscall.FreeLibrary(user32)
	defer syscall.FreeLibrary(dxva2)
	http.HandleFunc("/s/hdmi", func(writer http.ResponseWriter, request *http.Request) {
		SetMonitorValue(InPutSource, HDMI)
		writer.WriteHeader(200)
	})
	http.HandleFunc("/s/dp", func(writer http.ResponseWriter, request *http.Request) {
		SetMonitorValue(InPutSource, DP)
		writer.WriteHeader(200)
	})
	http.HandleFunc("/b/", func(writer http.ResponseWriter, request *http.Request) {
		id := strings.TrimPrefix(request.URL.Path, "/b/")
		atoi, err := strconv.Atoi(id)
		if err != nil || atoi <= 0 || atoi > 100 {
			writer.WriteHeader(400)
			writer.Write([]byte("Brightness: " + id + " is illegal"))
		} else {
			SetMonitorValue(Brightness, atoi)
			writer.WriteHeader(200)
		}
	})
	http.HandleFunc("/c/", func(writer http.ResponseWriter, request *http.Request) {
		id := strings.TrimPrefix(request.URL.Path, "/c/")
		atoi, err := strconv.Atoi(id)
		if err != nil || atoi <= 0 || atoi > 100 {
			writer.WriteHeader(400)
			writer.Write([]byte("Contrast: " + id + " is illegal"))
		} else {
			SetMonitorValue(Contrast, atoi)
			writer.WriteHeader(200)
		}
	})
	http.HandleFunc("/v/", func(writer http.ResponseWriter, request *http.Request) {
		id := strings.TrimPrefix(request.URL.Path, "/v/")
		atoi, err := strconv.Atoi(id)
		if err != nil || atoi <= 0 || atoi > 100 {
			writer.WriteHeader(400)
			writer.Write([]byte("Volume: " + id + " is illegal"))
		} else {
			SetMonitorValue(Volume, atoi)
			writer.WriteHeader(200)
		}
	})
	http.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		m := make(map[string]int, 0)
		m["Volume"] = GetMonitorValue(Volume)
		m["Brightness"] = GetMonitorValue(Brightness)
		m["Contrast"] = GetMonitorValue(Contrast)
		marshal, _ := json.Marshal(m)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		writer.Write(marshal)
	})
	log.Fatal(http.ListenAndServe(":1888", nil))
}
