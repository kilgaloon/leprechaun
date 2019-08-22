# Idea

## Daemon

Daemon is long living process that manages services and communicates with them. Daemon is responsible also for communication between user space and services, if user types command in terminal ex: `leprechaun --cmd='{agent} {command} {args}...` daemon passed this to `api.Resolver` then api resolves to which http endpoint to send request and get response.

```
+-----------------------------------------+
|                                         |
|                       +---------------+ |
|       DAEMON          |               | |
|                       |      API      | |       +--------------------------------+
|                       |               | |       |                                |
|                       +---------------+ |       |                                |
|                                         +------>+                                |
|                                         |       |       SERVICE                  |
|                                         |       |                                |
|                                         |       |                                |
+----------------+--------------------+---+       +--------------------------------+
                 |                    |
                 |                    |
                 v                    |
+----------------+---------------+    |   +--------------------------------+
|                                |    |   |                                |
|                                |    |   |                                |
|                                |    +-->+                                |
|       SERVICE                  |        |       SERVICE                  |
|                                |        |                                |
|                                |        |                                |
+--------------------------------+        +--------------------------------+


```

### API

To be registered in Daemon api, agent needs to satisfy `api.Registrator` interface which register http handles.

```
leprechaun --cmd="scheduler info"
              +
              |
              |
              |              API
      +----------------------------------------------------------+
      |       |                                                  |
      |   +---+----+                                             |
      |   | Client | Resolver(c Cmd)                             |
      |   +---+----+                                             |
      |       |                                                  |
      |       |          HTTP.Get("/scheduler/info")  +--------+ |
      |       | Info(c)+----------------------------->+ Server | |
      |       |        ^                              +---+----+ |
      |       |        |           Response               |      |
      |       |        +----------------------------------+      |
      |       |                                                  |
      +----------------------------------------------------------+
              v
         +----+----+
         | StdOut  |
         +---------+

```

## Service

Service is basically an agent just its not standalone and daemon is aware of his existence. When agent is added as service then daemon is responsible for this process and it handles all signals for stop, start and pause. When daemon is killed, all other services are also killed.

To satisfy Service interface agent need to satisfy `api.Registrator`, agent uses method `RegisterAPIHandles` to register http routes that Daemon api will register and will know how to communicate with that service. 

```
                AddService
                     +
                     |
                     |
+---------------+    |    +---------------+
|               |    |    |               |
|     AGENT     +----+--->+   SERVICE     |
|               |         |               |
+---------------+         +---------------+

```

## Agent

Agent is responsible for defining all methods of how it will be started, paused, stopped, it expose default api handles to handle workers, basically how it will build it self and how it will process recipes.

Agent is also barebone or boilerplate of something that will later become service


### Workers

Workers are one part of agent and they are responsible processing all bash commands that are coming from recipes they determine how they will be executed in which order and when. 

Number of currently active workers can be set in configuration and also number of workers in queue (Workers in queue will sit there and wait for first available slot to be free and then start working on recipe)

When worker is finished working on recipe it will notify through channel and it will be remove from stack, also there is error channel in case worker errors for some reason email will be sent with name of the worker and error message

```
+--------------------------------------+                                     +------------------------+
|                    +-------------+   |                                     |------------------------|
|                    |    STACK    +------------------------------------------|  worker  ||  worker  |-<----------+
|                    +-------------+   |                                     |------------------------|           |
|                    +-------------+   |                                     +------------------------+           |
|                    |    QUEUE    |   |                                           |                              |
|                    +------+------+   |                                           |                              |
|  WORKERS                  |          |                                           |                              |
|                           |          |                                           |                              |
|                           |          |   DoneChan    +-+                         |                              |
|                           |          +---------------+ +<------------------------+                              |
|                           |          |               +-+                                                        |
|                           |          |                                                                          |
|                           |          |                                                                          |
+--------------------------------------+                                                                          |
                            |                                                                                     |
                            |                                                                                     |
                            |                                                                                     |
                            |                                                                                     |
                            |                                                                                     |
                            |                                           +---------------------+                   |
                            |                                           |---------------------|                   |
                            +--------------------------------------------|  worker  | worker |--------------------+
                                                                        |---------------------|
                                                                        +---------------------+

```

