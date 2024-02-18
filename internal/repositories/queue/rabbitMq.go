package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitMQRepository struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     amqp.Queue
	closeCh   chan *amqp.Error
	url       string
	queueName string
}

func NewRabbitMQRepository(url, queueName string) (*RabbitMQRepository, error) {
	repo := &RabbitMQRepository{
		url: url,
	}
	repo.queueName = queueName
	err := repo.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	repo.NotifyClose()
	return repo, nil
}

func (r *RabbitMQRepository) Connect() error {
	var err error
	r.conn, err = amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close() // Закрываем соединение, если канал не открылся
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	r.queue, err = r.channel.QueueDeclare(
		r.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		r.channel.Close() // Закрываем канал, если очередь не объявлена
		r.conn.Close()
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	return nil
}

func (r *RabbitMQRepository) Close() error {
	if r.channel != nil {
		err := r.channel.Close()
		if err != nil {
			return err
		}
	}
	if r.conn != nil {
		err := r.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQRepository) Publish(task []byte) error {
	select {
	case err := <-r.closeCh:
		return fmt.Errorf("publish failed: %w", err)
	default:
		err := r.channel.Publish(
			"",           // exchange
			r.queue.Name, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        task,
			})

		return err
	}
}

func (r *RabbitMQRepository) Consume() (<-chan []byte, error) {
	select {
	case err := <-r.closeCh:
		return nil, err
	default:
		msgs, err := r.channel.Consume(
			r.queue.Name, // queue
			"",           // consumer
			true,         // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			return nil, err
		}

		taskChan := make(chan []byte)
		go func() {
			for msg := range msgs {
				taskChan <- msg.Body
			}
		}()

		return taskChan, nil
	}
}

// ReconnectDelay - начальная задержка перед попыткой переподключения.
const ReconnectDelay = 5 * time.Second

// MaxReconnectDelay - максимальная задержка перед попыткой переподключения.
const MaxReconnectDelay = 30 * time.Second

// ReconnectBackoff - коэффициент увеличения задержки.
const ReconnectBackoff = 2

// NotifyClose - метод для прослушивания событий закрытия соединения и канала.
func (r *RabbitMQRepository) NotifyClose() {
	go func() {
		for {
			select {
			case err := <-r.conn.NotifyClose(make(chan *amqp.Error)):
				r.handleClose(err)
			case err := <-r.channel.NotifyClose(make(chan *amqp.Error)):
				r.handleClose(err)
			}
		}
	}()
}

// handleClose - метод для обработки событий закрытия соединения и канала.
func (r *RabbitMQRepository) handleClose(err *amqp.Error) {
	if err != nil {
		fmt.Printf("Connection closed: %v\n", err)
	}
	r.Reconnect()
}

// Reconnect - метод для переподключения к RabbitMQ.
func (r *RabbitMQRepository) Reconnect() {
	delay := ReconnectDelay
	for {
		err := r.Connect()
		if err == nil {
			log.Printf("Successfully reconnect")

			return
		}
		log.Printf("Failed to reconnect: %v, retrying in %v\n", err, delay)
		time.Sleep(delay)
		delay *= ReconnectBackoff
		if delay > MaxReconnectDelay {
			delay = MaxReconnectDelay
		}
	}
}
