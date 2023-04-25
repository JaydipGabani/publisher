package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type RawMsgs struct {
	Records []PubsubMsg `json:"records,omitempty"`
}

type KustoIngestionClient struct {
    mutex sync.Mutex
    requestUri string
    pubSubMsgs []PubsubMsg
    counter int
    batchSize int
}

func NewKustoIngestionClient() KustoIngestionClient {
    targetUri := os.Getenv("TARGETURI")
	if targetUri == "" {
		log.Fatalf("Failed to set target uri env var")
	}

	functionKey := os.Getenv("FUNCTIONKEY")
	if functionKey == "" {
		log.Fatalf("Failed to set function key env var")
	}

    INGESTIONSERVICEBATCHSIZE := os.Getenv("INGESTIONSERVICEBATCHSIZE")
    batchSize, err := strconv.Atoi(INGESTIONSERVICEBATCHSIZE)
    if INGESTIONSERVICEBATCHSIZE == "" || err != nil {
        log.Fatalf("Failed to set INGESTIONSERVICEBATCHSIZE %v", err)
    }

    requestUri := targetUri + "?code=" + functionKey;

    pubSubMsgs := make([]PubsubMsg, 0, batchSize)

	return KustoIngestionClient { 
        requestUri:requestUri,
        pubSubMsgs: pubSubMsgs,
        batchSize: batchSize,
    }
}

func (client *KustoIngestionClient) SendAsync(msg PubsubMsg) error {
	client.mutex.Lock()
    defer client.mutex.Unlock()

    client.pubSubMsgs = append(client.pubSubMsgs, msg);
    client.counter += 1

    if client.counter == client.batchSize {
        msgs := make([]PubsubMsg, client.batchSize, client.batchSize)
        copy(msgs, client.pubSubMsgs)

        client.counter = 0
        client.pubSubMsgs = make([]PubsubMsg, 0, client.batchSize)
        return client.SendAsyncBatch(msgs)
    }

	return nil
}

func (client *KustoIngestionClient) SendAsyncBatch(msgs []PubsubMsg) error {
    rawMsgs := RawMsgs{Records: msgs}

	buf, err := json.Marshal(rawMsgs)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = http.Post(client.requestUri, "application/json", bytes.NewBuffer(buf))
		if err == nil {
			break
		}
		fmt.Printf("Error publishing count %d with error: %v \n", i, err)
	}

	if err != nil {
		return err
	}

    fmt.Println("Response Status:", resp.Status)

    return nil
}
