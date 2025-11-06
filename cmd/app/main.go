package main

import (
	"delivery/cmd"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

func main() {
	config := getConfig()

	compositionRoot := cmd.NewCompositionRoot(config)
	defer compositionRoot.CloseAll()

	startWebServer(compositionRoot, config.HttpPort)
}

func getConfig() cmd.Config {
	_ = godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	config := cmd.Config{
		HttpPort:                  os.Getenv("HTTP_PORT"),
		DbHost:                    os.Getenv("DB_HOST"),
		DbPort:                    os.Getenv("DB_PORT"),
		DbUser:                    os.Getenv("DB_USER"),
		DbPassword:                os.Getenv("DB_PASSWORD"),
		DbName:                    os.Getenv("DB_NAME"),
		DbSslMode:                 os.Getenv("DB_SSLMODE"),
		GeoServiceGrpcHost:        os.Getenv("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:                 os.Getenv("KAFKA_HOST"),
		KafkaConsumerGroup:        os.Getenv("KAFKA_CONSUMER_GROUP"),
		KafkaBasketConfirmedTopic: os.Getenv("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrderChangedTopic:    os.Getenv("KAFKA_ORDER_CHANGED_TOPIC"),
	}

	return config
}

func startWebServer(_ *cmd.CompositionRoot, port string) {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\r\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}
