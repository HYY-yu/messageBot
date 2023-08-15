# messageBot


This is a bot for receive message from facebook messenger.


### Architecture

![截屏2023-08-15 12.53.41.png](./arch.png)

### Sequence

```mermaid
sequenceDiagram
participant fb 
participant httpServer

fb ->> httpServer: POST /webhook
par response
    httpServer -->> fb: 200 OK
and goroutine
    httpServer ->> webhooker: HandleMessage()
end
loop webhokker.handlers
    webhooker ->> messageHandler: Handle()
end
```

### More

If we have powerful computer (8H16G + GPU)，We can use the [LocalAI](https://localai.io/) to run nlp model in local.
Now it will ask huggingface.io for help. 