# Linux Entropy Daemon

This repository hosts the code for a simple client application to seed the Linux kernel with entropy from a remote server via HTTP. It attempts to adhere to the UNIX philosophy of doing one thing and doing it well.

This project is a Go rewrite of [entropyserver](https://hg.sr.ht/~mas90/entropyserver) by Malcolm Scott. Much of the original credit for code logic must go to him.

## Build

To compile the code execute the following in a terminal. This will produce three `entropyd` binaries for the `arm`, `arm64` and `amd64` architectures in different directories under `bin`.

```bash
make clean && make build
```

## Usage

```
Usage:
    entropyd [-4 | -6] [--dry-run] [--min MIN] [--max MAX] [-t TARGET] [-p POLL]

Flags:

    -4                          Force the use of IPv4.
    -6                          Force the use of IPv6.
    --min                       Minimum amount of entropy in bits to request (default 64).
    --max                       Maximum amount of entropy in bits to request (default 8128).
    -t, --target                Target amount of entropy in bits to store in the kernel entropy pool (default 3072).
    -p, --poll                  Interval in milliseconds at which to poll the kernel entropy pool (default 200).
    --dry-run                   Make a request for 512 bits of entropy, write to STDOUT and exit.
    -v, --version               Print version and exit.
```

## Copyright

Copyright (c) 2018 Malcolm Scott.\
Copyright (c) 2019–2021 Charlie Jonas.\
The code here is released under Version 2 of the GNU General Public License.\
See the LICENSE file for full details.
