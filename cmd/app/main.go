package main

import (
	"delivery/cmd"
	httpin "delivery/internal/adapters/in/http"
	"delivery/internal/generated/servers"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
)

func main() {
	config := getConfig()

	cr := cmd.NewCompositionRoot(config)
	defer cr.CloseAll()

	runCronJobs(cr)
	startKafkaConsumer(cr)
	startWebServer(cr, config.HttpPort)
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

func runCronJobs(cr *cmd.CompositionRoot) {
	c := cron.New()

	_, err := c.AddJob("@every 1s", cr.NewAssignOrdersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}

	_, err = c.AddJob("@every 1s", cr.NewMoveCouriersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}

	c.Start()
}

func startWebServer(cr *cmd.CompositionRoot, port string) {

	handlers, err := httpin.NewServerHandlers(
		cr.NewAllCouriersQueryHandler(),
		cr.NewIncompleteOrdersQueryHandler(),
		cr.NewCreateCourierCommandHandler(),
		cr.NewCreateOrderCommandHandler(),
	)

	if err != nil {
		log.Fatalf("Failed to create http handlers: %v", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
	// 		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
	// 		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
	// 		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\r\n",
	// 	CustomTimeFormat: "2006-01-02 15:04:05.00000",
	// }))

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	registerSwaggerOpenApi(e)
	registerSwaggerUi(e)
	servers.RegisterHandlers(e, servers.NewStrictHandler(handlers, []servers.StrictMiddlewareFunc{}))

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUi(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func startKafkaConsumer(cr *cmd.CompositionRoot) {
	go func() {
		if err := cr.NewBasketConfirmedEventsConsumer().Consume(); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()
}
