package main

//dependencies
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/google/uuid"
)

var (
	PUBSUB_NAME = "pubsub"
	TOPIC_NAME  = "audit"
	BROKERNAME = os.Getenv("BROKERNAME")
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
	// Additional Metadata for benchmarking
	BrokerName            string            `json:"brokerName,omitempty"`
	Timestamp             string            `json:"timestamp,omitempty"`
	UserAgent             string            `json:"userAgent,omitempty"`
	CorrelationId         string            `json:"correlationId,omitempty"`
	PublishedTimestamp    string            `json:"publishedTimestamp,omitempty"`
	ReceivedTimestamp     string            `json:"receivedTimestamp,omitempty"`
    BatchCorrelationId    string            `json:"batchCorrelationId,omitempty"`
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

func (r *Dapr) Send() {
	var kustoIngestionClient =  NewKustoIngestionClient()

	value := os.Getenv("NUMBER")
	log.Printf("sending %s messages", value)
	total, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error getting number: %v", err)
		total = 1000000
	}
	ctx := context.Background()

    // Always run
	for true {
		log.Println("starting publish")
		start_time := time.Now()

        // Each batch has a separate correlation id
        batchCorrelationId := uuid.New().String();

		for i := 0; i < total; i ++ {
			msg := getObj(strconv.Itoa(i), batchCorrelationId)

			for _, c := range r.client {

				// Send to ingestion service to track drop rate
				err = kustoIngestionClient.SendAsync(msg)
				if err != nil {
					// Continue if ingestion error, to do track count of these
					log.Printf("Ingestion error for batchCorrelationId: %s, %v", batchCorrelationId, err)
				}

                // Set the published timestamp
                msg.PublishedTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
                jsonData, err := json.Marshal(msg)
                if err != nil {
                    log.Fatalf("error marshaling data: %v", err)
                }

                // Using Dapr SDK to publish a topic and send to Subscriber
                // What is the publisher's streaming rate? Does the batching above significantly impact this rate
                if err := c.client.PublishEvent(ctx, PUBSUB_NAME, TOPIC_NAME, jsonData); err != nil {
                    panic(err)
                }
			}
		}

		end_time := time.Now()
		log.Printf("total time it took %v", end_time.Sub(start_time))
		
        fmt.Println("Sleep starting now: ", time.Now())
		time.Sleep(15*time.Minute);
	}
}
