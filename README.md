# keylogger

A keylogger created for a BITSkrieg assignment.

## App

### Features

- Captures keyboard events globally
- Stores logs periodically inside a hidden directory in the user's home directory
- Works well on Windows

TODO: support Linux
TODO: add networking to upload logs to the server

## Server

Exposes an API endpoint to log batches of keys to a Firebase Firestore database

### API

#### `POST /api/log/save`

##### Required fields

`ts`: UNIX timestamp (in milliseconds)
`d`: Data, encoded in Base64 format

##### Optional fields

`c`: Clipboard data
`t`: Target ID. If it is not sent, then the server creates one and sends it back to the client, which is expected to be sent in subsequent requests from the client.

##### Example request body

```json
{
  "ts": 1731586746340,
  "d": "eyJfIjpbImEiLCJiIiwiYyIsImQiXX0=",
  "c": ["hello", "api"],
  "t": "14a901d8-e488-45e5-bdfe-2963317718f2"
}
```

##### Response

The server only responds in HTTP error codes. No response body is sent.

HTTP codes sent by the server are: `200 OK`, `400 Bad Request`, `500 Internal Server Error`

Body is sent only only when the client doesn't set `t` in its request to the server. Example response:

```json
{
  "t": "14a901d8-e488-45e5-bdfe-2963317718f2"
}
```
