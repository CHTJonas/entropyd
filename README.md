# Linux Entropy Daemon

This repository hosts the code for a simple client application to seed the Linux kernel with entropy from a remote server via HTTP. It attempts to adhere to the UNIX philosophy of doing one thing and doing it well.

## Build

To compile the code execute the following in a terminal. This will produce three `entropyd` binaries for the `arm`, `arm64` and `amd64` architectures in different directories under `bin`.

```bash
make clean && make build
```

## Usage

```
Usage of entropyd:
  -max int
        maximum amount of entropy (in bits) in a HTTP request (default 8192)
  -min int
        minimum amount of entropy (in bits) in a HTTP request (default 512)
  -poll int
        interval (in milliseconds) at which to poll the kernel entropy pool (default 200)
  -url string
        URL of the remote entropy server (default "https://entropy.malc.org.uk/entropy/")
```

## Copyright

Copyright (c) 2019–2020 Charlie Jonas.
The code here is released under Version 2 of the GNU General Public License.
See the LICENSE file for full details.
