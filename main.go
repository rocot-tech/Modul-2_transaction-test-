package main

import (
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64

	mu sync.Mutex
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (ps *PaymentSystem) AddUser(user *User) {
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(transaction Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, transaction)
}

func (u *User) Deposit(amount float64) {

	u.mu.Lock()
	defer u.mu.Unlock()

	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.Balance < amount {
		return fmt.Errorf("недостаточно средств")
	}

	u.Balance -= amount

	return nil
}

func (ps *PaymentSystem) ProcessingTransactions(transaction Transaction) error {

	fromUser, ok := ps.Users[transaction.FromID]
	if !ok {
		return fmt.Errorf("Отправитель не найден")
	}

	toUser, ok := ps.Users[transaction.ToID]
	if !ok {
		return fmt.Errorf("Получатель не найден")
	}

	err := fromUser.Withdraw(transaction.Amount)
	if err != nil {
		return err
	}

	toUser.Deposit(transaction.Amount)

	return nil
}

func (ps *PaymentSystem) Worker(
	ch <-chan Transaction,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for t := range ch {

		err := ps.ProcessingTransactions(t)

		if err != nil {
			fmt.Println("ошибка транзакции:", err)
		}
	}
}

func main() {

	ps := &PaymentSystem{
		Users:            make(map[string]*User),
		TransactionQueue: []Transaction{},
	}

	fmt.Println("Создаю UserID: 1 с балансом 1000")

	user1 := &User{
		ID:      "1",
		Name:    "User1",
		Balance: 1000,
	}

	fmt.Println("Создаю UserID: 2 с балансом 500")

	user2 := &User{
		ID:      "2",
		Name:    "User2",
		Balance: 500,
	}

	ps.AddUser(user1)
	ps.AddUser(user2)

	transaction1 := Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 200,
	}

	transaction2 := Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 50,
	}

	ps.AddTransaction(transaction1)
	ps.AddTransaction(transaction2)

	ch := make(chan Transaction, len(ps.TransactionQueue))

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {

		wg.Add(1)

		go ps.Worker(ch, &wg)
	}

	for _, t := range ps.TransactionQueue {
		ch <- t
	}

	close(ch)

	wg.Wait()

	fmt.Println("Итого:")
	fmt.Printf("У первого пользователя получилось %.0f\n", ps.Users["1"].Balance)
	fmt.Printf("У второго пользователя получилось %.0f\n", ps.Users["2"].Balance)
}
