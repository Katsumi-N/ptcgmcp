# Pokemon TCG Management Card Project (PTCGMCP)

PTCGMCP is a server application that provides APIs for retrieving Pokemon card information and managing decks for the Pokemon Trading Card Game.

## Project Overview

This project consists of two main components:

1. **API Server**: Handles HTTP requests and provides JSON responses for card searches, card details, and deck management.
2. **MCP (Model Context Protocol) Server**: Manages the data models and business logic for the application.

## Features

### Card Information
- Search for Pokemon cards by name, type, and other parameters
- Get detailed information about specific cards:
  - Pokemon cards (stats, abilities, moves)
  - Trainer cards
  - Energy cards

### Deck Management
- Create custom decks with a mix of Pokemon, Trainer, and Energy cards
- List all saved decks
- View detailed information about a specific deck
- Edit existing decks (rename, change cards, adjust quantities)
- Delete decks
- Validate decks against game rules

## API Endpoints

### Card Information Endpoints
- `GET /v1/cards/search?q={query}&card_type={type}` - Search for cards by name and type
- `GET /v1/cards/detail/pokemon/{id}` - Get details about a specific Pokemon card
- `GET /v1/cards/detail/trainer/{id}` - Get details about a specific Trainer card
- `GET /v1/cards/detail/energy/{id}` - Get details about a specific Energy card

### Deck Management Endpoints
- `GET /v1/decks` - List all decks
- `GET /v1/decks/detail/{id}` - Get details about a specific deck
- `POST /v1/decks/create` - Create a new deck
- `POST /v1/decks/validate` - Validate a deck against game rules
- `POST /v1/decks/edit/{id}` - Edit an existing deck
- `DELETE /v1/decks/delete/{id}` - Delete a deck

## Technology Stack

- Backend:
  - Go (API server)
  - Kotlin with Ktor framework (MCP server)
- Databases:
  - MySQL for persistent storage
  - Meilisearch for fast card searching

## Getting Started

### Running the API Server
1. Navigate to the `api` directory
2. Configure your database settings in the config file
3. Run `go run cmd/main.go`

### Running the MCP Server
1. Navigate to the `mcp` directory
2. Run one of the following commands:
   - `./gradlew run` - Run the server directly
   - `./gradlew buildFatJar` - Build an executable JAR with all dependencies
   - `./gradlew runDocker` - Run using a local Docker image

### Testing the API
You can use the provided HTTP request examples in `/api/request.http` to test the API endpoints.

## Project Structure
- `/api` - Go API server
- `/mcp` - Kotlin MCP server

## License
This project is proprietary software.