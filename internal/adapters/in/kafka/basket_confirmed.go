package kafka

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/generated/queues/basketconfirmedpb"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
)

type BasketConfirmedConsumer interface {
	Consume() error
	Close() error
}

var _ BasketConfirmedConsumer = &basketConfirmedConsumer{}

type basketConfirmedConsumer struct {
	topic                     string
	consumerGroup             sarama.ConsumerGroup
	createOrderCommandHandler commands.CreateOrderCommandHandler
	ctx                       context.Context
	cancel                    context.CancelFunc
}

func NewBasketConfirmedConsumer(brokers []string, group string, topic string,
	createOrderCommandHandler commands.CreateOrderCommandHandler) (BasketConfirmedConsumer, error) {
	if brokers == nil || len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}
	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &basketConfirmedConsumer{
		topic:                     topic,
		consumerGroup:             consumerGroup,
		createOrderCommandHandler: createOrderCommandHandler,
		ctx:                       ctx,
		cancel:                    cancel,
	}, nil
}

func (c *basketConfirmedConsumer) Close() error {
	c.cancel()
	return c.consumerGroup.Close()
}

func (c *basketConfirmedConsumer) Consume() error {
	handler := &consumerGroupHandler{
		createOrderCommandHandler: c.createOrderCommandHandler,
	}

	for {
		err := c.consumerGroup.Consume(c.ctx, []string{c.topic}, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
			return err
		}
		if c.ctx.Err() != nil {
			return nil
		}
	}
}

type consumerGroupHandler struct {
	createOrderCommandHandler commands.CreateOrderCommandHandler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := context.Background()
		fmt.Printf("Received: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		var event basketconfirmedpb.BasketConfirmedIntegrationEvent
		err := json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			session.MarkMessage(message, "") // Всё равно отметить как прочитанное
			continue
		}

		createOrderCommand, err := commands.NewCreateOrderCmd(
			uuid.MustParse(event.BasketId), event.Address.Street, int(event.Volume))
		if err != nil {
			log.Printf("Failed to create createOrder command: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		err = h.createOrderCommandHandler.Handle(ctx, createOrderCommand)
		if err != nil {
			log.Printf("Failed to handle createOrder command: %v", err)
		}

		// После успешной обработки сообщения — отметить его
		session.MarkMessage(message, "")
	}

	return nil
}
