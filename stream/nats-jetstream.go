package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PeerIslands/aci-fx-go/model/dto/request"
	"github.com/PeerIslands/aci-fx-go/model/dto/response"
	"github.com/PeerIslands/aci-fx-go/model/entity"
	"github.com/PeerIslands/aci-fx-go/service/bal"
	"github.com/PeerIslands/aci-fx-go/service/dal"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dbService = &dal.MongoDbService[entity.ForexData]{}

var fxService = &bal.Fx_service{
	DbService: dbService,
}

func Connect() {
	nc, cerr := nats.Connect(os.Getenv("NATS_URI"))

	if cerr != nil {
		log.Fatal("cerr")
	}

	js, jerr := jetstream.New(nc)

	if jerr != nil {
		log.Fatal("jerr")
	}
	stream, serr := js.Stream(context.Background(), os.Getenv("STREAM"))

	if serr != nil {
		log.Fatal("serr")
	}

	cons, _ := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:       os.Getenv("CONSUMER"),
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: os.Getenv("LISTEN_SUBJECT"),
		MaxWaiting:    0,
	})

	cc, err := cons.Consume(func(msg jetstream.Msg) {
		go func() {
			handle(msg, js)
		}()
	})
	if err != nil {
		print("error")
	}
	defer cc.Stop()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	fmt.Println("Shutting down gracefully...")
}

func handle(msg jetstream.Msg, js jetstream.JetStream) {
	var receivedTime = time.Now().UnixMilli()
	hn, _ := os.Hostname()
	var message request.NatConvertRequest
	unmarshalErr := json.Unmarshal(msg.Data(), &message)
	if unmarshalErr != nil {
		return
	}
	/*conversionResponse := fxService.GetConvertedRate(message.TenantID, message.BankID, message.Amount, message.BaseCurrency, message.TargetCurrency, message.Tier)
	if conversionResponse.Status == "Success" {
		conversionResponse.Data.InitiatedOn = message.InitiatedOn
		conversionResponse.Data.TimeTaken = time.Now().UnixMilli() - message.InitiatedOn
		conversionResponse.Data.ReceivedTime = receivedTime - message.InitiatedOn
	}*/

	conversionResponse := response.ConversionResponse{
		Amount:          0,
		ConvertedAmount: 0,
		BaseCurrency:    "",
		TargetCurrency:  "",
		InitiatedOn:     message.InitiatedOn,
		TimeTaken:       time.Now().UnixMilli() - message.InitiatedOn,
		ReceivedTime:    receivedTime - message.InitiatedOn,
		Rate:            0,
		HostName:        hn,
	}
	processedData, _ := json.Marshal(conversionResponse)
	_, err := js.PublishAsync(os.Getenv("PUBLISH_SUBJECT"), processedData)
	if err != nil {
		return
	}

	ackErr := msg.Ack()
	log.Println("completed")
	if ackErr != nil {
		return
	}

}
