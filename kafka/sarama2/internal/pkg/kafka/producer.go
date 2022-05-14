package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"experiment_go/kafka/sarama1/internal/model"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
	SyncProducer  sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
}

func NewProducer(cfg model.SaramaConfig) *KafkaProducer {
	tlsConfig := createTlsConfiguration(cfg)
	producer := KafkaProducer{
		SyncProducer:  initSyncProducer(cfg, tlsConfig),
		AsyncProducer: initAsyncProducer(cfg, tlsConfig),
	}
	return &producer
}

func (kp *KafkaProducer) Close() error {
	if err := kp.SyncProducer.Close(); err != nil {
		log.Println("[KAFKA] Failed to shutdown sync producer cleanly: ", err)
	}
	if err := kp.AsyncProducer.Close(); err != nil {
		log.Println("[KAFKA] Failed to shutdown async producer cleanly: ", err)
	}
	return nil
}

func createTlsConfiguration(cfg model.SaramaConfig) (t *tls.Config) {
	if cfg.CertFile != "" && cfg.KeyFile != "" && cfg.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			log.Fatal("[KAFKA] Failed to load cert: ", err)
		}

		caCert, err := os.ReadFile(cfg.CaFile)
		if err != nil {
			log.Fatal("[KAFKA] Failed to load ca file: ", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: cfg.VerifySSL,
		}
	}
	return t
}

func initSyncProducer(cfg model.SaramaConfig, tlsConfig *tls.Config) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	producer, err := sarama.NewSyncProducer(cfg.BrokersHost, config)
	if err != nil {
		log.Fatalln("[KAFKA] Failed to start Sarama Sync producer: ", err)
	}
	return producer
}

func initAsyncProducer(cfg model.SaramaConfig, tlsConfig *tls.Config) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	producer, err := sarama.NewAsyncProducer(cfg.BrokersHost, config)
	if err != nil {
		log.Fatalln("[KAFKA] Failed to start Sarama Async producer: ", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("[KAFKA] Failed to write message async: ", err)
		}
	}()

	return producer
}

// SendMessage

// close producer
