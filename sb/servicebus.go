package sb

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

//Producer is used to send message
func Producer(connectionString, topic string, interval, timeout time.Duration, bytes int, doneCh chan os.Signal) {
	t := getTopic(connectionString, topic)

	for {
		select {
		case <-doneCh:
			return
		case <-time.After(interval):
			log.Println("sending message")
			deadline := time.Now().Add(timeout)
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()

			s, err := t.NewSender(ctx)
			if err != nil {
				log.Println("unable to generate sender ", err)
				break
			}

			msg := servicebus.NewMessageFromString(generateData(bytes))
			if err := s.Send(ctx, msg); err != nil {
				log.Println("error sending message:  ", err)
			}

			log.Println("message sent")
		}
	}
}

func Consumer(connectionString, topic, subscription string, timeout time.Duration, doneCh chan os.Signal) {
	t := getTopic(connectionString, topic)
	client, err := t.NewSubscription(subscription)
	if err != nil {
		log.Fatal("unable to generate subscription:  ", err)
	}

	for {
		select {
		case <-doneCh:
			return
		default:
			expire := time.Now().Add(timeout)
			ctx, cancel := context.WithDeadline(context.Background(), expire)
			defer cancel()
			var rx servicebus.HandlerFunc = msgHandler
			if err := client.ReceiveOne(ctx, rx); err != nil {
				log.Println("error receiveing message:  ", err)
			}
		}
	}
}

func msgHandler(ctx context.Context, msg *servicebus.Message) error {
	log.Println("message received")
	log.Printf("msg:  %+v", msg)
	return msg.Complete(ctx)
}

func getTopic(connectionString, topic string) *servicebus.Topic {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		log.Fatal("unable to generate service bus namespace:  ", err)
	}

	t, err := ns.NewTopic(topic)
	if err != nil {
		log.Fatal("unable to generate service bus topic:  ", err)
	}

	return t
}

//the below is from stack over flow (link below)
//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func generateData(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
