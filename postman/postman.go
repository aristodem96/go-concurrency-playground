package postman

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var mailByPostman = map[int]string{
	1:  "Family letter",
	2:  "Business letter",
	3:  "Love letter",
	4:  "Invitation letter",
	5:  "Complaint letter",
	6:  "Thank you letter",
	7:  "Apology letter",
	8:  "Congratulation letter",
	9:  "Condolence letter",
	10: "Recommendation letter",
}

func postman(ctx context.Context, transferPoint chan<- string, n int, mail string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("I am postman number:", n, "and I have been stopped")
			return
		default:
			fmt.Println("I am postman number:", n, "take a letter")
			time.Sleep(1 * time.Second)
			fmt.Println("I am postman number:", n, "and I have brought the letter: ", mail)

			select {
			case transferPoint <- mail:
				fmt.Println("I am postman number:", n, "and I have sent the letter: ", mail)
			case <-ctx.Done():
				fmt.Println("I am postman number:", n, "and I have been stopped")
				return
			}
		}
	}
}

func PostmanPool(ctx context.Context, postmanCount int) <-chan string {
	mailTransferPoint := make(chan string)

	wg := &sync.WaitGroup{}

	for i := 1; i <= postmanCount; i++ {
		wg.Go(func() {
			postman(ctx, mailTransferPoint, i, postmanToMail(i))
		})
	}

	go func() {
		wg.Wait()
		close(mailTransferPoint)
	}()

	return mailTransferPoint
}

func postmanToMail(postmanNumber int) string {
	mail, ok := mailByPostman[postmanNumber]
	if !ok {
		return "Unknown letter"
	}
	return mail
}
