package main

import "fmt"

func main(){
    n := 5 
    result1 := sum_to_n_a(n)
    fmt.Printf("Sum of numbers up to %d using sum_to_n_a() is: %d\n", n, result1)
    result2 := sum_to_n_b(n)
    fmt.Printf("Sum of numbers up to %d using sum_to_n_b() is: %d\n", n, result2)
    result3 := sum_to_n_c(n)
    fmt.Printf("Sum of numbers up to %d using sum_to_n_c() is: %d\n", n, result3)
}

func sum_to_n_a(n int) int {
	// your code here
    sum:= 0
    for i := 1; i <= n; i++{
        sum += i 
    }
    return sum
}

// sum_to_n_a has a time complexity of O(n) because it iterates from 1 to n once.

func sum_to_n_b(n int) int {
	// your code here
    return n * (n + 1) / 2
}

// sum_to_n_b has a time complexity of O(1) because it calculates using a formula rather than iteration. 

func sum_to_n_c(n int) int {
	if n == 0 {
        return 0
    }

    return n + sum_to_n_c(n-1)
}

// sum_to_n_c has a time complexity of O(n) because of the recursive call stack