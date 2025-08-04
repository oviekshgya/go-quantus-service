package rabbitMQ

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

type RabbitMQImpl struct {
	RabbitConn    *amqp.Connection
	RabbitChannel *amqp.Channel
	MaxRetry      int
}

func NewRabbitMQConnection() *RabbitMQImpl {
	var err error
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", viper.GetString("RABBIT_USER"), viper.GetString("RABBIT_PASSWORD"), viper.GetString("RABBIT_HOST"), viper.GetString("RABBIT_PORT")))
	if err != nil {
		log.Fatal("Gagal menghubungkan ke RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Gagal membuka channel RabbitMQ:", err)
	}
	// Declare main queue with DLX
	args := amqp.Table{
		"x-dead-letter-exchange":    "payment-dlx-exchange",
		"x-dead-letter-routing-key": "dlq-key",
	}

	err = ch.ExchangeDeclare(viper.GetString("RABBIT_EXCHANGE_NAME"), "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Gagal deklarasi exchange utama:", err)
	}

	_, err = ch.QueueDeclare(viper.GetString("RABBIT_QUEUE_NAME"), true, false, false, false, args)
	if err != nil {
		log.Fatal("Gagal deklarasi queue utama:", err)
	}

	err = ch.QueueBind(viper.GetString("RABBIT_QUEUE_NAME"), viper.GetString("RABBIT_ROUTING_KEY"), viper.GetString("RABBIT_EXCHANGE_NAME"), false, nil)
	if err != nil {
		log.Fatal("Gagal bind queue utama:", err)
	}

	log.Println("✅ Terhubung dan queue disiapkan")
	return &RabbitMQImpl{
		RabbitChannel: ch,
		RabbitConn:    conn,
		MaxRetry:      3,
	}
}

func (s *RabbitMQImpl) Publish(exchange, routeKey string, data []byte, retryCount int) error {
	err := s.RabbitChannel.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Headers:     amqp.Table{"x-retry": retryCount},
		},
	)
	if err != nil {
		log.Println("❌ Gagal publish:", err)
	}
	return err
}

func (s *RabbitMQImpl) Consume(queueName string, handler func(body []byte) error) {
	msgs, err := s.RabbitChannel.Consume(
		queueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Gagal memulai consumer RabbitMQ:", err)
	}

	go func() {
		for d := range msgs {
			retry := 0
			if header, ok := d.Headers["x-retry"]; ok {
				retry, _ = strconv.Atoi(fmt.Sprintf("%v", header))
			}

			log.Printf("📥 Terima pesan: %s (retry: %d)", d.Body, retry)

			if err := handler(d.Body); err != nil {
				if retry >= s.MaxRetry {
					log.Println("💀 Retry max, pesan ditolak masuk DLQ")
					d.Nack(false, false) // false: jangan requeue → masuk DLQ
				} else {
					log.Printf("🔁 Re-publish pesan (retry %d)", retry+1)
					s.Publish("payment-exchange", "payment.key", d.Body, retry+1)
					d.Ack(false) // sudah ditangani (akan re-publish manual)
				}
			} else {
				log.Println("✅ Pesan berhasil diproses")
				d.Ack(false)
			}
		}
	}()

	log.Println("🚀 Consumer berjalan...")
	select {}
}
