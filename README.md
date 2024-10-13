# tcp-server
# TCP Client-Server Application

## Overview

This project is a TCP client-server application that simulates a high transaction system using a Redis-like storage mechanism. The client reads user records from the storage, sends them to a server, and receives unique transaction IDs in response. The system manages transactions per second (TPS) to simulate real-world load conditions.

## Features

- TCP client that communicates with a server over multiple connections.
- Redis-like storage that simulates user data storage and retrieval.
- Configurable transactions per second (TPS) based on a JSON configuration file.
- Concurrent handling of requests using Goroutines for efficient processing.

## Architecture

The application consists of three main components:

1. **Redis-like Storage**: An in-memory key-value store that simulates storing user records.
2. **TCP Client**: Manages the sending of requests to the server based on TPS settings.
3. **TCP Server**: Listens for incoming connections, processes requests, and sends back transaction IDs.

## Requirements

- Go 1.16 or later
- Basic understanding of Go and TCP networking

## Setup

### Clone the Repository

```bash
git clone https://github.com/your-username/tcp-client-server.git
cd tcp-client-server
