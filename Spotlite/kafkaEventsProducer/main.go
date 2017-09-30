package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	insight "github.com/apganesh/Insight_Project/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogo/protobuf/proto"
)

var (
	producer *kafka.Producer
	doneChan chan bool
	topic    string
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

func init() {
	doneChan = make(chan bool)
}

func handleRiderMessage(rec []string) {
	//fmt.Println("The record is: ", rec)
	rd := new(insight.Rider)
	rd.Id, _ = strconv.ParseInt(rec[1], 10, 64)
	rd.SLat, _ = strconv.ParseFloat(rec[2], 64)
	rd.SLng, _ = strconv.ParseFloat(rec[3], 64)
	rd.SLat, _ = strconv.ParseFloat(rec[4], 64)
	rd.SLng, _ = strconv.ParseFloat(rec[5], 64)
	rd.Timestamp, _ = strconv.ParseInt(rec[6], 10, 64)

	msg, err := proto.Marshal(rd)
	if err != nil {
		log.Fatal(err)
		return
	}
	one := make([]byte, 1)

	one[0] = 0x1
	msg = append(one, msg...)
	sendKafkaMessage(string(msg))
}

func handleDriverMessage(rec []string) {
	//fmt.Println("The record is: ", rec)
	dd := new(insight.Driver)
	dd.Id, _ = strconv.ParseInt(rec[1], 10, 64)
	dd.Lat, _ = strconv.ParseFloat(rec[2], 64)
	dd.Lng, _ = strconv.ParseFloat(rec[3], 64)
	dd.Radius, _ = strconv.ParseFloat(rec[4], 64)
	dd.Timestamp, _ = strconv.ParseInt(rec[5], 10, 64)
	dd.Status, _ = strconv.ParseInt(rec[6], 10, 64)
	msg, err := proto.Marshal(dd)
	if err != nil {
		log.Fatal(err)
		return
	}
	one := make([]byte, 1)

	one[0] = 0x0
	msg = append(one, msg...)
	sendKafkaMessage(string(msg))
}

func handleMessage(rec []string) {
	if rec[0] == "1" {
		handleDriverMessage(rec)
	} else {
		handleRiderMessage(rec)
	}
}

func readDrivers(fname string) error {
	driverFile, err := os.Open(fname)
	if err != nil {
		fmt.Println("Cannot open file at location: ", fname)
		return err
	}
	defer driverFile.Close()

	r := csv.NewReader(driverFile)
	r.LazyQuotes = true

	for {
		if rec, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		} else {
			handleMessage(rec)
		}
	}
	fmt.Println("Done reading the data")
	return nil
}

func readDriverFile() error {
	fname := "../data/drivers.log"
	driverFile, err := os.Open(fname)
	if err != nil {
		fmt.Println("Cannot open file at location: ", fname)
		return err
	}
	defer driverFile.Close()

	r := csv.NewReader(driverFile)
	r.LazyQuotes = true

	ticker := time.Tick(1 * time.Second)

	totalcount := 0
	// read 5 lines and send it out
	for {
		select {
		case <-ticker:
			count := 0
			for {
				if count == 1000 {
					fmt.Println("Finished reading drivers: ", totalcount)
					break
				}
				if rec, err := r.Read(); err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
					return err
				} else {
					count++
					totalcount++
					handleDriverMessage(rec)
				}
			}

		}
	}
	return nil
}

func readRiderFile() error {
	fname := "../data/riders.log"
	riderFile, err := os.Open(fname)
	if err != nil {
		fmt.Println("Cannot open file at location: ", fname)
		return err
	}
	defer riderFile.Close()

	r := csv.NewReader(riderFile)
	r.LazyQuotes = true

	ticker := time.Tick(1 * time.Second)
	totalcount := 0
	// read 5 lines and send it out
	for {
		select {
		case <-ticker:
			count := 0
			for {
				if count == 10 {
					fmt.Println("Finished reading riders: ", totalcount)
					break
				}
				if rec, err := r.Read(); err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
					return err
				} else {
					count++
					totalcount++
					handleRiderMessage(rec)
				}
			}

		}
	}
	return nil
}

func sendKafkaMessage(mesg string) {
	producer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(mesg)}
}

func setupKafka(broker, topic string) error {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return err
	}

	fmt.Printf("Created Producer %v\n", producer.String())

	go func() {
		for e := range producer.Events() {
			//fmt.Println("Producer got an event ... ")
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
	return err

}

func main() {
	/*
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n",
				os.Args[0])
			os.Exit(1)
		}
	*/

	//broker := os.Args[1]
	topic = insight.DefaultConfig.KafkaTopics
	broker := insight.DefaultConfig.KafkaBrokers
	//fmt.Println("broker and topic ", broker, topic)

	err := setupKafka(broker, topic)
	if err != nil {
		fmt.Printf("Failed to setup Kafka: %s\n", err)
		return
	}

	//readDrivers("../data/driver_events.log")
	go readDriverFile()
	time.Sleep(2 * time.Second)
	go readRiderFile()
	_ = <-doneChan
	producer.Close()
}
