/** Given code **/
// package main

// import "fmt"

// func main() {
// 	cnp := make(chan func(), 10)
// 	for i := 0; i < 4; i++ {
// 		go func() {
// 			for f := range cnp {
// 				f()
// 			}
// 		}()
// 	}
// 	cnp <- func() {
// 		fmt.Println("HERE1")
// 	}
// 	fmt.Println("Hello")
// }

/** Completed code **/
package main

import (
	"fmt"
	"sync"
)

func main() {
	cnp := make(chan func(), 10)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for f := range cnp {
				fmt.Printf("Worker %d executing task\n", workerID)
				f()
			}
		}(i + 1)
	}
	cnp <- func() {
		fmt.Println("HERE1")
	}
	close(cnp)
	wg.Wait()

	fmt.Println("Hello")
}

/** Example application **/
// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type CustomerRequest struct {
// 	ID      int
// 	Message string
// }

// func processCustomerRequest(request CustomerRequest) {
// 	fmt.Printf("Processing request #%d: %s\n", request.ID, request.Message)
// 	time.Sleep(time.Second) // Simulating processing time
// 	fmt.Printf("Request #%d processed\n", request.ID)
// }

// func main() {
// 	numGoroutines := 4
// 	requestChannel := make(chan func(), 10)

// 	var wg sync.WaitGroup

// 	for i := 0; i < numGoroutines; i++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
// 			for processFunc := range requestChannel {
// 				processFunc()
// 			}
// 		}(i + 1)
// 	}

// 	for i := 1; i <= 10; i++ {
// 		requestID := i
// 		requestChannel <- func() {
// 			request := CustomerRequest{
// 				ID:      requestID,
// 				Message: fmt.Sprintf("Help needed for request #%d", requestID),
// 			}
// 			processCustomerRequest(request)
// 		}
// 	}

// 	close(requestChannel)

// 	wg.Wait()

// 	fmt.Println("All customer requests processed.")
// }
