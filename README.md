# Stream Engine

Stream Engine is a small utility written in Go that opens a WebSocket server on top of a text stream from `STDIN`. This is particularly useful for streaming content to multiple subscribers from a script that, for whatever reason, you can (or should) only run in one process at a time.

This was originally written to stream sensor data from a serial connection to multiple browsers for a project in graduate school. I have used it since for similar purposes, as well as to efficiently stream the regular output of relatively expensive computations to multiple consumers.

## Example

```
while sleep 2; do date; done | ./streamengine
```

Now, in a few web browser windows:

```
const ws = new WebSocket('ws://localhost:8080');

ws.addEventListener('message', e => {
    console.log(e.data);
});
```

All of the browser windows will receive the same timestamps at the same time.

## Install

This is confirmed working with Go 1.9 and untested on lower versions (though it has been only lightly modified since it was written in Go 1.6, so it likely runs on fairly old versions).

```
go get github.com/jonahsmith/streamengine
```

Note that this has `gorilla/websocket` as a dependency (loaded into `vendor/`). I've also added [`dep`](https://github.com/golang/dep) metadata in the hope that `dep` is eventually integrated into the Go toolchain.

## Flags

The `-port` flag can be used to specify the port on which `streamengine` serves. Defaults to `:8080`.

## How does this differ from `websocketd`?

[`websocketd`](https://github.com/joewalnes/websocketd) is awesome, but it forks a new process for each WebSocket connection. This is useful for bidirectional communication, but is not a great fit when streaming from necessarily unitary source of data in applications that do not require browser-to-server communication.