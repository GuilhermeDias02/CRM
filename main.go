package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Votre choix: ")

	input, _ := reader.ReadString('\n')

	fmt.Println(input)
}