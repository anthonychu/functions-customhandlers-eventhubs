package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]interface{}
}

func eventHubTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var invokeRequest InvokeRequest

	d := json.NewDecoder(r.Body)
	d.Decode(&invokeRequest)

	messagesJson := invokeRequest.Data["eventHubMessages"]
	var outerMessages string
	var messages []string
	_ = json.Unmarshal(messagesJson, &outerMessages)
	_ = json.Unmarshal([]byte(outerMessages), &messages)

	outputs := map[string]interface{}{"": ""}
	invokeResponse := InvokeResponse{outputs, nil, ""}

	for _, message := range messages {
		invokeResponse.Logs = append(invokeResponse.Logs, "eventid:"+message)
	}

	responseJson, _ := json.Marshal(invokeResponse)
	// fmt.Println(string(responseJson))
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/EventHubTrigger", eventHubTriggerHandler)
	fmt.Println("Go server Listening on: ", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
