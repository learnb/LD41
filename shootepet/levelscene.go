// Copyright 2018 Bryan Learn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shootepet

import (
        "fmt"
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	rshootepet "github.com/learnb/ld41/resources/images/shootepet"
)

/* Definitions & Initialization */

var (
    tilesImage *ebiten.Image
    pet *Entity
    owner *Entity
)

func init() {
        /* Create Characters */
        pet = &Entity{x: 32*12, y: 32*1, resoures: [3]float32{0.0, 0.0, 0.0}}
        owner = &Entity{x: 32*5, y: 32*5, resoures: [3]float32{0.0, 0.0, 0.0}}

        /* Preload Sprites */
	img, _, err := image.Decode(bytes.NewReader(rshootepet.Tiles_png))
	if err != nil {
		panic(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

        img, _, err = image.Decode(bytes.NewReader(rshootepet.Owner_png))
	if err != nil {
		panic(err)
	}
	owner.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        owner.setSizeByImage()

        img, _, err = image.Decode(bytes.NewReader(rshootepet.Pet_png))
	if err != nil {
		panic(err)
	}
	pet.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        pet.setSizeByImage()
}

type LevelScene struct {
	count int
}

// Map Info

const (
        tileSize = 32
        tileXNum = 8
)

var (
        tilelayers = [][]int{ /* ints represent index of reshootepet.Tiles_png */
                /* First Layer - Main */
                {
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 1, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 3, 3, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 3, 3, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 4, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 1, 6, 10,
                    10, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 10,
                },
                /* Second Layer - Top */
               /* {
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, 7, 7, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, 7, 7, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
                },*/
        }
        tileprops = [][]int{ /* ints represent tile property booleans */
                /* Tile properties:
                    0: passable
                    1: impassable
                    2: bullet-passable
                    3: character-passable
                */
                {
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                },
        }
)

/* Update */

func (s *LevelScene) Update(state *GameState) error {
        if state.Input.TriggeredSecondary() {
                state.SceneManager.GoTo(&GameOverScene{})
                return nil
        }
        if state.Input.TriggeredMain() {
                pet.moveToCell(1,1)
                return nil
        }

        // update owner
        s.updateOwner(state)

        // update pet
        s.updatePet(state)

        // check collision
        if owner.doesCollideWith(pet) {
                state.SceneManager.GoTo(&GameOverScene{})
                return nil
        }

	return nil
}

func (s *LevelScene) updateOwner(state *GameState) error {
        if state.Input.StateForUp() > 0 {
                owner.y -= 3
        }
        if state.Input.StateForDown() > 0 {
                owner.y += 3
        }
        if state.Input.StateForLeft() > 0 {
                owner.x -= 3
        }
        if state.Input.StateForRight() > 0 {
                owner.x += 3
        }

        return nil
}


func (s *LevelScene) updatePet(state *GameState) error {
        pet.x -= 1
        pet.y += 1

        return nil
}

/* Draw */

func (s *LevelScene) Draw(r *ebiten.Image) {
        /* Debug */
        ebitenutil.DebugPrint(r, "\nNothing here yet :(")
	message := "~ Level Scene ~"
	x := 0
	y := ScreenHeight - 48
	drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)

        /* Draw Map */
        s.drawMap(r)

        /* Draw Characters */
        s.drawChars(r)

	message = fmt.Sprintf("Dist: %0.2f", owner.distanceTo(pet))
	drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)

        ebitenutil.DebugPrint(r, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

func (s *LevelScene) drawMap(r *ebiten.Image) {
        const xNum = ScreenWidth / tileSize
        for _, l := range tilelayers {
                for i, t := range l {
                        op := &ebiten.DrawImageOptions{}
                        op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))

                        sx := (t % tileXNum) * tileSize
                        sy := (t / tileXNum) * tileSize
                        rect := image.Rect(sx, sy, sx+tileSize, sy+tileSize)
                        op.SourceRect = &rect
                        r.DrawImage(tilesImage, op)
                }
        }
}

func (s *LevelScene) drawChars(r *ebiten.Image) {
        op := &ebiten.DrawImageOptions{}
        //x, y := owner.pos()
        //op.GeoM.Translate(float64(x), float64(y))
        op.GeoM.Translate(owner.posf())
        r.DrawImage(owner.image, op)

        op = &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(pet.x), float64(pet.y))
        r.DrawImage(pet.image, op)
}

/* Other */

func NewGameScene() *LevelScene {
        return &LevelScene{
            count: 0,
        }
}

