package miner

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func miner(
	ctx context.Context,
	transferPoint chan<- int,
	n int,
	power int,
) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("I am miner number: ", n, " and I have been stopped")
			return
		default:
			fmt.Println("I am miner ", n, " and I am mining ")
			time.Sleep(1 * time.Second)
			fmt.Println("I am miner number: ", n, " and I have mined! ", power)

			select {
			case transferPoint <- power:
				fmt.Println("I am miner number: ", n, " and I have sent the Coal: ", power)
			case <-ctx.Done():
				fmt.Println("I am miner number: ", n, " and I have been stopped")
				return
			}
		}
	}
}

func MinerPool(ctx context.Context, minerCount int) <-chan int {
	coalTransferPoint := make(chan int)
	wg := &sync.WaitGroup{}
	for i := 1; i <= minerCount; i++ {
		wg.Go(func() {
			miner(ctx, coalTransferPoint, i, i*10)
		})
	}

	go func() {
		wg.Wait()
		close(coalTransferPoint)
	}()

	return coalTransferPoint
}
