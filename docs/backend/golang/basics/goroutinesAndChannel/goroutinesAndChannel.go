package main

import (
	"fmt"
	"sync"
)

type accountReqest struct {
	name    string
	balance chan float64
	Type    string
	Amount  float64
	Result  chan struct{}
}

var (
	Account = make(map[string]float64)
	reqChan = make(chan accountReqest)
	wg      sync.WaitGroup
)

func creatAccount(name string, amount float64) {
	done := make(chan struct{})

	reqChan <- accountReqest{
		name:   name,
		Type:   "creat",
		Amount: amount,
		Result: done,
	}

	<-done
}

func transation(name string, Type string, amount float64) {
	done := make(chan struct{})

	reqChan <- accountReqest{
		name:   name,
		Type:   Type,
		Amount: amount,
		Result: done,
	}

	<-done
}

func getBalence(name string) float64 {
	bal := make(chan float64)

	reqChan <- accountReqest{
		name:    name,
		Type:    "query",
		balance: bal,
	}
	return <-bal
}

func accountManger() {
	for req := range reqChan {
		switch req.Type {
		case "deposit":
			Account[req.name] += req.Amount
			req.Result <- struct{}{}

		case "withdrawal":
			Account[req.name] -= req.Amount
			req.Result <- struct{}{}

		case "creat":
			Account[req.name] = req.Amount
			req.Result <- struct{}{}
		case "query":
			req.balance <- Account[req.name]

		}

	}
}

func printBalence(name string) {
	// getBalence(name)
	fmt.Printf("Balance of %s: %.2f\n", name, getBalence(name))
}

func main() {
	go accountManger()

	wg.Add(2)
	// Create an account for Alice with an initial balance of 100

	go func() { creatAccount("Alice", 100); wg.Done() }()
	go func() { creatAccount("Bob", 300); wg.Done() }()

	wg.Wait()

	wg.Add(2)
	go func() { transation("Alice", "deposit", 300); wg.Done() }()
	go func() { transation("Alice", "deposit", 300); wg.Done() }()
	wg.Wait()

	wg.Add(2)
	go func() { printBalence("Alice"); wg.Done() }()
	go func() { printBalence("Bob"); wg.Done() }()
	wg.Wait()

}
