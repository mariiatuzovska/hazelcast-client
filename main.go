package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/core"
	"github.com/urfave/cli"
)

var (
	ServiceName = "hazelcast-examples"
	Version     = "0.0.2"
)

func main() {
	app := cli.NewApp()
	app.Name = ServiceName
	app.Usage = "command line client"
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
			Name:   "pessimistic-lock-example",
			Usage:  "Presents pessimistic-lock-example for map",
			Action: pessimisticLockExample,
		},
		{
			Name:   "optimistic-lock-example",
			Usage:  "Presents optimistic-lock-example for map",
			Action: optimisticLockExample,
		},
		{
			Name:   "task",
			Action: task,
		},
	}
	app.Run(os.Args)
}

var (
	nameLockExampleMap = "lock.example.map"

	capitalsMap = map[string]string{
		"Belarus":        "Minsk",
		"Belgium":        "Brussels",
		"Brazil":         "Brasilia",
		"Bulgaria":       "Sofia",
		"Chile":          "Santiago",
		"China":          "Beijing",
		"Cyprus":         "Nicosia",
		"Egypt":          "Cairo",
		"Estonia":        "Tallinn",
		"Germany":        "Berlin",
		"Greece":         "Athens",
		"Japan":          "Tokyo",
		"Morocco":        "Rabat",
		"Poland":         "Warsaw",
		"Slovakia":       "Bratislava",
		"Spain":          "Madrid",
		"Thailand":       "Bangkok",
		"Ukraine":        "Kyiv",
		"United Kingdom": "London",
	}
)

func mapExample(c *cli.Context) error {

	Alice := make(chan MapWriterChan) // writer will write AlicesWriterChan.Key : AlicesWriterChan.Value into map "Map"
	Bob := make(chan MapWriterChan)   // writer will write BobsWriterChan.Key : BobsWriterChan.Value into map "Map"
	pinger := make(chan bool)

	go MapWriter("Alice", "192.168.0.102:5701", Alice, pinger)
	go MapWriter("Bob", "192.168.0.102:5702", Bob, pinger)
	go MapReader("Carl", "192.168.0.102:5703", pinger)

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
	config.NetworkConfig().AddAddress("192.168.0.102:5701")

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

func pessimisticLockExample(c *cli.Context) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5701")
	config.SetClientName("c1")
	client1, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5702")
	config.SetClientName("c2")
	client2, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5703")
	config.SetClientName("c3")
	client3, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	Map1, _ := client1.GetMap("new.map")
	Map2, _ := client2.GetMap("new.map")
	Map3, _ := client3.GetMap("new.map")
	Map1.PutIfAbsent(int64(1), int64(0))
	times := make(chan bool, 10)

	var Routine = func(Map core.Map, t chan bool) {
		for k := 0; k < 1000; k++ {
			if k%100 == 0 {
				fmt.Println(k)
			}
			Map.Lock(int64(1))
			value, _ := Map.Get(int64(1))
			v := value.(int64)
			v++
			Map.Put(int64(1), v)
			Map.Unlock(int64(1))
		}
		t <- true
	}

	go Routine(Map1, times)
	go Routine(Map2, times)
	go Routine(Map3, times)

	// wait
	<-times
	<-times
	<-times

	val, _ := Map1.Get(int64(1))
	log.Println(fmt.Sprintf("RESULT = %d", val.(int64)))
	Map1.Delete(int64(1))

	client1.Shutdown()
	client2.Shutdown()
	client3.Shutdown()

	return nil
}

func optimisticLockExample(c *cli.Context) error {

	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5701")
	config.SetClientName("c1")
	client1, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5702")
	config.SetClientName("c2")
	client2, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5703")
	config.SetClientName("c3")
	client3, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	Map1, _ := client1.GetMap("new.map")
	Map2, _ := client2.GetMap("new.map")
	Map3, _ := client3.GetMap("new.map")
	Map1.PutIfAbsent(int64(1), int64(0))
	times := make(chan bool, 10)

	var Routine = func(Map core.Map, t chan bool) {
		for k := 0; k < 1000; k++ {
			if k%100 == 0 {
				fmt.Println(k)
			}
			for {
				oldValue, _ := Map.Get(int64(1))
				newValue := oldValue.(int64) + 1
				time.Sleep(time.Duration(10) * time.Millisecond)
				ok, _ := Map.ReplaceIfSame(int64(1), oldValue, newValue)
				if ok {
					break
				}
			}
		}
		t <- true
	}

	go Routine(Map1, times)
	go Routine(Map2, times)
	go Routine(Map3, times)

	// wait
	<-times
	<-times
	<-times

	val, _ := Map1.Get(int64(1))
	log.Println(fmt.Sprintf("RESULT = %d", val.(int64)))
	Map1.Delete(int64(1))

	client1.Shutdown()
	client2.Shutdown()
	client3.Shutdown()

	return nil
}

func queueExample(c *cli.Context) error {

	Alice := make(chan QueueWriterChan) // writer will write Value into queue
	Bob := make(chan QueueWriterChan)
	pinger := make(chan bool)

	go QueueWriter("Alice", "192.168.0.102:5701", Alice, pinger)
	go QueueWriter("Bob", "192.168.0.102:5702", Bob, pinger)
	go QueueReader("Carl", "192.168.0.102:5703", pinger)

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

	go TopicWriter("Alice", "192.168.0.102:5701", topicMessageListener, Alice)
	go TopicWriter("Bob", "192.168.0.102:5702", topicMessageListener, Bob)

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

func task(c *cli.Context) error {
	config := hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5701")
	config.SetClientName("Alice")
	Alice, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5702")
	config.SetClientName("Bob")
	Bob, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = hazelcast.NewConfig()
	config.GroupConfig().SetName("dev")
	config.GroupConfig().SetPassword("dev-pass")
	config.NetworkConfig().AddAddress("192.168.0.102:5703")
	config.SetClientName("Bob")
	Carl, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	AlicesMap, _ := Alice.GetMap("task.map")
	BobsMap, _ := Bob.GetMap("task.map")
	CarlsMap, _ := Carl.GetMap("task.map")
	values, _ := CarlsMap.Values()
	fmt.Println(len(values))
	for i := 0; i < 333; i++ {
		j := i * 3
		AlicesMap.PutIfAbsent(j, j)
		BobsMap.PutIfAbsent(j+1, j+1)
		CarlsMap.PutIfAbsent(j+2, j+2)
	}
	Alice.Shutdown()
	Bob.Shutdown()

	values, _ = CarlsMap.Values()
	fmt.Println(len(values))
	for i := 0; i < 1000; i++ {
		CarlsMap.Delete(i)
	}

	// Alice.Shutdown()
	// Bob.Shutdown()
	Carl.Shutdown()

	return nil
}
