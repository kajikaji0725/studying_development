package main

import "fmt"

func main() {
	s := "fae323jfa{}{}fas24232"
	S := ""
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			S += string(s[i])
		}
	}
	fmt.Println(S)

}
