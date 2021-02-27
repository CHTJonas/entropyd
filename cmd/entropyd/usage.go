package main

const usage = `Usage:
    entropyd [-4 | -6] [--dry-run] [--min MIN] [--max MAX] [-t TARGET] [-p POLL]

Flags:

    -4                          Force the use of IPv4.
    -6                          Force the use of IPv6.
    --min                       Minimum amount of entropy in bits to request (default 256).
    --max                       Maximum amount of entropy in bits to request (default 6080).
    -t, --target                Target amount of entropy in bits to store in the kernel entropy pool (default 3072).
    -p, --poll                  Interval in milliseconds at which to poll the kernel entropy pool (default 200).
    --dry-run                   Make a request for 512 bits of entropy, write to STDOUT and exit.
    -v, --version               Print version and exit.`
