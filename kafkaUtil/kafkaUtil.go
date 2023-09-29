// this utility class is used to write log messages from logHarbour to Kafka
package logHarbour

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var KafkaProducer *kafka.Producer
var topic = "logHarbour"
var logId []byte

// initialize Kafka Producer with a log ID which will be a key to all messages written to Kafka.
// This key should ideally be combination of appName+moduleName+systemName
func KafkaInit(id string) {
	logId = []byte(id)
	kafkaStart()
}

// start kafka producer.
func kafkaStart() {
	KafkaProducer, _ = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "0.0.0.0:9092"})
	//fmt.Println("->> New Kafka Producer <<-")
	go func() {
		//fmt.Println("->> New Kafka Event Listener <<-")
		for e := range KafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()
	KafkaProducer.Flush(100)
}

func sendMsgToKafka(msg []byte) {

	if KafkaProducer == nil {
		kafkaStart()
	}
	//deliveryChan := make(chan kafka.Event)

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
		Key:            logId,
	}

	// produce the message
	KafkaProducer.Produce(message, nil)
	//defer KafkaProducer.Close()
}

type KafkaWriter struct {
}

func (e KafkaWriter) Write(msg []byte) (int, error) {
	sendMsgToKafka(msg)
	return len(msg), nil
}
