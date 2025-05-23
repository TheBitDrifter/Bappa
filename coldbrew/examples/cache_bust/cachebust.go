package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"

	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
)

//go:embed assets/*
var assets embed.FS

const (
	sceneOneName = "s1"
	sceneTwoName = "s2"
)

// This example is setup so the scene change triggers the cache protocol
func main() {
	client := coldbrew.NewClient(
		640,
		360,
		5, 10,
		10,
		assets,
	)

	client.SetTitle("Scene Managment Example")
	client.SetMinimumLoadTime(20)

	err := client.RegisterScene(
		sceneOneName,
		640,
		360,
		sceneOnePlan,
		[]coldbrew.RenderSystem{},
		[]coldbrew.ClientSystem{},
		[]blueprint.CoreSystem{},
	)
	if err != nil {
		log.Fatal(err)
	}
	err = client.RegisterScene(
		sceneTwoName,
		640,
		360,
		sceneTwoPlan,
		[]coldbrew.RenderSystem{},
		[]coldbrew.ClientSystem{},
		[]blueprint.CoreSystem{},
	)
	if err != nil {
		log.Fatal(err)
	}

	client.RegisterGlobalRenderSystem(coldbrew_rendersystems.GlobalRenderer{})
	client.RegisterGlobalClientSystem(basicTransferSystem{}, &coldbrew_clientsystems.CameraSceneAssignerSystem{})

	client.ActivateCamera()

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}

func sceneOnePlan(width, height int, sto warehouse.Storage) error {
	spriteArchetype, err := sto.NewOrExistingArchetype(
		spatial.Components.Position,
		client.Components.SpriteBundle,
		client.Components.CameraIndex,
	)
	if err != nil {
		return err
	}

	err = spriteArchetype.Generate(1,
		input.Components.ActionBuffer,

		spatial.NewPosition(255, 20),
		client.NewSpriteBundle().
			AddSprite("images/sprite.png", true),

		client.CameraIndex(0),
	)
	err = blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/sky.png", 0.1, 0.1).
		AddLayer("images/far.png", 0.3, 0.3).
		AddLayer("images/mid.png", 0.4, 0.4).
		AddLayer("images/near.png", 0.8, 0.8).
		Build()
	if err != nil {
		return err
	}
	return nil
}

func sceneTwoPlan(width, height int, sto warehouse.Storage) error {
	err := blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/sky2.png", 0.1, 0.1).
		Build()
	if err != nil {
		return err
	}
	return nil
}

type transfer struct {
	target       coldbrew.Scene
	playerEntity warehouse.Entity
}

type basicTransferSystem struct{}

func (basicTransferSystem) Run(cli coldbrew.Client) error {
	var pending []transfer
	sceneCache := cli.Cache()

	for activeScene := range cli.ActiveScenes() {
		if !activeScene.Ready() {
			continue
		}
		cursor := activeScene.NewCursor(blueprint.Queries.CameraIndex)
		for range cursor.Next() {
			if inpututil.IsKeyJustPressed(ebiten.Key1) {

				currentPlayerEntity, err := cursor.CurrentEntity()
				if err != nil {
					return err
				}

				// Simple toggle between scenes
				var sceneTargetName string
				if activeScene.Name() == sceneOneName {
					sceneTargetName = sceneTwoName
				} else {
					sceneTargetName = sceneOneName
				}

				targetSceneIndex, found := sceneCache.GetIndex(sceneTargetName)
				if !found {
					log.Println("Target scene not found:", sceneTargetName)
					return fmt.Errorf("target scene '%s' not found in cache", sceneTargetName)
				}

				targetScene := sceneCache.GetItem(targetSceneIndex)

				transfer := transfer{
					target:       targetScene,
					playerEntity: currentPlayerEntity,
				}
				pending = append(pending, transfer)
			}
		}
	}

	for _, transfer := range pending {
		cli.ActivateScene(transfer.target, transfer.playerEntity)
	}
	for activeScene := range cli.ActiveScenes() {
		cursor := activeScene.NewCursor(blueprint.Queries.CameraIndex)
		if cursor.TotalMatched() == 0 {
			// No player entities with input buffers left in this scene
			cli.DeactivateScene(activeScene)
		}
	}

	return nil
}
