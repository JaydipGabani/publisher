package main

//dependencies
import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dapr/go-sdk/client"
)

// code
var (
	PUBSUB_NAME = "pubsub"
	TOPIC_NAME  = "audit"
)

type Client struct {
	client client.Client
}

// Dapr represents driver for interacting with pub sub using dapr.
type Dapr struct {
	// Array of clients to talk to different endpoints
	client []Client
}

type PubsubMsg struct {
	Key                   string            `json:"key,omitempty"`
	ID                    string            `json:"id,omitempty"`
	Details               interface{}       `json:"details,omitempty"`
	EventType             string            `json:"eventType,omitempty"`
	Group                 string            `json:"group,omitempty"`
	Version               string            `json:"version,omitempty"`
	Kind                  string            `json:"kind,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Namespace             string            `json:"namespace,omitempty"`
	Message               string            `json:"message,omitempty"`
	EnforcementAction     string            `json:"enforcementAction,omitempty"`
	ConstraintAnnotations map[string]string `json:"constraintAnnotations,omitempty"`
	ResourceGroup         string            `json:"resourceGroup,omitempty"`
	ResourceAPIVersion    string            `json:"resourceAPIVersion,omitempty"`
	ResourceKind          string            `json:"resourceKind,omitempty"`
	ResourceNamespace     string            `json:"resourceNamespace,omitempty"`
	ResourceName          string            `json:"resourceName,omitempty"`
	ResourceLabels        map[string]string `json:"resourceLabels,omitempty"`
}

func main() {
	r := &Dapr{}
	tmpc, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	newClient := Client{
		client: tmpc,
	}
	r.client = append(r.client, newClient)
	log.Println("client initialized")

	r.Send()
	time.Sleep(300 * time.Hour)
}

func getObj(i string) interface{} {
	return PubsubMsg{
		Key:                   i,
		ID:                    "2023-04-06T21:17:05Z",
		Details:               map[string]interface{}{"missing_labels": []interface{}{"test"}},
		EventType:             "violation_audited",
		Group:                 "constraints.gatekeeper.sh",
		Version:               "v1beta1",
		Kind:                  "K8sRequiredLabels",
		Name:                  "pod-must-have-test",
		Namespace:             "",
		Message:               "you must provide labels: {\"test\"}",
		EnforcementAction:     "deny",
		ConstraintAnnotations: map[string]string(nil),
		ResourceGroup:         "",
		ResourceAPIVersion:    "v1",
		ResourceKind:          "Pod",
		ResourceNamespace:     "nginx",
		ResourceName:          "dywuperf-deployment-10kpods-69bd64c867-h2wdx",
		ResourceLabels:        map[string]string{"app": "dywuperf-app-100kpods", "pod-template-hash": "69bd64c867"}}
}

func (r *Dapr) Send() {
	value := os.Getenv("NUMBER")
	log.Printf("sending %s messages", value)
	total, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error getting number: %v", err)
		total = 100000
	}
	ctx := context.Background()
	start_time := time.Now()
	log.Println("starting publish")
	
	for i := 0; i < total; i++ {
		test := getObj(strconv.Itoa(i))
		jsonData, err := json.Marshal(test)
		if err != nil {
			log.Fatalf("error marshaling data: %v", err)
		}

		for _, c := range r.client {
			//Using Dapr SDK to publish a topic
			// log.Println("Published data: " + string(jsonData))
			if err := c.client.PublishEvent(ctx, PUBSUB_NAME, TOPIC_NAME, jsonData); err != nil {
				panic(err)
			}
		}
	}
	end_time := time.Now()
	log.Printf("total time it took %v", end_time.Sub(start_time))
}
