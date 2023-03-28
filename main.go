package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Vicenteefenequis/codebank/infrastructure/grpc/server"
	"github.com/Vicenteefenequis/codebank/infrastructure/kafka"
	"github.com/Vicenteefenequis/codebank/infrastructure/repository"
	"github.com/Vicenteefenequis/codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := setupDb()
	defer db.Close()

	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(processTransactionUseCase)

}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))
	return producer
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("db"),
		os.Getenv("host"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("error connection database")
	}

	return db

}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	fmt.Println("Rodando gRPC Server")
	grpcServer.Serve()
}

// cc := domain.NewCreditCard()
// 	cc.Number = "1234"
// 	cc.Name = "Vicente"
// 	cc.ExpirationMonth = 7
// 	cc.ExpirationMonth = 2023
// 	cc.CVV = 123
// 	cc.Limit = 1000
// 	cc.Balance = 0
