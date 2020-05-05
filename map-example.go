package main

import (
	"fmt"
	"log"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
)

type MapWriterChan struct {
	Key   interface{}
	Value interface{}
}

var nameExampleMap = "map.capitals"

func MapWriter(name, address string, writerChan chan MapWriterChan, pinger chan bool) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress(address)
	config.SetClientName(name)

	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		temp, open := <-writerChan
		if temp.Key != nil && open {
			Map, _ := client.GetMap(nameExampleMap)
			Map.Put(temp.Key, temp.Value)
			log.Println(fmt.Sprintf(" %s | %s | WRITES | %s : %s", name, address, temp.Key, temp.Value))
			pinger <- true
		} else {
			break
		}
	}

	client.Shutdown()
	log.Println(fmt.Sprintf(" %s | %s | SHUT DOWN", name, address))
	return nil
}

func MapReader(name, address string, pinger chan bool) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress(address)
	config.SetClientName(name)

	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	thisMap := make(map[interface{}]interface{})

	for {
		ok, open := <-pinger
		if open && ok {
			Map, _ := client.GetMap(nameExampleMap)
			keysSet, _ := Map.KeySet()
			var keys = make([]interface{}, 0)
			for _, k := range keysSet {
				keys = append(keys, k)
			}
			for _, key := range keys {
				if _, exist := thisMap[key]; !exist {
					value, _ := Map.Get(key)
					thisMap[key] = value
					log.Println(fmt.Sprintf(" %s | %s | READS | %s : %s", name, address, key, value))
				}
			}
		} else {
			break
		}
	}

	client.Shutdown()
	log.Println(fmt.Sprintf(" %s | %s | SHUT DOWN", name, address))

	return nil
}
