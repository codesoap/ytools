`ytools` is a set of simple tools to interact with YouTube via the terminal.

# Usage
```console
$ ytools-search Black Mambo
 1: Glass Animals - Black Mambo (Lyric Video)
 2: Glass Animals - Black Mambo
 3: Madrugada-Black Mambo
...
$ ytools-info 2
Glass Animals - Black Mambo
14.643.046 Views  ▲ 81.526  ▼ 1.626  Published on 17 Feb 2015

Our new record “How To Be a Human Being” featuring “Youth” and
...
$ mpv $(ytools-pick 2)
Playing: https://www.youtube.com/watch?v=H7bqZIpC3Pg
...
$ # Without an argument the comments of the last picked video are shown:
$ ytools-comments
=== Comment #1: ===
If you're here then you have good taste in music... ;)

=== Comment #2: ===
Atyphical
...
$ ytools-recommend
 1: Glass Animals - Cane Shuga
 2: Black Coast - TRNDSTTR (Lucian Remix)
 3: Glass Animals - Season 2 Episode 3 (Official Video)
...
```

For more information take a look at `man ytools`.

# Installation
The `ytools` have been tested on OpenBSD 6.4 and Xubuntu 18.04, but
will probably work on any POSIX compliant operating system, that
is [supported by go](https://github.com/golang/go/wiki/GoArm#introduction).

The easiest way to try `ytools` out is to use the prebuilt binaries,
which are available for OpenBSD and Linux (amd64 only). If you want
to properly install `ytools` on your system, I recommend building
them yourself.

## Using prebuilt binaries
```shell
# Download and extract the binaries:
wget "https://github.com/codesoap/ytools/releases/download/v1.0/ytools_bin_$(uname -s)_amd64.tar.gz"
tar -xzf "ytools_bin_$(uname -s)_amd64.tar.gz"

# Use ytools by calling the binaries directly, like so:
./ytools-search Efence - Spaceflight
```

If you want to read the manual for `ytools`, download `man/ytools.7`
and place it in `/usr/local/man/man7/`.

## Building
1. `mkdir -p "$HOME/go/src/github.com/codesoap/" && cd "$HOME/go/src/github.com/codesoap/"`
   (adapt if you've set a different `$GOPATH`)
2. `git clone https://github.com/codesoap/ytools.git && cd ytools/`
3. `go get ./...` to install dependencies (`golang.org/x/net/html`)
4. `make install` to install ytools (if you just want the binaries do
   `make all`)

To uninstall ytools call `make uninstall`.

## Porcelain
You can find some convenient scripts in `porcelain/`. They are *not*
installed by `make install`. Place the ones you like in `$HOME/bin/`
(make sure this directory is in your `$PATH`).
