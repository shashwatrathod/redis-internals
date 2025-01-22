# Redis Internals

Redis Internals is a Go project that mimics the functionality of Redis, inspired by DiceDb. This project aims to provide a basic implementation of Redis-like features, including RESP (REdis Serialization Protocol) decoding, TCP server handling, and command-line flag configuration. It is a learning project inspired by DiceDb and aims to provide a basic understanding of Redis-like functionality in Go.

## Table of Contents

- Features
- Installation
- Usage
- Technologies Used
- References

## Features

- RESP decoding for various data types (integers, strings, bulk strings, and arrays)
- Error handling
- See [supported commands](#supported-commands) for a list of available commands.

## Installation

To install and run the project, follow these steps:

1. Clone the repository:

   ```sh
   git clone https://github.com/shashwatrathod/redis-internals.git
   cd redis-internals
   ```

2. Install the dependencies:
    ```sh
    go mod download
    ```

3. Build the project
    ```sh
    go build -o redis-internals
    ```

## Usage

To run the Redis Internals server, use the following command:

    ```sh
    redis-internals -host <host> -port <port>
    ```

This project is meant to be a drop-in replacement for an actual Redis server. So you can interact
with the server using the `redis-cli` utility.

```sh
redis-cli -h <host> -p <port> PING hello
```
See [supported commands](#supported-commands) for a list of available commands to try!

## Supported Commands

- [PING](https://redis.io/docs/latest/commands/ping/)

## Technologies Used

- **Go**: The primary programming language used for this project.
- **RESP**: REdis Serialization Protocol, used for encoding and decoding data.

## References

- **DiceDb**: This project is inspired by [DiceDb](https://github.com/dicedb/dicedb), a lightweight in-memory database.
- **Redis Documentation**: For more information on Redis and its features, refer to the [Redis documentation](https://redis.io/documentation).
