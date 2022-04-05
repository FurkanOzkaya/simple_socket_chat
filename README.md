# Go ~ Simple Real Time Messaging

Simple Real Time Messaging One to One.

GoRoutine is not used in this example please check other go_socket example.

<br/>

![WEBSOCKET POSTMAN](./images/postman1.PNG)

<br/>

![WEBSOCKET POSTMAN2](./images/postman2.PNG)


## Information:

Three  operation:

- connect: store user and it's connection to map
- message: send message to another user
- disconnet: delete user from map

```
switch {
        case wsmessage.Operation == CONNECT:
            pass
        case wsmessage.Operation == MESSAGE:
            pass
        case wsmessage.Operation == DISCONNECT:
            pass
        default:
            pass
        }
```

Storing User with connection info in struct.
```
type UserPool struct {
	Clients map[string]*websocket.Conn
}
```

Flow:

- Message comes to server.
- Server checks target user and finds target user's connection
- Server redirect message to founded connection
