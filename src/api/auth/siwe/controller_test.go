package siwe

import (
	"fmt"
	"net/url"
	"testing"
)

func TestController_Challenge(t *testing.T) {

	domainURL, _ := url.Parse("localhost:3000")
	fmt.Println(domainURL)
	fmt.Println(domainURL.Hostname())
}
