package main

import (
	"fmt"
	"log"
	"sync"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/core"
)

type TopicWriterChan struct {
	Value string
}

type TopicMessageListener struct {
	wg *sync.WaitGroup
}

var (
	nameTopicExample = "topic.capitals"
)

func TopicWriter(name, address string, topicMessageListener *TopicMessageListener, writerChan chan TopicWriterChan) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress(address)
	config.SetClientName(name)

	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		log.Println(err)
		return nil
	}

	Topic, _ := client.GetTopic(nameTopicExample)
	Topic.AddMessageListener(topicMessageListener)

	for {
		temp, open := <-writerChan
		if temp.Value != "" && open {
			Topic.Publish(temp.Value)
			log.Println(fmt.Sprintf(" %s | %s | WRITES | %s", name, address, temp.Value))
		} else {
			break
		}
	}

	Topic.Destroy()
	client.Shutdown()
	log.Println(fmt.Sprintf(" %s | %s | SHUT DOWN", name, address))
	return nil
}

func (l *TopicMessageListener) OnMessage(message core.Message) error {
	fmt.Println("Got message: ", message.MessageObject())
	fmt.Println("Publishing Time: ", message.PublishTime())
	l.wg.Done()
	return nil
}
