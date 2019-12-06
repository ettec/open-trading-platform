package common

import "fmt"

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckWithMsg(e error, msg string) {
	if e != nil {
		panic(fmt.Errorf(msg, e))
	}
}