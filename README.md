`ytools` is a set of simple tools to interact with youtube via the terminal.

# Usage
```console
$ ytools-search Black Mambo
 1: Glass Animals - Black Mambo (Lyric Video)
 2: Glass Animals - Black Mambo
 3: Madrugada-Black Mambo
 4: The Black Mamba - Believe (Official Music Video)
 5: Madrugada - Black Mambo [Official Music Video] [2001]
 6: Black Mambo & Ácido Pantera  - Efecto Manglar (Official Video)
 7: Glass Animals - Hazey
 8: Black Mambo
 9: Black Mambo - Ritual (Official Music Video)
10: Glass Animals - Black Mambo (Live From Capitol Studios)
11: The Black Mamba & Aurea - Wonder Why (Lyric Video)
12: The Black Mamba - It Ain't You (Official Video)
$ mpv $(ytools-pick 2)
Playing: https://www.youtube.com/watch?v=H7bqZIpC3Pg
...
```

# Todo
- `ytools-info` for displaying title, views, likes, description, ...
- `ytools-recommend` to list recommended videos for a search result
