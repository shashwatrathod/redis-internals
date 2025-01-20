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
- TCP server to handle incoming connections
- Command-line flag configuration for server settings

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

## Technologies Used

- **Go**: The primary programming language used for this project.
- **RESP**: REdis Serialization Protocol, used for encoding and decoding data.

## References

- **DiceDb**: This project is inspired by [DiceDb](https://github.com/dicedb/dicedb), a lightweight in-memory database.
- **Redis Documentation**: For more information on Redis and its features, refer to the [Redis documentation](https://redis.io/documentation).
