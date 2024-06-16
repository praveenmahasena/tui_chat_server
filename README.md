# TUI Chat Server

This is a simple chat server app written in [Golang](https://go.dev/dl/) which uses TCP underneath to make connections and communicate.



## Features
- TUI chat server
- Can connect via multiple terminals

### Pre Request
- [Golang](https://go.dev/dl/)
- [Telnet](https://www.telnet.org/htm/faq.htm)

### Installing

```
bash
git@github.com:praveenmahasena/tui_chat_server.git
```

### Start
```
bash
cd tui_chat_server
make server
make run
```

By dialing `PORT=42069` you could connect to the server locally with telnet and start sending chat messages
