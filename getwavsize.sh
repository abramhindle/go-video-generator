soxi $1 | fgrep samples | awk '{print int(30 * $5 / 44100)}'
