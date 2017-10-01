package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	insight "github.com/apganesh/Insight_Project/Spotlite/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
)

var (
	driverConsumer *kafka.Consumer
	rpcClient      insight.MatcherClient
	sigChan        chan bool
)

const (
	address = "localhost:50051"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func checkWarning(err error) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return false
	}
	return true
}

// RPC call for adding Driver to the database
func addDriver(client insight.MatcherClient, driver *insight.Driver) {
	_, err := client.AddDriver(context.Background(), driver)
	if err != nil {
		log.Fatalf("Could not add Driver: %v", err)
	}
}

// RPC Call for Matching Driver
func getDriver(client insight.MatcherClient, rider *insight.Rider) {
	resp, err := client.GetDriver(context.Background(), rider)
	if err != nil {
		log.Fatalf("Could not get Driver: %v", err)
	}
	fmt.Println("Got driver: ", resp.Id, resp.Lat, resp.Lng)
}

func setupKafkaConsumer(broker string, group string) (*kafka.Consumer, error) {
	var err error
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	return consumer, err
}

func driverCallback(buf []byte) {
	driver := new(insight.Driver)
	proto.Unmarshal(buf, driver)
	//fmt.Println(driver.GetId(), driver.GetLat(), driver.GetLng(), driver.GetRadius(), driver.GetStatus())
	addDriver(rpcClient, driver)
}

func riderCallback(buf []byte) {
	rider := new(insight.Rider)
	proto.Unmarshal(buf, rider)
	fmt.Println(rider.Id, rider.SLat, rider.SLng)
	getDriver(rpcClient, rider)
}

func handleMessage(msg []byte) {
	if msg[0] == 0 {
		driverCallback(msg[1:])
	} else {
		riderCallback(msg[1:])
	}
}
func runConsumer(consumer *kafka.Consumer) {
	run := true

	for run == true {
		select {
		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Unassign()
			case *kafka.Message:
				//fmt.Println("Got message: ", string(e.Value))
				handleMessage(e.Value)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v %v\n", e, time.Now().Unix())
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}
	}
}

func main() {

	broker := insight.DefaultConfig.KafkaBrokers
	driverTopics := insight.DefaultConfig.KafkaTopics
	groupId := insight.DefaultConfig.KafkaGroup

	var err error
	sigChan = make(chan bool)

	driverConsumer, err = setupKafkaConsumer(broker, groupId)
	checkError(err)
	fmt.Println(driverConsumer.String())

	// riderTopics := []string{"r1"}

	err = driverConsumer.SubscribeTopics(driverTopics, nil)
	checkError(err)

	fmt.Printf("Created Consumer: %v \n", driverConsumer.String())

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to the microService: %v", err)
	}
	defer conn.Close()

	// Creates a new CustomerClient
	rpcClient = insight.NewMatcherClient(conn)

	runConsumer(driverConsumer)

	fmt.Printf("Closing consumer\n")
	<-sigChan
	driverConsumer.Close()
}
