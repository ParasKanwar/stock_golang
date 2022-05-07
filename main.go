package main

import (
	"fmt"
	customalgorithms "paraskanwar/stock_golang/custom_algorithms"
)

func main() {
	fib := customalgorithms.Supports([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 2)
	fmt.Println(fib)
}
