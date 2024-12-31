# Reactor - Visual Flow State Manager

Reactor is a visual flow state management tool that allows users to create, manage, and visualize state machines through an intuitive web interface. It provides a drag-and-drop interface for creating states and defining transitions between them.

This is a WIP implementation of the architecture described in [Dynamic State Management: Composable State Machines with Primitive Operations](https://medium.com/picus-security-engineering/dynamic-state-management-composable-state-machines-with-primitive-operations-7580ccdd0a3d).

## Core Concept

The core concept centers on building states from chains of reusable primitive operations. New system behaviors can be created by composing existing primitives in different configurations, with minimal need for new code. When new functionality is required, developers need only implement new primitive operations, which then become available for use across all states.

### Key Architectural Components

- **Primitive Operations**: Foundational, reusable building blocks that perform single, well-defined operations
- **State Definitions**: Configuration-driven state definitions that combine primitive chains
- **Chain of Responsibility**: Two-level implementation for primitive chaining and state execution
- **Execution Context**: Shared context that flows between primitives and states

## Features

- **Visual Flow Editor**: Drag-and-drop interface for creating and managing states
- **State Management**: Create, update, and delete states with custom properties
- **Edge Management**: Define success/failure transitions between states
- **Primitive Actions**: Assign primitive actions to states with execution order
- **Persistent Storage**: All states and configurations are saved to SQLite database
- **Real-time Updates**: Changes are immediately reflected in the UI
- **Composable States**: Build complex workflows by combining primitive operations
- **Reusable Primitives**: Create new behaviors by composing existing primitives

## Technology Stack

### Backend (Go)
- **Web Framework**: Standard Go HTTP server
- **Database**: GORM with SQLite
- **Architecture**: Clean architecture with separated concerns

### Frontend (React + TypeScript)
- **UI Framework**: React
- **Flow Visualization**: React Flow
- **Type Safety**: TypeScript
- **Styling**: CSS-in-JS

## Project Structure
