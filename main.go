package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var streamName = "DEMOSTREAM"
var consumerName = "DEMOSTREAM_CONSUMER"

func main() {
	admin := getJs("./admin.seed")
	publisher := getJs("./publisher.seed")
	subscriber := getJs("./subscriber.seed")

	err := admin.DeleteStream(streamName)

	if err != nil {
		fmt.Println(err)
	}
	_, err = admin.AddStream(&nats.StreamConfig{Name: streamName, MaxMsgSize: 1024 * 1024, Subjects: []string{"demostream.>"}, MaxMsgs: 10, MaxBytes: (1024 * 10), Storage: nats.MemoryStorage})

	if err != nil {
		panic(err)
	}

	_, err = admin.AddConsumer(streamName, &nats.ConsumerConfig{
		Name:          consumerName,
		Durable:       consumerName,
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: "demostream" + ".>",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println()

	data, err := GenerateRandomBytes(1024 * 1024)
	if err != nil {
		panic(err)
	}
	noOfMsgsToSend := 100

	//subscribe
	sub, err := subscriber.PullSubscribe("demostream.>", consumerName, nats.BindStream(streamName))
	if err != nil {
		panic(err)
	}
	msgsRcd := 0
	go func() {
		for {
			fmt.Printf("trying to get msg %d\n", msgsRcd)
			msgs, err := sub.Fetch(1)
			if err != nil {
				fmt.Println(err)
				fmt.Println("no message still waiting")
				time.Sleep(time.Second)
				continue
			}
			for _, msg := range msgs {
				fmt.Printf("got msg %d\n", msgsRcd)
				msg.Ack()
			}
			if err != nil {
				fmt.Println(msgsRcd)
				panic(err)
			}
			msgsRcd++
		}
	}()
	//publish
	for i := 0; i < noOfMsgsToSend; i++ {
		ack, err := publisher.Publish(fmt.Sprintf("demostream.%d", i), data)
		if err != nil {
			fmt.Println(i)
			panic(err)
		}

		fmt.Printf("pubd msg %d on stream %s\n", i, ack.Stream)
		time.Sleep(time.Second * 3)
	}
	if noOfMsgsToSend != msgsRcd {
		panic(fmt.Sprintf("no of msg rcd %d", msgsRcd))
	}
}

func getJs(seed string) nats.JetStreamContext {
	nkuo, err := nats.NkeyOptionFromSeed(seed)
	if err != nil {
		panic(err)
	}
	nc, err := nats.Connect(nats.DefaultURL, nkuo)
	if err != nil {
		panic(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	return js
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
