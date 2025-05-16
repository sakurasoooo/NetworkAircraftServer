# NetworkAircraftServer

NetworkAircraftServer is a Go-based server for a multiplayer aircraft game. It handles game logic, player connections, and real-time game state synchronization.

## Features

* **Real-time Multiplayer:** Supports multiple players connecting and interacting in a shared game world.
* **Game Entities:** Manages players, a boss entity, and projectiles (player rockets and boss rockets). [cite: 1]
* **Game Logic:**
    * Handles player movement and attacks. [cite: 1]
    * Controls boss behavior, including movement and attacks. [cite: 1]
    * Manages rocket movement, collision detection (rockets hitting the boss or players), and health updates. [cite: 1]
    * Resets the game when the boss is defeated. [cite: 1]
* **Network Communication:**
    * Uses TCP for client-server communication. [cite: 1]
    * Receives client requests (e.g., move, attack) in JSON format. [cite: 1]
    * Broadcasts game state updates to all connected clients in JSON format at regular intervals (every 200ms). [cite: 1]

## Architecture

The server is built in Go and is structured as follows:

* **`main.go`**: Contains the main server logic, including:
    * Listening for incoming TCP connections on `localhost:8080`. [cite: 1]
    * Handling new client connections and managing a list of connected clients. [cite: 1]
    * A game loop that updates the game state and broadcasts it to clients every 200 milliseconds. [cite: 1]
    * Functions to initialize and reset the game. [cite: 1]
    * Handlers for different client request types (`move`, `attack`, `hit`). [cite: 1]
    * Logic for updating positions and states of players, boss, and rockets. [cite: 1]
* **`game/` package**: Defines the structures and functions related to game entities and mechanics:
    * **`player.go`**: Defines the `Player` struct and related functions like `NewPlayer`. Players have attributes like health, attack, position, and a unique UUID.
    * **`boss.go`**: Defines the `Boss` struct and its behavior, including movement patterns and attack timers. The boss also has health, attack, position, and a UUID. It can generate random movement patterns and targets players for attacks.
    * **`playerRocket.go`**: Defines the `PlayerRocket` struct. Player rockets are created when a player attacks and move towards the boss. [cite: 1]
    * **`bossRocket.go`**: Defines the `BossRocket` struct. Boss rockets are created by the boss and target a random player. [cite: 1]
    * **`request.go`**: Defines the `ClientRequest` struct used for decoding messages from clients.
    * **`Vec2`**: A common struct likely used for representing 2D positions and vectors (defined in `boss.go` but used across game entities).

## Getting Started

### Prerequisites

* Go (refer to Go's official documentation for installation instructions)

### Running the Server

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <your-repository-url>
    cd NetworkAircraftServer-main
    ```
2.  **Run the server:**
    ```bash
    go run main.go
    ```
    The server will start and listen on `localhost:8080`. [cite: 1]

## How it Works

1.  The server starts and listens for TCP connections. [cite: 1]
2.  When a client connects, a new player is created and added to the game. [cite: 1] Initial player data (UUID, name, position) is sent to the client. [cite: 1]
3.  Clients can send requests to the server, such as:
    * `move`: To update their player's intended movement. [cite: 1]
    * `attack`: To make their player fire a rocket. [cite: 1]
4.  The server continuously updates the game state in a loop (every 200ms): [cite: 1]
    * Updates boss position and handles its attacks (spawning boss rockets). [cite: 1]
    * Updates player positions based on their `NextMove` and handles their attacks (spawning player rockets). [cite: 1]
    * Updates player rocket positions, checks for hits against the boss, and updates health accordingly. [cite: 1]
    * Updates boss rocket positions, checks for hits against players, and updates health accordingly. [cite: 1]
    * Removes players or rockets with health less than or equal to zero. [cite: 1]
    * If the boss's health is zero, it is removed, and the game can be reset. [cite: 1]
5.  After each game state update, the server broadcasts the current state of all relevant entities (boss, players, player rockets, boss rockets) to all connected clients as a JSON array. [cite: 1]

## Data Structures

### Sent to Client on Connection:

```json
{
  "uuid": 12345,
  "name": "PlayerName",
  "position": {"X": 1.0, "Y": 2.5}
}
```

### Client Requests (Example: Move):
```json
{
  "type": "move",
  "x": 0.5, // Horizontal movement input
  "y": -0.2, // Vertical movement input
  "uuid": 12345 // Player's UUID
}
```

### Broadcast Game State
```json
{
  "type": "player",
  "health": 80,
  "attack": 1, // This might be static or part of player's base stats
  "position": {"X": 5.0, "Y": -3.0},
  "uuid": 12345,
  "username": "PlayerName"
}
```
### Example for a boss
```json
{
  "type": "boss",
  "health": 95,
  "attack": 10,
  "position": {"X": 0.0, "Y": 0.0},
  "uuid": 67890
}
```
