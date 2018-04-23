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
        "math/rand"
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
    mapGraph *Graph

    pet *Entity
    petPath []int
    petTarget int

    owner *Entity

    hBullet *Bullet
    aBullet *Bullet
    eBullet *Bullet
    activeWeapon int
)



func init() {
        /* Create Characters */
        pet = &Entity{x: 32*12, y: 32*1, resoures: [3]float64{0.0, 0.0, 0.0}, speed: 3.0}
        owner = &Entity{x: 32*5, y: 32*5, resoures: [3]float64{0.0, 0.0, 0.0}, speed: 1.5}
        activeWeapon = 0

        /* Create Weapons & Ammo */
        hBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resoures: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            active: false,
        }
        aBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resoures: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            active: false,
        }
        eBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resoures: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            active: false,
        }

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

        // Hunger Bullet
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Bullet_png))
	if err != nil {
		panic(err)
	}
	hBullet.ent.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        hBullet.ent.setSizeByImage()

        // Attention Bullet
        img, _, err = image.Decode(bytes.NewReader(rshootepet.BulletA_png))
	if err != nil {
		panic(err)
	}
	aBullet.ent.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        aBullet.ent.setSizeByImage()

        // Exercise Bullet
        img, _, err = image.Decode(bytes.NewReader(rshootepet.BulletE_png))
	if err != nil {
		panic(err)
	}
	eBullet.ent.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        eBullet.ent.setSizeByImage()




        /* Define Graph */
        tileprops = buildTileProps()
        mapGraph = &Graph{graph: tileprops, xMax: 15, yMax: 10}
        buildCollisionMap()

        /* Get Pet moving */
        petPath = makeNewPetPath(0, 0) //Test inital fails safely
        x, y := mapGraph.Indx2Coord(petTarget)
        fmt.Printf("Pet target: (%d, %d)\n", x,y)
        fmt.Printf("path length: %d\n", len(petPath))
        for _, v := range petPath {
            x, y := mapGraph.Indx2Coord(v)
            fmt.Printf("(%d, %d)\n", x,y)
        }

        // test pathfinding impassable
        testl := mapGraph.getNeighbors(1)
        fmt.Printf("Neighbors of 1: %d %d %d %d\n", testl[0], testl[1], testl[2], testl[3])
        fmt.Printf("Neighbors real: %d %d %d %d\n", -1, 2, -1, 16)
}

func  (s *LevelScene) Init() {
    //ebiten.init() //TODO have Scenes implement Init func so they can be called again
}

type LevelScene struct {
	count int
}

type Bullet struct {
        ent Entity
        vec [2]float64
        target [2]float64
        active bool
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
                    9, 1, 9, 9, 9, 9, 9, 9, 9, 1, 1, 1, 9, 9, 9,
                    9, 1, 5, 5, 5, 5, 5, 5, 5, 5, 5, 1, 5, 1, 9,
                    9, 1, 5, 5, 5, 5, 5, 5, 5, 5, 5, 1, 1, 1, 9,
                    9, 1, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 9,
                    9, 1, 5, 5, 5, 5, 5, 3, 3, 5, 5, 5, 5, 5, 9,
                    9, 1, 5, 5, 5, 5, 5, 3, 3, 5, 5, 5, 5, 5, 9,
                    9, 5, 5, 5, 5, 5, 5, 5, 5, 3, 5, 5, 5, 5, 9,
                    1, 1, 1, 1, 1, 4, 5, 5, 5, 5, 3, 5, 5, 5, 9,
                    9, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 3, 1, 1, 9,
                    9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
                },
        }
        tileprops = []int { /* ints represent tile property booleans */
                /* Tile properties:
                    0: passable
                    1: impassable
                    2: bullet-passable
                    3: character-passable
                */
                }
        collisionMap = map[int]bool{}
)

func buildCollisionMap() {
        for i, v := range tileprops {
                if v != 0 {
                        collisionMap[i] = true
                }
        }
}

func buildTileProps() []int {
        a := make([]int, len(tilelayers[0]))
        for i, v := range tilelayers[0] {
                if v == 1 || v == 4 || v == 3  { // if node is impassable
                        a[i] = 1
                } else {
                        a[i] = 0
                }
        }
        return a
}

/* Update */

func (s *LevelScene) Update(state *GameState) error {
        // check of game over
        //if state.Input.TriggeredSecondary() {
        //        state.SceneManager.GoTo(&GameOverScene{})
        //        return nil
        //}


        //if state.Input.TriggeredMain() {
        //        pet.moveTowardCell(1,1)
        //        return nil
        //}

        // update owner
        s.updateOwner(state)  // handles user input

        // update pet
        s.updatePet(state)

        // update bullets
        s.updateBullets(state)

        // check collision
        if hBullet.ent.doesCollideWith(pet) {
                state.SceneManager.GoTo(&GameOverScene{})
                return nil
        }


	return nil
}

func (s *LevelScene) updateOwner(state *GameState) error {
        // movement input
        targetX, targetY := owner.pos()

        if state.Input.StateForUp() > 0 {
                targetY -= owner.speed
        }
        if state.Input.StateForDown() > 0 {
                targetY += owner.speed
        }
        if state.Input.StateForLeft() > 0 {
                targetX -= owner.speed
        }
        if state.Input.StateForRight() > 0 {
                targetX += owner.speed
        }

        owner.moveTowardPoint(targetX, targetY)

        // action input
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) { // shoot
                x, y := ebiten.CursorPosition()
                switch activeWeapon {
                case 0:
                    if !hBullet.active { // cannot shoot until bullet 'dies'
                        FireAt(x, y) //units in screen pixels
                    }
                case 1:
                    if !aBullet.active { // cannot shoot until bullet 'dies'
                        FireAt(x, y) //units in screen pixels
                    }
                case 2:
                    if !eBullet.active { // cannot shoot until bullet 'dies'
                        FireAt(x, y) //units in screen pixels
                    }

                }
        }
        if state.Input.TriggeredSecondary() { // change weaspon
                RotateWeapons()
        }

        return nil
}


func (s *LevelScene) updatePet(state *GameState) error {
        px, py := mapGraph.Indx2Coord(petTarget)
        pet.moveTowardCell(px, py)

        if pet.isAtCell(px, py) {
              // pop next cell from petPath
              if len(petPath) >= 1 {
                  petTarget, petPath = petPath[len(petPath)-1], petPath[:len(petPath)-1]
              } else { // make new path
                  // find new random destination
                  dx, dy := rand.Intn(mapGraph.xMax-1), rand.Intn(mapGraph.yMax-1)
                  //fmt.Printf("New Point: (%d, %d)\n", dx, dy)
                  petPath = makeNewPetPath(dx, dy)
              }
        }

        return nil
}

func (s *LevelScene) updateBullets(state *GameState) error {
        if hBullet.active {
            if hBullet.ent.isOnScreen() {   // not out of bounds yet
                hBullet.ent.moveByVecComponents(hBullet.vec[0], hBullet.vec[1])
            } else {                                                    // at target; clear
                hBullet.ent.x, hBullet.ent.y = owner.centerPos()
                hBullet.active = false
            }
        } else {  // keep bullet on owner
            hBullet.ent.x, hBullet.ent.y = owner.centerPos()
        }
        //
        if aBullet.active {
            if aBullet.ent.isOnScreen() {   // not out of bounds yet
                aBullet.ent.moveByVecComponents(aBullet.vec[0], aBullet.vec[1])
            } else {                                                    // at target; clear
                aBullet.ent.x, aBullet.ent.y = owner.centerPos()
                aBullet.active = false
            }
        } else {  // keep bullet on owner
            aBullet.ent.x, aBullet.ent.y = owner.centerPos()
        }
        //
        if eBullet.active {
            if eBullet.ent.isOnScreen() {   // not out of bounds yet
                eBullet.ent.moveByVecComponents(eBullet.vec[0], eBullet.vec[1])
            } else {                                                    // at target; clear
                eBullet.ent.x, eBullet.ent.y = owner.centerPos()
                eBullet.active = false
            }
        } else {  // keep bullet on owner
            eBullet.ent.x, eBullet.ent.y = owner.centerPos()
        }

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

        /* Draw Bullets */
        s.drawBullets(r)

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
        op.GeoM.Translate(owner.pos())
        r.DrawImage(owner.image, op)

        op = &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(pet.x), float64(pet.y))
        r.DrawImage(pet.image, op)
}

func (s *LevelScene) drawBullets(r *ebiten.Image) {
        switch activeWeapon {
        case 0:
            if hBullet.active {
                op := &ebiten.DrawImageOptions{}
                op.GeoM.Translate(hBullet.ent.pos())
                r.DrawImage(hBullet.ent.image, op)
            }
        case 1:
            if aBullet.active {
                op := &ebiten.DrawImageOptions{}
                op.GeoM.Translate(aBullet.ent.pos())
                r.DrawImage(aBullet.ent.image, op)
            }
        case 2:
            if eBullet.active {
                op := &ebiten.DrawImageOptions{}
                op.GeoM.Translate(eBullet.ent.pos())
                r.DrawImage(eBullet.ent.image, op)
            }
        }
}

/* Other */

func NewGameScene() *LevelScene {
        return &LevelScene{
            count: 0,
        }
}

func Point2MapCell(x, y float64) (int, int) {
    cX := int(x/tileSize)
    cY := int(y/tileSize)
    return cX, cY
}

/* Pet Helper Functions */

func makeNewPetPath(dstX, dstY int) []int {
    px, py := pet.pos()
    cx, cy := Point2MapCell(px, py)

    // check if dst is impassable
    if collisionMap[mapGraph.Coord2Indx(dstX, dstY)] { // if collision, change dst to a neighboring tile
        nl := mapGraph.getNeighbors(mapGraph.Coord2Indx(dstX, dstY))
        for _, v := range nl { // use first available neighbor
            if v != -1 {
                dstX, dstY = mapGraph.Indx2Coord(v)
            }
        }
    }

    l := mapGraph.Astar( mapGraph.Coord2Indx(cx, cy), mapGraph.Coord2Indx(dstX, dstY) )
    if len(l) >= 1 {
        _, l = l[len(l)-1], l[:len(l)-1]            // pop src (current) node
    }
    if len(l) >= 1 {
        petTarget, l = l[len(l)-1], l[:len(l)-1]    // pop & set next node
    }
    return l
}


/* Owner Helper Functions */

func FireAt(x, y int) {
    switch activeWeapon {
    case 0:
        hBullet.target[0] = float64(x)
        hBullet.target[1] = float64(y)
        hBullet.vec[0], hBullet.vec[1] = hBullet.ent.getVecComponents(hBullet.target[0], hBullet.target[1])
        hBullet.active = true
    case 1:
        aBullet.target[0] = float64(x)
        aBullet.target[1] = float64(y)
        aBullet.vec[0], aBullet.vec[1] = aBullet.ent.getVecComponents(aBullet.target[0], aBullet.target[1])
        aBullet.active = true
    case 2:
        eBullet.target[0] = float64(x)
        eBullet.target[1] = float64(y)
        eBullet.vec[0], eBullet.vec[1] = eBullet.ent.getVecComponents(eBullet.target[0], eBullet.target[1])
        eBullet.active = true
    }
}

func RotateWeapons() {
    if activeWeapon < 2 {
        activeWeapon += 1
    } else {
        activeWeapon = 0
    }
}
