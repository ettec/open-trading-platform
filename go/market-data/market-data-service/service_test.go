package main

import (
	"strings"
	"testing"
)

func Test_service_Subscribe(t *testing.T) {

	p := "iexg-market-data-gateway-0"
	idx := strings.LastIndex(p, "-")
	r := []rune(p)
	serviceName := string( r[0:idx])
	podId := string(r[idx+1:len(p)])

	log.Printf("%v and %v", serviceName, podId)

}
