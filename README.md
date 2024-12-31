# Reactor - Visual Flow State Manager

Reactor is a visual flow state management tool that allows you to create, manage, and visualize state machines through an intuitive web interface. Built with React Flow and Go, it provides a powerful way to design and execute complex state workflows.

## Features

- 🎨 Visual drag-and-drop interface for state machine design
- 🔄 Success/failure transition paths for each state
- 🧩 Composable primitive operations
- 💾 Persistent storage with SQLite
- 🔌 RESTful API for state management
- ⚡ Real-time state updates

## Architecture

Reactor implements a composable state machine architecture where:
- States are composed of primitive operations
- Each state has success and failure transitions
- Business logic is isolated in primitive operations
- State flow is configuration-driven

### Tech Stack

- Frontend: React + TypeScript + React Flow
- Backend: Go + Gorilla Mux
- Database: SQLite + GORM

## Getting Started

### Prerequisites

- Go 1.20+
- Node.js 16+
- npm or yarn

### Installation

1. Clone the repository:
```
git clone https://github.com/aliatli/reactor.git
cd reactor
```
2. Install&Run backend:
```
go mod download
go run cmd/web/main.go
```
3. Install&Run frontend:
```
cd web
npm install
npm run dev
```

4. Open your browser and navigate to `http://localhost:5173`

## Usage

1. **Creating States**
   - Click "Add New State" to create a new state
   - Enter a name for the state
   - Click "Create"

2. **Configuring Primitives**
   - Click on a state to open the primitive panel
   - Select primitives to be executed in that state
   - Click "Save" to update the state

3. **Creating Transitions**
   - Drag from a state's success (green) or failure (red) handle to another state
   - Transitions are automatically saved

4. **Saving the Flow**
   - Click "Save Flow" to persist the entire state machine
### Project Structure
```
├── cmd/
│ └── web/ # Application entry point
├── internal/
│ ├── api/ # HTTP handlers and routing
│ ├── core/ # Core domain types
│ ├── db/ # Database operations
│ ├── executor/ # State machine execution
│ └── models/ # Database models
├── examples/
│ └── primitives/ # Example primitive operations
└── web/ # Frontend React application
```