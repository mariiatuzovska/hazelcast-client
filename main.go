package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = ServiceName
	app.Usage = fmt.Sprintf("%s command line client", ServiceName)
	app.Description = ""
	app.Version = Version
	app.Copyright = "2020, mariiatuzovska"
	app.Authors = []cli.Author{cli.Author{Name: "Mariia Tuzovska"}}
	app.Commands = []cli.Command{
		{
			Name:   "map-example",
			Usage:  "Presents map-example",
			Action: mapExample,
		},
		{
			Name:   "queue-example",
			Usage:  "Presents queue-example",
			Action: queueExample,
		},
		{
			Name:   "topic-example",
			Usage:  "Presents topic-example",
			Action: topicExample,
		},
		{
			Name:   "lock-example",
			Usage:  "Presents lock-example",
			Action: lockExample,
		},
	}
	app.Run(os.Args)
}

var (
	ServiceName = "hazelcast-go-client"
	Version     = "0.0.1"

	capitalsMap = map[string]string{
		"Belarus":  "Minsk",
		"Belgium":  "Brussels",
		"Brazil":   "Brasilia",
		"Bulgaria": "Sofia",
		"Chile":    "Santiago",
		"China":    "Beijing",
		"Cyprus":   "Nicosia",
		"Egypt":    "Cairo",
		"Estonia":  "Tallinn",
		"Germany":  "Berlin",
		"Greece":   "Athens",
		"Japan":    "Tokyo",
		"Morocco":  "Rabat",
		"Poland":   "Warsaw",
		"Slovakia": "Bratislava",
		"Spain":    "Madrid",
		"Thailand": "Bangkok",
		"Ukraine":  "Kyiv",
		"United Kingdom	": "London",
	}
)

func mapExample(c *cli.Context) error {

	Alice := make(chan MapWriterChan) // writer will write AlicesWriterChan.Key : AlicesWriterChan.Value into map "Map"
	Bob := make(chan MapWriterChan)   // writer will write BobsWriterChan.Key : BobsWriterChan.Value into map "Map"
	pinger := make(chan bool)

	go MapWriter("Alice", "192.168.0.100:5701", Alice, pinger)
	go MapWriter("Bob", "192.168.0.100:5702", Bob, pinger)
	go MapReader("Carl", "192.168.0.100:5703", pinger)

	i := 0
	for key, value := range capitalsMap {
		time.Sleep(time.Duration(5) * time.Second)
		if i&1 == 0 {
			Alice <- MapWriterChan{key, value}
		} else {
			Bob <- MapWriterChan{key, value}
		}
		i++
	}

	close(Alice)
	close(Bob)
	pinger <- false

	time.Sleep(time.Duration(2) * time.Second)

	close(pinger)

	// delete all entries

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.100:5701")

	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	Map, err := client.GetMap(nameExampleMap)
	if err != nil {
		log.Println(err)
		return nil
	}
	for key := range capitalsMap {
		err = Map.Delete(key)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	client.Shutdown()

	return nil
}

func queueExample(c *cli.Context) error {

	Alice := make(chan QueueWriterChan) // writer will write Value into queue
	Bob := make(chan QueueWriterChan)
	pinger := make(chan bool)

	go QueueWriter("Alice", "192.168.0.100:5701", Alice, pinger)
	go QueueWriter("Bob", "192.168.0.100:5702", Bob, pinger)
	go QueueReader("Carl", "192.168.0.100:5703", pinger)

	i := 0
	for _, value := range capitalsMap {
		time.Sleep(time.Duration(5) * time.Second)
		if i&1 == 0 {
			Alice <- QueueWriterChan{value}
		} else {
			Bob <- QueueWriterChan{value}
		}
		i++
	}

	close(Alice)
	close(Bob)
	pinger <- false
	time.Sleep(time.Duration(2) * time.Second)
	close(pinger)

	return nil
}

func topicExample(c *cli.Context) error {

	var wg sync.WaitGroup
	wg.Add(100)
	topicMessageListener := &TopicMessageListener{&wg}
	Alice := make(chan TopicWriterChan)
	Bob := make(chan TopicWriterChan)

	go TopicWriter("Alice", "192.168.0.100:5701", topicMessageListener, Alice)
	go TopicWriter("Bob", "192.168.0.100:5702", topicMessageListener, Bob)

	i := 0
	for _, value := range capitalsMap {
		time.Sleep(time.Duration(5) * time.Second)
		if i&1 == 0 {
			Alice <- TopicWriterChan{value}
		} else {
			Bob <- TopicWriterChan{value}
		}
		i++
	}

	return nil
}

func lockExample(c *cli.Context) error {
	return nil
}
