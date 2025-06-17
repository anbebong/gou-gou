package main

import "fmt"

func main() {
	a, b := 7, 3

	fmt.Println("a + b =", a+b)
	fmt.Println("a - b =", a-b)
	fmt.Println("a * b =", a*b)
	fmt.Println("a / b =", a/b)
	fmt.Println("a % b =", a%b)

	fmt.Println("a == b:", a == b)
	fmt.Println("a > b:", a > b)
	fmt.Println("a <= b:", a <= b)

	c := a > 5 && b < 5
	fmt.Println("a > 5 && b < 5:", c)

	fmt.Println("Bitwise a & b:", a&b)
	fmt.Println("Bitwise a | b:", a|b)
	fmt.Println("Bitwise a ^ b:", a^b)
	fmt.Println("a << 1:", a<<1)
	fmt.Println("a >> 1:", a>>1)
}
