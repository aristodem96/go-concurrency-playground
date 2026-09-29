package blacksmith

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

const totalPrice = 10

var tools = []string{"Hammer", "Axe", "Sword", "Shield", "Pickaxe"}

func blacksmith(coalTransferPoint <-chan int, toolTransferPoint chan<- string, leftover *atomic.Int64, n int) {
	coalSummary := 0

	for c := range coalTransferPoint {
		fmt.Println("I am blacksmith number:", n, "and I start working with coal")
		coalSummary += c
		fmt.Println("I am blacksmith number:", n, "and I have received coal: ", c)
		fmt.Println("Total coal received: ", coalSummary, "By blacksmith number: ", n)
		makeTool(&coalSummary, toolTransferPoint, n)
	}

	if coalSummary > 0 {
		fmt.Println("I am blacksmith number:", n, "and I have coal left over: ", coalSummary)
	}
	leftover.Add(int64(coalSummary))
}

// Pool starts blacksmithCount blacksmiths that turn coal into tools.
// The returned counter holds coal that was not enough for another tool;
// it is final only after the tools channel has been closed.
func Pool(coalTransferPoint <-chan int, blacksmithCount int) (<-chan string, *atomic.Int64) {
	toolTransferPoint := make(chan string)
	leftover := &atomic.Int64{}
	wg := &sync.WaitGroup{}
	for i := 1; i <= blacksmithCount; i++ {
		wg.Go(func() {
			blacksmith(coalTransferPoint, toolTransferPoint, leftover, i)
		})
	}
	go func() {
		wg.Wait()
		close(toolTransferPoint)
	}()

	return toolTransferPoint, leftover
}

func makeTool(coalSum *int, toolTransferPoint chan<- string, n int) {
	for *coalSum >= totalPrice {
		tool := tools[rand.IntN(len(tools))]

		*coalSum -= totalPrice
		time.Sleep(time.Second * 1)
		toolTransferPoint <- tool

		fmt.Println("I am blacksmith number:", n, "and I have made a tool: ", tool)
	}
}
