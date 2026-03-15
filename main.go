package main

import (
	"fmt"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {

	if u.Balance < amount {
		return fmt.Errorf("недостаточно средств")
	}

	u.Balance -= amount
	return nil
}

func main() {

	user1 := &User{
		ID:      "1",
		Name:    "User1",
		Balance: 1000,
	}

	user2 := &User{
		ID:      "2",
		Name:    "User2",
		Balance: 500,
	}

	fmt.Println("Начальные балансы:")
	fmt.Println(user1.Name, ":", user1.Balance)
	fmt.Println(user2.Name, ":", user2.Balance)

	fmt.Println("\n Пополнение баланса:")

	user1.Deposit(200)
	fmt.Println(user1.Name, "после пополнения:", user1.Balance)

	user2.Deposit(100)
	fmt.Println(user2.Name, "после пополнения:", user2.Balance)

	fmt.Println("\n Снятие средств:")

	err := user1.Withdraw(300)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(user1.Name, "после снятия:", user1.Balance)
	}

	err = user2.Withdraw(700)
	if err != nil {
		fmt.Println("Ошибка у", user2.Name+":", err)
	} else {
		fmt.Println(user2.Name, "после снятия:", user2.Balance)
	}

	fmt.Println("\n Итоговые балансы:")
	fmt.Println(user1.Name, ":", user1.Balance)
	fmt.Println(user2.Name, ":", user2.Balance)
}
