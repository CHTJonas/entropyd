package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/chtjonas/entropyd/pkg/logging"
	"github.com/chtjonas/entropyd/pkg/pool"
	"github.com/chtjonas/entropyd/pkg/providers/malc"
)

// Limit of ioctl requests is 1024 bytes, including header.
const maxDataBytes = 1016

// Must be no higher than that used by the server.
const maxReqBits = maxDataBytes * 8

// Software version defaults to the value below but is overridden by the compiler in Makefile.
var version = "dev-edge"

func init() {
	if runtime.GOOS != "linux" {
		fmt.Println("entropyd can only run on Linux")
		os.Exit(1)
	}
}

func main() {
	verThenExit := flag.Bool("version", false, "print software version and exit")
	ipv4Ptr := flag.Bool("4", false, "force the use of IPv4 for the HTTP connection")
	ipv6Ptr := flag.Bool("6", false, "force the use of IPv6 for the HTTP connection")
	minBitsPtr := flag.Int("min", 64, "minimum amount of entropy (in bits) in a HTTP request")
	maxBitsPtr := flag.Int("max", maxReqBits, "maximum amount of entropy (in bits) in a HTTP request")
	targetBitsPtr := flag.Int("target", 3072, "target amount of entropy (in bits) to store in the kernel entropy pool")
	pollIntervalPtr := flag.Int("poll", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	doDryRunPtr := flag.Bool("dry-run", false, "makes a request for 512 bits of entropy but writes to stdout instead of the kernel entropy pool")
	flag.Parse()

	// Print version and exit if the user asked us to.
	if *verThenExit {
		path := os.Args[0]
		fmt.Println(path, "version", version)
		os.Exit(0)
	}

	// Setup User-Agent header and IP protocol version.
	ua := "entropyd/" + version + " (+https://github.com/CHTJonas/entropyd)"
	ipv := ""
	if *ipv4Ptr {
		ipv = "tcp4"
	}
	if *ipv6Ptr {
		ipv = "tcp6"
	}

	// Instantiate the actual entropy client and open the Linux kernel entropy pool.
	cl := malc.NewEntropyClient(*minBitsPtr, *maxBitsPtr, ua, ipv)
	pl := pool.OpenPool()
	defer pl.Cleardown()

	// Perform an dry-run and exit if the user asked us to.
	if *doDryRunPtr {
		entropy, err := cl.FetchEntropy(16)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Entropy: %s", entropy.Data)
		os.Exit(0)
	}

	logging.Log("entropyd started successfully",
		logging.LogString("path", os.Args[0]),
		logging.LogString("version", version),
	)

	interval := time.Duration(*pollIntervalPtr)
	pl.Run(interval, *targetBitsPtr, *maxBitsPtr, cl)
}
