package ldtk

import (
	"log"
	"strings"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

// LoadTiles loads tiles for the specified level into the storage
func (p *LDtkProject) LoadTiles(levelName string, sto warehouse.Storage) error {
	level, exists := p.parsedLevels[levelName]
	if !exists {
		log.Printf("Level '%s' not found", levelName)
		return nil
	}

	// Process each tile layer
	layerIndex := 0
	for _, tileData := range level.TileRawData {
		// Find the tileset info
		var tilesetPath string
		var tileSize int
		for i := range p.Defs.Tilesets {
			if p.Defs.Tilesets[i].UID == tileData.TilesetDefUID {
				tilesetPath = p.Defs.Tilesets[i].RelPath
				tileSize = p.Defs.Tilesets[i].TileGridSize
				break
			}
		}

		if tilesetPath == "" {
			continue
		}

		index := strings.Index(tilesetPath, "tilesets")
		tilesetPath = "images/" + tilesetPath[index:]

		// Create a single entity with a sprite bundle for this layer
		archetype, err := sto.NewOrExistingArchetype(
			client.Components.SpriteBundle,
			spatial.Components.Position,
		)
		if err != nil {
			return err
		}

		// Create a sprite bundle for the whole layer
		blueprint := client.NewSpriteBundle().
			AddSprite(tilesetPath, true).
			WithPriority(10 + layerIndex). // Higher layers get higher priority
			WithOffset(vector.Two{X: 0, Y: 0})

		// Add all tiles to the TileSet
		for _, tile := range tileData.GridTiles {
			// Calculate source position in tiles
			sourceX := tile.Src[0] / tileSize
			sourceY := tile.Src[1] / tileSize

			// Extract flip flags
			flippedX := (tile.F & 1) != 0
			flippedY := (tile.F & 2) != 0

			// Add the tile to the tileset
			blueprint.Blueprints[0].TileSet = append(blueprint.Blueprints[0].TileSet, client.Tile{
				SourceX:  sourceX,
				SourceY:  sourceY,
				TileID:   tile.T,
				FlippedX: flippedX,
				FlippedY: flippedY,
				X:        float64(tile.Px[0]),
				Y:        float64(tile.Px[1]),
			})
		}

		// Create a single entity for the entire layer
		err = archetype.Generate(1,
			spatial.NewPosition(0, 0), // Origin position
			blueprint,
		)
		if err != nil {
			return err
		}

		layerIndex++
	}

	return nil
}
