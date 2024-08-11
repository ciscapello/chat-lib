package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient структура для работы с RabbitMQ
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQClient создает нового клиента RabbitMQ
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// DeclareQueue объявляет очередь, если её ещё нет
func (client *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	queue, err := client.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return queue, err
}

// Publish отправляет сообщение в указанную очередь
func (client *RabbitMQClient) Publish(queueName string, message []byte) error {

	err := client.channel.Publish(
		"",        // exchange
		queueName, // routing key (название очереди)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return err
}

// ConsumeMessages начинает потребление сообщений из указанной очереди
func (client *RabbitMQClient) ConsumeMessages(queueName string, handler func([]byte) error) error {
	queue, err := client.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	msgs, err := client.channel.Consume(
		queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	log.Printf("Waiting for messages from queue: %s. To exit press CTRL+C", queueName)
	<-forever

	return nil
}

// Close закрывает соединение и канал
func (client *RabbitMQClient) Close() {
	client.channel.Close()
	client.conn.Close()
}
