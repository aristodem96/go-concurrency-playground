package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aristodem96/go-concurrency-playground/blacksmith"
	"github.com/aristodem96/go-concurrency-playground/miner"
	"github.com/aristodem96/go-concurrency-playground/postman"
)

const workDuration = 5 * time.Second

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), workDuration)
	defer cancel()

	coalForBlacksmith := make(chan int)

	coalTransferPoint := miner.MinerPool(ctx, 3)
	mailTransferPoint := postman.PostmanPool(ctx, 3)
	toolTransferPoint, leftoverCoal := blacksmith.Pool(coalForBlacksmith, 5)

	initTime := time.Now()

	// Each result is written by exactly one goroutine and read only after
	// wg.Wait(), so no extra synchronization is needed.
	var coal int
	var mails []string
	var tools []string

	wg := &sync.WaitGroup{}

	wg.Go(func() {
		defer close(coalForBlacksmith)
		for c := range coalTransferPoint {
			coal += c
			coalForBlacksmith <- c
		}
	})

	wg.Go(func() {
		for m := range mailTransferPoint {
			mails = append(mails, m)
		}
	})

	wg.Go(func() {
		for t := range toolTransferPoint {
			tools = append(tools, t)
		}
	})

	wg.Wait()

	fmt.Println("Total coal mined: ", coal)
	fmt.Println("Total mails delivered: ", len(mails))
	fmt.Println("Total tools made: ", len(tools))
	fmt.Println("Coal left over: ", leftoverCoal.Load())
	fmt.Println("Total time taken: ", time.Since(initTime))
}
