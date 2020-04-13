package main

import (
	"fmt"
	"log"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
)

type QueueWriterChan struct {
	Value string
}

var (
	nameExampleQueue = "queue.capitals"
)

func QueueWriter(name, address string, writerChan chan QueueWriterChan, pinger chan bool) error {

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

	for {
		temp, open := <-writerChan
		if temp.Value != "" && open {
			Queue, _ := client.GetQueue(nameExampleQueue)
			Queue.Put(temp.Value)
			log.Println(fmt.Sprintf(" %s | %s | WRITES | %s", name, address, temp.Value))
			pinger <- true
		} else {
			break
		}
	}

	client.Shutdown()
	log.Println(fmt.Sprintf(" %s | %s | SHUT DOWN", name, address))
	return nil
}

func QueueReader(name, address string, pinger chan bool) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress(address)
	config.SetClientName(name)

	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for {
		ok, open := <-pinger
		if open && ok {
			Queue, _ := client.GetQueue(nameExampleQueue)
			value, _ := Queue.Take()
			log.Println(fmt.Sprintf(" %s | %s | READS | %s", name, address, value))
		} else {
			break
		}
	}

	client.Shutdown()
	log.Println(fmt.Sprintf(" %s | %s | SHUT DOWN", name, address))

	return nil
}
