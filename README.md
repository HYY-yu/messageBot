# messageBot


This is a bot for receive message from facebook messenger.


### Architecture

```mermaid
classDiagram
    direction BT
    namespace Repository{
        class MessageRepository {
            <<interface>>
            +Save()
            +Saves()
            +Read()
        }
        class MessageTemplateRepository {
            <<interface>>
            +QueryOne()
        }
    }
    note for MessageRepository "The Repository is a bridge from code to database. "
    
    class  Webhooker{
        -[]MessageHandler handlers
        +AddMessageHandler()
        +HandleMessage()
        +Shutdown()
    }
    class MessageHandler{
        <<interface>>
        +Handler()
    }
    class DBHandler{
        -MessageRepository repo
        +Handler()
    }
    class NLPHandler{
        +Handler()
    }
    class BotHandler{
        -MessageTemplateRepository repo
        +Handler()
    }
    DBHandler ..|> MessageHandler
    NLPHandler ..|> MessageHandler
    BotHandler ..|> MessageHandler
    MessageHandler "*"-->"1" Webhooker
    
    class Messenger {
        +SendMessage()
    }
    note for Messenger "wrapped the FB messenger API "

    BotHandler ..> Messenger
    BotHandler ..|> MessageTemplateRepository
    DBHandler ..|> MessageRepository
    
```

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