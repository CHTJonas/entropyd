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
$ entropyd --help

Usage of entropyd:
  -4    force the use of IPv4 for the HTTP connection
  -6    force the use of IPv6 for the HTTP connection
  --dry-run
        makes a request for 512 bits of entropy but writes to stdout instead of the kernel entropy pool
  --max int
        maximum amount of entropy (in bits) in a HTTP request (default 8128)
  --min int
        minimum amount of entropy (in bits) in a HTTP request (default 64)
  --poll int
        interval (in milliseconds) at which to poll the kernel entropy pool (default 200)
  --target int
        target amount of entropy (in bits) to store in the kernel entropy pool (default 3072)
  --url string
        URL of the remote entropy server (default "https://entropy.malc.org.uk/entropy/")
```

## Copyright

Copyright (c) 2018 Malcolm Scott.\
Copyright (c) 2019–2020 Charlie Jonas.\
The code here is released under Version 2 of the GNU General Public License.\
See the LICENSE file for full details.
