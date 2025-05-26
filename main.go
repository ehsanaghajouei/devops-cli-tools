package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "cli-tools"}

	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(listenPortsCmd)
	rootCmd.AddCommand(checkPortCmd)
	rootCmd.AddCommand(dnsLookupCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var httpCmd = &cobra.Command{
	Use:   "http http(s)://[url]",
	Short: "Send a GET request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(body))
		fmt.Println("")
		fmt.Println("HTTP Status Code:", resp.StatusCode)
	},
}

var listenPortsCmd = &cobra.Command{
	Use:   "listen-ports",
	Short: "Check which ports are listening on localhost",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		for port := 1; port <= 65535; port++ {
			wg.Add(1)
			go checkPort(port, &wg)
		}
		wg.Wait()
	},
}

func checkPort(port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := "127.0.0.1:" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err == nil {
		fmt.Printf("Port %d is listening\n", port)
		conn.Close()
	}
}

var checkPortCmd = &cobra.Command{
	Use:   "check-port [host] [port]",
	Short: "Check if a specific port is open on a host",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port := args[1]
		address := net.JoinHostPort(host, port)

		timeout := 5 * time.Second
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			fmt.Printf("Connection to %s failed: %s\n", address, err)
			return
		}
		defer conn.Close()

		fmt.Printf("Successfully connected to %s\n", address)
	},
}

var dnsLookupCmd = &cobra.Command{
	Use:   "dns-lookup [dns_record] [dns_server]",
	Short: "Perform a DNS lookup using system or custom DNS server",
	Args:  cobra.RangeArgs(1, 2), // Accepts 1 or 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		dnsRecord := args[0]
		var resolver *net.Resolver

		if len(args) == 2 {
			// Use custom DNS server
			dnsServer := args[1]
			resolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("udp", net.JoinHostPort(dnsServer, "53"))
				},
			}
		} else {
			// Use system default DNS settings
			resolver = net.DefaultResolver
		}

		ips, err := resolver.LookupIPAddr(context.Background(), dnsRecord)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DNS lookup failed: %v\n", err)
			os.Exit(1)
		}

		for _, ip := range ips {
			fmt.Printf("%s IN A %s\n", dnsRecord, ip.String())
		}
	},
}
