# Muted Websocket

**Muted Websocket** *is a backend service built in Go that bridges communication between muted hypixel players to non-muted players in real-time.*

## Deploy
- **Solution to deployment**: I personally use [Cloud Run](https://console.cloud.google.com/run) and the provided docker file to make deployment easy.

## Authentication (Thank you [appable](https://github.com/appable0))
- Client adds header with generated server id (from mod) and requests to join the server id.
- Server checks if player did join serverid.
- Once authenticated, server generates a session id which gets stored on a map with the client as the index.

## Credits
- [appable0](ttps://github.com/appable0) Guidance on security between client and server
- [NEU Ursa-Minor](https://github.com/NotEnoughUpdates/ursa-minor/tree/master) Implementation on authentication

## License
---
This project is licensed under the MIT license.
