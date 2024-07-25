# cute-little-authenticator

Toy project to learn how authenticator apps work.

## How it works basically

### The motivation
Using authenticator apps, I'm intrigued at how my app and the GitHub server seem to communicate with each other and agree on, very synchronously, 6 random digits which change every 30 seconds. Even more amazingly, my phone doesn't even have internet.

Filled with curiosity, I embarked on the journey to understand this black magic.

### The juice
So, very simply, this is how it works:

1. The server and our app share a secret key (the QR code you scan).
2. There's an algorithm that takes in the secret key and the current time, and spits out a 6-digit number. Given the same inputs, it always gives the same output.
3. When we authenticate, the server and the app calculate the number independently, using the shared secret key.
4. After we entered the numbers, the server checks if they match with the server's numbers. If they match, it means that we own the device that has the secret key -> we know the secret key, and the server will joyfully trust us and let us in.

The algorithm is called [TOTP](https://datatracker.ietf.org/doc/html/rfc6238#section-1.2), which is based on another called [HMAC](https://datatracker.ietf.org/doc/html/rfc2104). What's interesting about HMAC is that, it is created for a different purpose - to ensure that any message you receive over unsecure channel is not tampered with. But in this case, we use HMAC to hash the current time - a message both parties know about. So, HMAC is not used the way it is meant to be used?

I guess that in the end, the server just wants to make sure that the client knows the secret key, and among different ways to achieve this, HMAC is both efficient and user-friendly. For example, we can use a block cipher, but what would the message be? We want a code (or ciphertext) that is short, changes every 30 seconds, and maybe fixed-length so that attacker cannot take advantage when the code is shorter. That sounds like just the job for a hash function. But then we also want the hash to take in a secret key. So, HMAC seems like a good choice after all.