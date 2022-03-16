package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/copier"
	"stream-bullet/job"
	"stream-bullet/server"
)

func consumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"bullets", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			var data server.Data
			json.Unmarshal(msg.Value, &data)
			bullet := &job.Bullet{
				Ip:      data.Ip,
				User:    data.User,
				Id:      data.Id,
				Content: data.Content,
				Start:   data.Start,
			}
			copier.Copy(&bullet, data)
			job.Insert(bullet)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
