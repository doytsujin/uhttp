package uhttp_test

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dnesting/uhttp"
)

func ExampleClient_sSDP() {
	// This example performs an SSDP M-SEARCH to the local Multicast SSDP address.
	// It uses the uhttp.Client so as to receive multiple responses, but this same
	// transport could be used with a stock http.Client is you're OK only getting
	// the first response.
	client := uhttp.Client{
		Transport: &uhttp.Transport{
			// SSDP requires upper-case header names
			HeaderCanon: func(n string) string { return strings.ToUpper(n) },
			Repeat:      uhttp.RepeatAfter(50*time.Millisecond, 1),
		},
	}

	// Build M-SEARCH request
	req, _ := http.NewRequest("M-SEARCH", "", nil)
	req.URL.Host = "239.255.255.250:1900"
	req.URL.Path = "*"
	req.Header.Add("MAN", `"ssdp:discover"`)
	req.Header.Add("MX", "1")
	req.Header.Add("ST", "upnp:rootdevice")
	req.Header.Add("CPFN.UPNP.ORG", "Test")

	err := client.Do(req, 3*time.Second, func(sender net.Addr, resp *http.Response) error {
		fmt.Printf("From %s:\n", sender)
		resp.Write(os.Stdout)
		fmt.Println("---")
		return nil
	})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
