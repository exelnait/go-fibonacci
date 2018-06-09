package main

import (
	"bufio"
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	"encoding/json"
)

type FibonacciNumber struct {
	Value int `json: "value"`
	Iteration  int `json: "iteration"`
}

type InputResult struct {
	Value int
	Error bool
}

func fibonacci() func() int {
	current, next := 1, 0
	return func() int {
		if next < 0 {
			current, next = 1, 0
		}
		current, next = current + next, current

		return current - next
	}
}

func getInput(input chan InputResult) {
	for {
		res := new(InputResult)
		res.Error = false

		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		result = strings.Replace(result, "\n", "", -1)
		if err != nil {
			fmt.Println(err)
			res.Error = true
		}

		resultInt, err := strconv.Atoi(result)
		if err != nil {
			fmt.Println(err)
			res.Error = true
		}

		res.Value = resultInt

		input <- *res
	}
}

func main() {
	var FibonacciSequence []FibonacciNumber
	calculateFibonacciNumber := fibonacci()

	var errorNumberCount int
	var rightNumberCount int
	var iteration int

	input := make(chan InputResult)

	go getInput(input)

	for {
		if rightNumberCount == 10 || errorNumberCount == 3 {
			break
		}

		fmt.Println("--> Enter fibonacci number:")

		currentNumber := calculateFibonacciNumber()

		fibonacciNumber := FibonacciNumber{Value: currentNumber, Iteration: iteration}
		fibonacciNumberJson, err := json.Marshal(fibonacciNumber)
		if err != nil {
			fmt.Println(err)
		}
		FibonacciSequence = append(FibonacciSequence, fibonacciNumber)

		select {
			case inputResult := <-input: {
				if inputResult.Value != currentNumber || inputResult.Error {
					rightNumberCount = 0
					fmt.Println("Error ", string(fibonacciNumberJson))
					errorNumberCount++
				} else {
					rightNumberCount++
				}
			}
			case <-time.After(10 * time.Second):
				fmt.Println("Time out ", string(fibonacciNumberJson))
				errorNumberCount++
		}

		iteration++
	}
}