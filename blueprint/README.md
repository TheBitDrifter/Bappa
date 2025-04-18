# Blueprint

Blueprint serves as the base layer of the Bappa Framework, defining common/shared component types and interfaces that other packages in the framework build upon. It establishes shared patterns and data structures used throughout the ecosystem.

Previously located [here](https://github.com/TheBitDrifter/blueprint) prior to monorepo (archived for Git history).

## Features

- **Core Interfaces**: Defines essential interfaces like `Scene` and `CoreSystem`
- **Common Components**: Provides standardized component definitions across multiple domains
- **Vector Mathematics**: Complete 2D vector implementation with comprehensive operations
- **Predefined Queries**: Ready-to-use queries for common component combinations
- **Background Utilities**: Tools for creating static and parallax backgrounds

## Installation

```bash
go get github.com/TheBitDrifter/bappa/blueprint
```

## Component Subpackages

Blueprint organizes components by domain:

- **blueprint/client**: Visual and audio components
  - `SpriteBundle`: Sprites with animation support
  - `SoundBundle`: Audio resources
  - `ParallaxBackground`: Multi-layered scrolling backgrounds

- **blueprint/input**: User interaction components
  - `InputBuffer`: Collection and management of user inputs
  - `StampedInput`: Inputs with timing and position information
- **blueprint/vector**: 2D vector mathematics

  - `Two`: Vector with extensive operations (add, subtract, rotate, etc.)
  - Vector interfaces for flexible implementation

## Quick Start

### Using Predefined Queries

```go
// Create a cursor for entities with position components
cursor := scene.NewCursor(blueprint.Queries.Position)

// Process matching entities
for range cursor.Next() {
    pos := spatial.Components.Position.GetFromCursor(cursor)
    // Process entity...
}
```

### Creating Backgrounds

```go
// Create a parallax background with multiple layers
builder := blueprint.NewParallaxBackgroundBuilder(storage)
builder.AddLayer("backgrounds/mountains.png", 0.2, 0.0) // Slow-moving background
builder.AddLayer("backgrounds/clouds.png", 0.5, 0.1)    // Mid-speed layer
builder.WithOffset(vector.Two{X: 0, Y: 20})             // Optional offset
builder.Build()

// Create a static background
blueprint.CreateStillBackground(storage, "backgrounds/scene.png")
```

## License

MIT License - see the [LICENSE](LICENSE) file for details.
