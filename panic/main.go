package main

import "fmt"

func main() {
	defer catch() // # play with this code, comment uncomment to see the effect

	sum, err := theFuncWillPanic()
	if err != nil {
		fmt.Println("err")
	}

	fmt.Println("hasil sum is: ", sum)
}

func catch() {
	if r := recover(); r != nil {
		fmt.Println("Panic occured :", r)
	} else {
		fmt.Println("Application running perfectly")
	}
}

func theFuncWillPanic() (int, error) {
	defer catch() // # play with this code, comment uncomment to see the effect

	sum := 0
	for i := 1; i < 100; i++ {
		if i == 60 {
			panic("this is panic happening") // create custom panic
		}
		sum += i
	}
	return sum, nil
}
