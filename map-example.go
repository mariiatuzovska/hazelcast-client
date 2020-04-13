package main

import (
	"fmt"
	"log"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
)

type MapWriterChan struct {
	Key   string
	Value string
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
		log.Println(err)
		return nil
	}

	for {
		temp, open := <-writerChan
		if temp.Key != "" && open {
			Map, err := client.GetMap(nameExampleMap)
			if err != nil {
				log.Println(err)
				return nil
			}
			_, err = Map.Put(temp.Key, temp.Value)
			if err != nil {
				log.Println(err)
				return nil
			}
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
		return nil
	}

	thisMap := make(map[string]string)

	for {
		ok, open := <-pinger
		if open && ok {
			Map, err := client.GetMap(nameExampleMap)
			if err != nil {
				log.Println(err)
				return nil
			}
			keysSet, err := Map.KeySet()
			var keys = []string{}
			for _, k := range keysSet {
				keys = append(keys, k.(string))
			}
			for _, key := range keys {
				if _, exist := thisMap[key]; !exist {
					value, err := Map.Get(key)
					if err != nil {
						log.Println(err)
						return nil
					}
					thisMap[key] = value.(string)
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
