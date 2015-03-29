# canvasflood

Send fire-and-forget UDP packets to an IP with (x, y, red, green, blue) and canvasflood will dutifully display it into the raw linux framebuffer!


## Compilation

Make sure your `GOPATH` is set.  I do `export GOPATH=$(pwd)` when I'm running in this working directory, but I'm not sure if that's kosher golang yet.  This env variable is needed to compile.

```shell
$ go get .
$ go build .
```

## Running Locally

This has only been tested with Linux so far.

First, switch over to `/dev/fb0` by pressing `ctrl-alt-<F1>`.  You should see the familiar text mode TTY you normally get when you do this.

Then, run `canvasflood`.  I reccomend doing this either in a tmux/screen session you can detach or a separate framebuffer, ssh session, etc.  The logging of the process interrupts the framebuffer rendering (although that is sort of a fun glitch effect unto itself!).

```
$ ./canvasflood
```

The screen will clear, and possibly still show your login prompt.  Send UDP packets with `x y red green blue`, where all parameters are integers. `x` and `y` are within the max screen width and height of the screen, and `red`, `green`, and `blue` are 0-255.

## TODOs

- change packet format to something cleaner (fewer bytes)
- error handling (:warning: this is really important, current daemon is hell of fragile)
- resize framebuffer
- JSON API for metadata about screen width/height, stats, etc.
- UDP packet protocol to get value of framebuffer color at a given x, y
- filters?
- accept more than one pixel at a time maybe?
