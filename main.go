package main
//dependencies
import (
	"context"
	"log"
	"math/rand"
	"time"
	"strconv"
	"github.com/dapr/go-sdk/client"
)

//code
var (
	PUBSUB_NAME = "pubsub"
	TOPIC_NAME  = "orders"
)

type Client struct {
	client          client.Client
}

// Dapr represents driver for interacting with pub sub using dapr.
type Dapr struct {
	// Array of clients to talk to different endpoints
	client []Client
}

func main() {
	r := &Dapr{}
	tmpc, err := client.NewClient()
	if err != nil {
		log.Println(err)
	}
	newClient := Client{
		client:          tmpc,
	}
	r.client = append(r.client, newClient)
	log.Println("client initialized")
	for {
		r.Send()
		log.Println("5 messages sent")
		time.Sleep(3 * time.Second)
	}
}


func (r *Dapr) Send() {
	
		ctx := context.Background()
	for  i := 0; i < 5; i++ {
		orderId := rand.Intn(1000-1) + 1
		
		for _, c := range r.client {
    //Using Dapr SDK to publish a topic
		if err := c.client.PublishEvent(ctx, PUBSUB_NAME, TOPIC_NAME, []byte(strconv.Itoa(orderId))); 
		err != nil {
			log.Println(err)
			if err
		}

		log.Println("Published data: " + strconv.Itoa(orderId))
		time.Sleep(2 * time.Second)
	}
}
}