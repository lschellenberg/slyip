# Pin Authentication 

## Basic Flow

The user requests a pin by providing an email address and a ECDSA address. An email
with a generated 6-digit pin will be sent to the email address. The user is supposed
to enter the pin and sign the message with the private key. The signature and the pin
are verified and a JWT is been sent back.

## Endpoints


### Request Pin

    POST /api/v1/auth/pin 
    
    Request Body
    {
        "email": "some@email....",
        "ecdsaPubKey": "0x31..."
    }

    Response Body
    {
	    "accountId": "0000-0000-....",   // either the of the newly created account or if email associated with an account the according resp. id 
	    "email": "some@email....",
	    "ecdsaPubKey": "0x31...",
	    "expiration": 60
    }

NOTE: the pin is being sent via email

### Redeem Pin

    POST /api/v1/auth/pin/redeem
    
    Request Body
    {
        "pin": "123456",
        "pinSignature": "0x12481", // signed with the private key of the given public key
        "audiences": ["https://api.respurce.com"]     
    }

    Response Body
    JWT Token

NOTE: the pin is stored in a pin pool in memory. 