<div align="center">
<img title="Brick Hauler" width="250" alt="Brick Hauler Gopher" src="./assets/brikhaulergoper.png">
</div>

# BrickHauler

A tool designed for fast and easy load testing of websites and web applications.

## Params

Parameters accepted by the command line:

- `--verb` (string): Specifies the HTTP verb to be used (GET, POST, PUT, PATCH, DELETE, etc.).
- `--uri` (string): The URL where the tests will be performed (e.g., <https://example.com>).
- `--concurrent` (int): The number of virtual users to launch requests concurrently.
- `--request` (int): The total number of requests to be sent by all users.
- `--cookie` (string): Cookie to be included in the requests (format: cookieName=cookieValue).
- `--proxy` (string): Url to the proxy that is going to take all the request.
- `--feed` (bool): Display real-time logs of the test.

## Usage

Here's the example you're probably looking for to understand how this tool works:

```bash
go run ./cmd/brickhauler --uri https://example.com --concurrent 2 --request 4 --feed
```

Do you want to use another http verb?

```bash
go run ./cmd/brickhauler --verb POST --uri https://example.com --concurrent 2 --request 4 --feed
```

Need to add a cookie? Here's how:

```bash
go run ./cmd/brickhauler --uri https://example.com --concurrent 2 --request 4 --cookie "foo=bar" --feed
```

Need to add multiple cookies? It's just as simple:

```bash
go run ./cmd/brickhauler --uri https://example.com --concurrent 2 --request 4 --cookie "foo=bar" --cookie "root=toor" --feed
```

Maybe you need a proxy:

```bash
go run ./cmd/brickhauler --uri https://example.com --concurrent 2 --request 4 --proxy "http://43.123.54.1:8080/" --feed
```

## Features

- Ability to choose the HTTP method for making requests.

- Simulation of virtual users acting independently, capable of making concurrent requests.

- Option to add cookies to the requests.

- Use proxies for doing all the requests.

## Building

To generate binaries for the main operating systems, we must execute the following commands:

```bash
chmod +x ./build.sh
./build.sh ./cmd/brickhauler
```

Or to build a single binary:

```bash
go build -o ./bin/brickhauler ./cmd/brickhauler
```
