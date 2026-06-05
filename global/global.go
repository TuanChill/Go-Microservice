package global

import (
	"database/sql"
	"fmt"

	"go_template/configs"
	"go_template/internal/controllers/initialization"
	"go_template/internal/models"
	pkg "go_template/pkg/setting"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	Cfg          models.Config
	DB           *sql.DB
	Cache        *redis.Client
	AdminSdk     *firebase.App
	MessageQueue *amqp.Connection
	SES          *sesv2.Client
	S3           *s3.Client
)

func init() {
	//* CONFIG
	var err error
	Cfg, err = configs.LoadConfig("configs/yaml")
	if err != nil {
		panic(err)
	}

	//* DATABASE
	DB, err = initialization.ConnectPG(Cfg)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		panic(err)
	}

	//* CACHE
	Cache, err = initialization.ConnectRedis(Cfg)
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		panic(err)
	}

	//* FIREBASE
	AdminSdk, err = pkg.InitializeApp()
	if err != nil {
		fmt.Printf("Error connecting to firebase: %v\n", err)
		panic(err)
	}

	//* RABBITMQ
	MessageQueue, err = initialization.ConnectRabbitMQ(Cfg.RabbitMQ.URL)
	if err != nil {
		fmt.Printf("Error connecting to RabbitMq: %v\n", err)
		panic(err)
	}

	//* SES
	SES, err = initialization.ConnectSES(Cfg)
	if err != nil {
		fmt.Printf("Error connecting to AWS SES: %v\n", err)
		panic(err)
	}

	//* S3
	S3, err = initialization.ConnectS3(Cfg)
	if err != nil {
		fmt.Printf("Error connecting to AWS S3: %v\n", err)
		panic(err)
	}
}
