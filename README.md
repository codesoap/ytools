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
$ # Without an argument recommendations for the last picked result are
$ # listed:
$ ytools-recommend
 1: Glass Animals - Cane Shuga
 2: Black Coast - TRNDSTTR (Lucian Remix)
 3: Glass Animals - Season 2 Episode 3 (Official Video)
...
```

For more information take a look at `man ytools`.

# Installation
The easiest way to try `ytools` is to use the prebuilt binaries, that
are available at the [releases
page](https://github.com/codesoap/ytools/releases).

If you want to properly install `ytools` on your system, I recommend
building them yourself:

```shell
git clone git@github.com:codesoap/ytools.git
cd ytools

# Execute as root to install:
make install

# To uninstall use (again as root):
# make uninstall

# If you don't want to run make as root and don't care for the man
# page, you could alternatively run the following. This will install the
# binaries to ~/go/bin/:
# go install ./...
```

## Porcelain
You can find some convenient scripts in `porcelain/`. They are *not*
installed by `make install`. Place the ones you like in `$HOME/bin/`
(make sure this directory is in your `$PATH`).
