package main

import (
	"fmt"
	"crypto/rand"
	"math/big"
	"math"
	"net/http"
	"runtime"
)


func doRequest(i int, jobs <-chan string) {

	for j := range jobs {
		resp, err := http.Get("http://localhost:8000/set/" + j )
		if err != nil {
			panic(fmt.Sprintf("Error getting %v: %v\n", j, err))
		} else {
			resp.Body.Close()
		}
		fmt.Printf(" %v ", i)
		runtime.Gosched()
	}
}

func main() {

	goCount := 50
	jobCount := 10000
	jobs := make(chan string)
	
	//var requested map[string]string
	
	thirtyTwo := math.Pow(2,32)-1
	max := *big.NewInt(int64(thirtyTwo))
	
	for i:=0; i < goCount; i++ {
		go doRequest(i,jobs)
		
	}

	for i:=0; i < jobCount; i++ {
		roff,_ := rand.Int(rand.Reader,&max)
		jobs <- fmt.Sprintf("%d",roff)
		runtime.Gosched()
	}
	
	var input string
    fmt.Scanln(&input)
    fmt.Println("done")
    
}
