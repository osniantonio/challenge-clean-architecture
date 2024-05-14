package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graph_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/osniantonio/challenge-clean-architecture/configs"
	"github.com/osniantonio/challenge-clean-architecture/internal/event/handler"
	"github.com/osniantonio/challenge-clean-architecture/internal/infra/graph"
	"github.com/osniantonio/challenge-clean-architecture/internal/infra/grpc/pb"
	"github.com/osniantonio/challenge-clean-architecture/internal/infra/grpc/service"
	"github.com/osniantonio/challenge-clean-architecture/internal/infra/web/webserver"
	"github.com/osniantonio/challenge-clean-architecture/internal/usecase"
	"github.com/osniantonio/challenge-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	db, err := openDatabase(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ch, err := openChannel(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	dispatcher := events.NewEventDispatcher()
	dispatcher.Register("OrderCreated", handler.NewOrderCreatedHandler(ch))

	startWebServer(config, db, dispatcher)

	createOrderUseCase := NewCreateOrderUseCase(db, dispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	if err := startGRPCServer(config, createOrderUseCase, listOrdersUseCase); err != nil {
		log.Fatalln(err.Error())
	}

	startGraphQLServer(config, *createOrderUseCase, *listOrdersUseCase)
}

func loadConfig() (*configs.Config, error) {
	config, err := configs.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func openDatabase(config *configs.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open(config.DBDriver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openChannel(config *configs.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(config.AMQPUrl)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func startWebServer(config *configs.Config, db *sql.DB, dispatcher *events.EventDispatcher) {
	server := webserver.NewWebServer(config.WebServerPort)
	handler := NewWebOrderHandler(db, dispatcher)
	server.Router.Get("/order", handler.List)
	server.Router.Post("/order", handler.Create)
	fmt.Printf("Starting web server on port %s\n", config.WebServerPort)
	go server.Start()
}

func startGRPCServer(config *configs.Config, createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) error {
	server := grpc.NewServer()
	service := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(server, service)
	reflection.Register(server)
	fmt.Printf("Starting gRPC server on port %s\n", config.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))
	if err != nil {
		return err
	}
	go server.Serve(lis)
	return nil
}

func startGraphQLServer(config *configs.Config, createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) {
	server := graph_handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					CreateOrderUseCase: createOrderUseCase,
					ListOrdersUseCase:  listOrdersUseCase,
				},
			},
		),
	)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)
	fmt.Printf("Starting GraphQL server on port %s\n", config.GraphQLServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.GraphQLServerPort), nil)
}
