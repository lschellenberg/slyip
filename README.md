# YIP - Identity Provider

**YIP** is an identity provider to authenticate users.
It uses *SIWE - Sign In With Ethereum*, Email-Pin &
Credential authentication. It is also intended to use
this service for deploying and syncing **SLYWallet** contracts.


## **YIP** Functions 

1. Authenticate a **SLYWallet** with SIWE (EIP1237 & EOA)
2. Deploy a **SLYWallet**
3. Synchronize a **SLYWallet** information for quick access


## Functionality

* [Pin Authentication](./docs/pin_authentication.md)

## Development

### Before You Start

The YIP service depends on a postgres database that is being developed with
[Goose](https://github.com/pressly/goose) and [Jet](https://github.com/go-jet/jet). You need to download the packages by following
the instructions on the web page.

    1. Create Certificates

    CERTDIR=DIRECTORY_OF_YOUR_CHOICE
    openssl genrsa -out $(CERTDIR)/app.rsa 4096
    openssl rsa -in $(CERTDIR)/app.rsa -pubout > $(CERTDIR)/app.rsa.pub

    2. Prepare Database 

    cd ./internal/goose/migrations 
    goose postgres "postgres://USER:PWD@DB_URL/postgres?sslmode=disable" up
    jet -dsn=postgres://USER:PWD@$DB_URL/postgres?sslmode=disable -schema=yip -path=./.gen

    [you can use the Makefile + .local.env as well...]

    3. Rename yip.json.example -> yip.json and modify values
    
    cp yip.json.example yip.json

    4. Run the service

    go run cmd/service/main.go

### Database/Migration

YIP uses **goose** for migration. The migration directory is at

    ./internal/goose/migrations

Create a new migration with

    cd ./internal/goose/migrations
    goose create NEW_MIGRATION_NAME sql

Make sure to create the ORM models. Run

    jet -dsn=postgres://USER:PWD@$DB_URL/postgres?sslmode=disable -schema=yip -path=./.gen

NOTE: The migration will be executed when the service starts the next time.


### Test End To End

Run the service

	go build -o YOUR_OUT_DIRECTORY/yip_service cmd/service/main.go
	chmod +x YOUR_OUT_DIRECTORY/yip_service 
	./YOUR_OUT_DIRECTORY/yip_service

In another terminal

    go test ./e2e

## CLI

Build the CLI

    go build -o $GOPATH/bin/yip cmd/cli/main.go


and use it with

    yip COMMAND

### Register Admin User

    yip register YIP_ADMIN_PASSWORD USER_EMAIL USER_PASSWORD

### Login Admin User

    yip login USER_EMAIL USER_PASSWORD AUDIENCE

**e.g. yip login somebody@mail.com mypassword*

### User Info

    yip info YIP_ADMIN_PASSWORD

### Remote Connect

Remote Sign In is designed to support websocket and a http poll mechanism. The http poll is run by
a POST on the same endpoint


    POST {{url}}/api/v1/auth/session

to rely on the same message type as the websocket connection. The type is

    WebsocketMessage

    {
      "messageType": string,
      "sessionId": uuid,
      "payload": JSONObject 
    }

and will serve for request and response. The payload structure depends on the messageType.

The error message is also a WebsocketMessage with the specific form
    
    WebsocketMessage / Error

    {
      "messageType": "session_error",
      "sessionId": uuid,
      "payload": {
            "code": string,
            "msg": string,
            "details": string
       } 
    }

A flow is sequence of message exchanges between remote party to achieve an objective.
A session always must be created and the sessionId connects the parties. Once the
objective is fulfilled the session can be closed. Any party can close the session at any
time with the message.

    WebsocketMessage / Close Session

    {
      "messageType": "session_close",
      "sessionId": uuid,
      "payload": {} 
    }

    RESPONSE

    {
        "messageType":"session_closed",
        "sessionId":"3ac30b8e-a9b1-4622-82c2-995d682aa943",
        "payload":{ 
        }
    }

    OR
    
    ERROR Response

#### Auth Flow

The authentication flow has 2 parties. The MApp & the Wallet. 

#### The MApp

The MApp starts the session with the message


    WebsocketMessage / Create Session

    {
        "messageType":"create_session",
        "sessionId":"",
        "payload": {
            "clientId":"d798c41c-1afe-465e-9143-75c5f111a1cb",
            "sessionType":"auth_session"
        }
    }

    RESPONSE

    {
        "messageType":"session_created",
        "sessionId":"3ac30b8e-a9b1-4622-82c2-995d682aa943",
        "payload":{
            "clientId":"d798c41c-1afe-465e-9143-75c5f111a1cb",
            "qrCodeContent":"http://localhost:8080/session?sid=3ac30b8e-a9b1-4622-82c2-995d682aa943\u0026cid=d798c41c-1afe-465e-9143-75c5f111a1cb\u0026flow=auth_session","sessionId":"3ac30b8e-a9b1-4622-82c2-995d682aa943",
            "sessionType":"auth_session"
        }
    }

    OR
    
    ERROR Response
 
The clientId is given by the service. The auth_session indicates the desired flow.
The response has the newly created session id. The qrCodeContent should be displayed
in the app in a QRCode. In the next steps the MApp polls the status by

    WebsocketMessage / Ping 

    {
        "messageType":"ping_token",
        "sessionId":"203ef152-8808-41c3-aaff-68d3546a27dc",
        "payload":{
        }
    }
 
    RESPONSE

    {
        "messageType":"ping_token_response",
        "sessionId":"203ef152-8808-41c3-aaff-68d3546a27dc",
        "payload":{
            "authState":"pending" | "failed" | "success",
            "token":null | {"token": string, ...}
        }
    }
    
    OR
    
    ERROR Response

When the authState is **success** the token is returned as well.

#### Wallet

The wallet app parses the qr code and sends.

    WebsocketMessage / Connect With Account

    {
        "messageType":"connect_with_account",
        "sessionId":"ec874cc1-ae0f-4ab1-86e1-0d762a752b52",
        "payload":{
            "eoa":"0x4E345039EE45217fC99a717a441384A46dD2b85C",
            "SLYWalletAddress":0x030E4BFabdF1d5463B92BBC4fA8cE8587c7BA079" | "",
            "chainId":"11155111"
        }
    }

    RESPONSE

    {
        "messageType":"eth_sign",
        "sessionId":"0b869d97-3b8a-4d30-b920-d67e97397742",
        "payload":{
            "address":"0x4E345039EE45217fC99a717a441384A46dD2b85C",
            "chainId":"11155111",
            "challenge":"localhost:3000 wants you to sign in with your Ethereum account:\n0x4E345039EE45217fC99a717a441384A46dD2b85C\n\n\nURI: https://localhost:3000\nVersion: 1\nChain ID: 11155111\nNonce: VkrowCqUhM82dT9U\nIssued At: 2024-04-17T19:36:24Z",
            "domain":"https://localhost:3000"
        }
    }
 
    OR
    
    ERROR Response

The result should be immediately used to provide the signature.

    // PSEUDO CODE
    let signature YOUR_WALLET.Sign(msg.payload.challenge)

and the result is send to the session.

    WebsocketMessage / ETH Sign Response

    {   
        "messageType":"eth_sign_response",
        "sessionId":"2fa215d6-0b05-45bb-b652-84108893d10d",
        "payload":{
            "message":"localhost:3000 wants you to sign in with your Ethereum account:\n0x4E345039EE45217fC99a717a441384A46dD2b85C\n\n\nURI: https://localhost:3000\nVersion: 1\nChain ID: 11155111\nNonce: HJHTk2Q8FJc69b7J\nIssued At: 2024-04-17T19:39:22Z",
            "signature":"0x127bda96533251697387bff7c71861d2cded6df4ddb7ee4a7f19754d32845ddf0c641b4ff70263b05aba9dc4fa09e79ececd2b64b42ad017bdc8e0a17728cb8e1c"
        }
    }

    RESPONSE

    {
        "messageType":"eth_sign_verification_response",
        "sessionId":"37b3f64e-39e5-40d9-b39f-0cd6f5c87409",
        "payload":{
            "domain":"localhost:3000",
            "originalAddress":"0x4E345039EE45217fC99a717a441384A46dD2b85C",
            "recoveredAddress":"0x4E345039EE45217fC99a717a441384A46dD2b85C",
            "uri":"https://localhost:3000"
        }
    }
    
    OR
    
    ERROR Response