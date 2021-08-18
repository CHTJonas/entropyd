# Linux Entropy Daemon

This repository hosts the code for a simple application to seed the Linux kernel entropy pool with random data from a remote server. It attempts to adhere to the UNIX philosophy of doing one thing and doing it well.

Parts of this project are a rewrite of [entropyserver](https://hg.sr.ht/~mas90/entropyserver) by [Malcolm Scott](https://www.cl.cam.ac.uk/~mas90/). Much of the credit for the idea behind this and the original code logic must go to him.

The provider from which entropy is obtained is designed to be modular and extensible, however the only backend currently supported is `entropy.malc.org.uk`. I hope to expand the number of providers in the future and to mix entropy from different providers with whitening.

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

## Installation

Pre-built statically-linked binaries for a variety of architectures are available to download from [GitHub Releases](https://github.com/CHTJonas/entropyd/releases). To compile from source you will need a suitable [Go toolchain installed](https://golang.org/doc/install):

```bash
git clone https://github.com/CHTJonas/entropyd.git
cd entropyd
make clean && make build
```

## Copyright

Copyright (c) 2018 Malcolm Scott.\
Copyright (c) 2019–2021 Charlie Jonas.\
The code here is released under Version 2 of the GNU General Public License.\
See the LICENSE file for full details.
