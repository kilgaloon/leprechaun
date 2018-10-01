### API

This package is used to provide API (socket/unix) between cli and client. API receives command from cli and return kind of informations through commands that are registered for Agents.

It implements `registry` interface

```
type registry interface {
    RegisterCommandSocket() *Registrator
}
```

When building new socket with `Register` method needs to recieve this interface and will start listening for commands from cli.

`Registrator` is Agent that is registered with name and commands

Reference for this can be found in client/commands.go
