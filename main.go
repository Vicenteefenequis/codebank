package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Vicenteefenequis/codebank/domain"
	"github.com/Vicenteefenequis/codebank/infrastructure/repository"
	"github.com/Vicenteefenequis/codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Vicente"
	cc.ExpirationMonth = 7
	cc.ExpirationMonth = 2023
	cc.CVV = 123
	cc.Limit = 1000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Print(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	return usecase.NewUseCaseTransaction(transactionRepository)
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("error connection database")
	}

	return db

}
