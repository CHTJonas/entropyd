package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/chtjonas/entropyd/pkg/pool"
	"github.com/chtjonas/entropyd/pkg/providers/malc"
)

var (
	// Software version defaults to the value below but is overridden by the compiler in Makefile.
	version = "dev-edge"

	// Command line flags.
	versionFlag      bool
	dryRunFlag       bool
	forceIPv4Flag    bool
	forceIPv6Flag    bool
	minBitsFlag      int
	maxBitsFlag      int
	targetBitsFlag   int
	pollIntervalFlag int
)

func init() {
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.BoolVar(&versionFlag, "v", false, "print software version and exit")
	flag.BoolVar(&versionFlag, "version", false, "print software version and exit")
	flag.BoolVar(&dryRunFlag, "dry-run", false, "makes a request for 512 bits of entropy but writes to stdout instead of the kernel entropy pool")
	flag.BoolVar(&forceIPv4Flag, "4", false, "force the use of IPv4")
	flag.BoolVar(&forceIPv6Flag, "6", false, "force the use of IPv6")
	flag.IntVar(&minBitsFlag, "min", 256, "minimum amount of entropy (in bits) to request")
	flag.IntVar(&maxBitsFlag, "max", 6080, "maximum amount of entropy (in bits) to request")
	flag.IntVar(&targetBitsFlag, "t", 3072, "kernel entropy pool target value (in bits)")
	flag.IntVar(&targetBitsFlag, "target", 3072, "kernel entropy pool target value (in bits)")
	flag.IntVar(&pollIntervalFlag, "p", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	flag.IntVar(&pollIntervalFlag, "poll", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	flag.Parse()
}

func main() {
	// Instantiate logger.
	logger := log.New(os.Stdout, "", 0)

	// Print version and exit if the user asked us to.
	if versionFlag {
		logger.Println(version)
		os.Exit(0)
	}

	// Setup User-Agent header.
	ua := "entropyd/" + version + " (+https://github.com/CHTJonas/entropyd)"

	// Set IP protocol version.
	var ipv string
	if forceIPv4Flag && !forceIPv6Flag {
		ipv = "tcp4"
	}
	if !forceIPv4Flag && forceIPv6Flag {
		ipv = "tcp6"
	}

	// Instantiate the entropy client.
	cl := malc.NewEntropyClient(minBitsFlag, maxBitsFlag, ua, ipv)

	// Perform an dry-run and exit if the user asked us to.
	if dryRunFlag {
		entropy, err := cl.FetchEntropy(512)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}
		logger.Printf("Entropy: %s", entropy.Data)
		os.Exit(0)
	}

	// Don't enter the main loop unless we're on Linux.
	if runtime.GOOS != "linux" {
		logger.Println("entropyd can only run on Linux")
		os.Exit(1)
	}

	// Open the Linux kernel entropy pool.
	pl, err := pool.OpenPool()
	if err != nil {
		logger.Printf("Failed to access kernel entropy pool: %v\n", err)
		os.Exit(10)
	}
	defer pl.Cleardown()

	pl.SetLogger(logger)
	logger.Println("entropyd", version, "started successfully")

	interval := time.Duration(pollIntervalFlag)
	err = pl.Run(interval, targetBitsFlag, maxBitsFlag, cl)
	if err != nil {
		fmt.Printf("Failed to seed kernel entropy pool: %v\n", err)
		os.Exit(125)
	}
}
