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
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	rshootepet "github.com/learnb/ld41/resources/images/shootepet"
)

/* Definitions & Initialization */

var (
    uiSqr *ebiten.Image
    blissfulImage *ebiten.Image
    happyImage *ebiten.Image
    concernedImage *ebiten.Image
    woefulImage *ebiten.Image
    tilesImage *ebiten.Image
    mapGraph *Graph

    pet *Entity
    petPath []int
    petTarget int
    emotion int /* 0: Blissful -> 3: Woeful */
    desireRate [3]float64
    petBall bool // is pet playing ball
    ballLanded bool // is ball on the ground

    owner *Entity

    hBullet *Bullet
    aBullet *Bullet
    eBullet *Bullet
    activeWeapon int
)

func init() {
    myInit()
}

func myInit() {
        /* Create Characters */
        pet = &Entity{x: 32*12, y: 32*1, resources: [3]float64{1.0, 1.0, 1.0}, speed: 3.0}
        desireRate = [3]float64{0.001, 0.001, 0.001}
        emotion = 0 // start blissful
        petBall = false
        ballLanded = false

        owner = &Entity{x: 32*5, y: 32*5, resources: [3]float64{0.0, 0.0, 0.0}, speed: 1.5}
        activeWeapon = 0

        /* Create Weapons & Ammo */
        hBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resources: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            count: 0,
            active: false,
        }
        aBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resources: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            count: 0,
            active: false,
        }
        eBullet = &Bullet{
            ent: Entity{x: 0, y: 0, resources: [3]float64{0.0, 0.0, 0.0}, speed: 15},
            vec: [2]float64{0,0},
            target: [2]float64{0,0},
            count: 0,
            active: false,
        }

        /* Preload Sprites */
	img, _, err := image.Decode(bytes.NewReader(rshootepet.Tiles_png))
	if err != nil {
		panic(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

        // owner
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Owner_png))
	if err != nil {
		panic(err)
	}
	owner.image, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        owner.setSizeByImage()

        // pet
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Blissful_png))
	if err != nil {
		panic(err)
	}
	blissfulImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Happy_png))
	if err != nil {
		panic(err)
	}
	happyImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Concerned_png))
	if err != nil {
		panic(err)
	}
	concernedImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        img, _, err = image.Decode(bytes.NewReader(rshootepet.Woeful_png))
	if err != nil {
		panic(err)
	}
	woefulImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
        //pet.setSizeByImage()
        pet.w, pet.h = 32, 32

        // UI
        uiSqr, _ = ebiten.NewImage(65, 35, ebiten.FilterDefault)

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
        //x, y := mapGraph.Indx2Coord(petTarget)
        //fmt.Printf("Pet target: (%d, %d)\n", x,y)
        //fmt.Printf("path length: %d\n", len(petPath))
        //for _, v := range petPath {
        //    x, y := mapGraph.Indx2Coord(v)
        //    fmt.Printf("(%d, %d)\n", x,y)
        //}

        // test pathfinding impassable
        //testl := mapGraph.getNeighbors(1)
        //fmt.Printf("Neighbors of 1: %d %d %d %d\n", testl[0], testl[1], testl[2], testl[3])
        //fmt.Printf("Neighbors real: %d %d %d %d\n", -1, 2, -1, 16)
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
        // check for slow rendering
        if ebiten.IsRunningSlowly(){
            return nil    // Skip
        }

        // check for game over
        if emotion == 4 {
        //if state.Input.TriggeredSecondary() {
                state.SceneManager.GoTo(&GameOverScene{})
                myInit()
                return nil
        }


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
        //if hBullet.ent.doesCollideWith(pet) {
        //        state.SceneManager.GoTo(&GameOverScene{})
        //        return nil
        //}


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

        // wall collision - check target is not a wall
        // cheat by bringing in bounding box a bit
        buff := 3.0
        wallStop := false
        cx, cy := Point2MapCell(targetX+buff, targetY+buff)
        if collisionMap[mapGraph.Coord2Indx(cx, cy)] { // top left
            wallStop = true
        }
        cx, cy = Point2MapCell(targetX+buff, targetY+float64(owner.h)-buff)
        if collisionMap[mapGraph.Coord2Indx(cx, cy)] { // bottom left
            wallStop = true
        }
        cx, cy = Point2MapCell(targetX+float64(owner.w)-buff, targetY+buff)
        if collisionMap[mapGraph.Coord2Indx(cx, cy)] { // top right
            wallStop = true
        }
        cx, cy = Point2MapCell(targetX+float64(owner.w)-buff, targetY+float64(owner.h)-buff)
        if collisionMap[mapGraph.Coord2Indx(cx, cy)] { // botttom right
            wallStop = true
        }

        if !wallStop {
            owner.moveTowardPoint(targetX, targetY)
        } else { // do nothing
            //owner.moveTowardPoint(-1*(targetX), -1*(targetY))
        }


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

        // movement
        if ballLanded {    // chase ball
            ballLanded = false
            bx, by := eBullet.ent.cellPos()
            //fmt.Printf("Pet Ball: making path\n")
            petPath = makeNewPetPath(bx, by)            // make path to ball
            petBall = true         // done until ball found or despawned
            px, py := mapGraph.Indx2Coord(petTarget)
            pet.moveTowardCell(px, py)
        } else {        // random movement
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
        }

        // collision
        if pet.doesCollideWith(&hBullet.ent) && hBullet.active {       // food get
            fmt.Println("hit H")
            pet.resources[0] += desireRate[0]*1000
            hBullet.active = false
        }
        if pet.doesCollideWith(&aBullet.ent) && aBullet.active {       // love get
            fmt.Println("hit A")
            pet.resources[1] += desireRate[1]*1000
            aBullet.active = false
        }
        if pet.doesCollideWith(&eBullet.ent) && eBullet.active {       // ball get
            fmt.Println("hit E")
            pet.resources[2] += desireRate[2]*1000
            eBullet.active = false
            petBall = false
        }


        // dynamic desires
        pet.resources[0] -= desireRate[0]
        if pet.resources[0] <= 0.0 {
            pet.resources[0] = 0.0
        }
        if pet.resources[0] >= 1.0 {
            pet.resources[0] = 1.0
        }
        pet.resources[1] -= desireRate[1]
        if pet.resources[1] <= 0.0 {
            pet.resources[1] = 0.0
        }
        if pet.resources[1] >= 1.0 {
            pet.resources[1] = 1.0
        }
        pet.resources[2] -= desireRate[2]
        if pet.resources[2] <= 0.0 {
            pet.resources[2] = 0.0
        }
        if pet.resources[2] >= 1.0 {
            pet.resources[2] = 1.0
        }
        emoSum := (pet.resources[0] + pet.resources[1] + pet.resources[2]) / 3.0

        if emoSum <= 0.0 {
            emotion = 4    // DEAD
        } else if emoSum < 0.3 {
            emotion = 3
        } else if emoSum < 0.5 {
            emotion = 2
        } else if emoSum < 0.7 {
            emotion = 1
        } else if emoSum < 0.9 {
            emotion = 0
        }

        // emotional state
        switch emotion {
        case 0:  // blissful
            // lower desire rates
        case 1: // happy
            // increase playful rate
        case 2: // concerned
            // increase desire rates
        case 3: // woeful
            // desire playful rate
        }


        return nil
}

func (s *LevelScene) updateBullets(state *GameState) error {
        // Hunger
        if hBullet.active {
            if hBullet.ent.isOnScreen() {   // not out of bounds yet
                hBullet.ent.moveByVecComponents(hBullet.vec[0], hBullet.vec[1])
            } else {                                                    // at target; clear
                hBullet.active = false
                hBullet.ent.x, hBullet.ent.y = owner.centerPos()
            }
        } else {  // keep bullet on owner
            hBullet.ent.x, hBullet.ent.y = owner.centerPos()
        }
        // Attention
        if aBullet.active {
            // keep bullet on owner
            aBullet.ent.x, aBullet.ent.y = owner.pos()
            aBullet.count -= 1
            if aBullet.count <= 0 {   // despawned; clear
                aBullet.active = false
                aBullet.ent.x, aBullet.ent.y = owner.pos()
            }
        }
        // Exercise
        if eBullet.active {
            if !eBullet.ent.isAtPoint(eBullet.target[0], eBullet.target[1]) {   // not at target yet
                eBullet.ent.moveTowardPoint(eBullet.target[0], eBullet.target[1])
            } else {                                                    // stay until pet collision / despawn
                // check if landing space is passable
                dstX, dstY := eBullet.ent.cellPos()
                if collisionMap[mapGraph.Coord2Indx(dstX, dstY)] { // if collision, change dst to a neighboring tile
                    nl := mapGraph.getNeighbors(mapGraph.Coord2Indx(dstX, dstY))
                    for _, v := range nl { // use first available neighbor
                        if v != -1 {
                            dstX, dstY = mapGraph.Indx2Coord(v)
                            // translate from cell grid to pixels
                            x, y := float64((dstX * tileSize) + tileSize/2), float64((dstY * tileSize) + tileSize/2)
                            eBullet.target[0], eBullet.target[1] = x, y   // move to neighboring tile
                        }
                    }
                }
                if !petBall {   // if active & at target & not started petBall yet
                    ballLanded = true  // signal landing
                    //fmt.Printf("Ball has landed\n")
                }

                eBullet.count -= 1    // at target so countdown to despawn
                if eBullet.count <= 0 {   // despawned; clear
                    eBullet.active = false
                    petBall = false
                    eBullet.ent.x, eBullet.ent.y = owner.centerPos()
                }
            }
        } else {  // keep bullet on owner
            eBullet.ent.x, eBullet.ent.y = owner.centerPos()
        }

    return nil
}


/* Draw */

func (s *LevelScene) Draw(r *ebiten.Image) {
        /* Debug */
        //ebitenutil.DebugPrint(r, "\nNothing here yet :(")
	message := ""
	x := 0
	y := ScreenHeight - 48
	//drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)

        /* Draw Map */
        s.drawMap(r)

        /* Draw Characters */
        s.drawChars(r)

        /* Draw Bullets */
        s.drawBullets(r)

        /* Display current weapon */
        switch activeWeapon {
        case 0:
            message = fmt.Sprintf("Food Blaster")
        case 1:
            message = fmt.Sprintf("Love Bomb")
        case 2:
            message = fmt.Sprintf("Ball Launcher")
        }

        s.drawUI(r)

	drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)

        //ebitenutil.DebugPrint(r, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
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
        // Draw Owner
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(owner.pos())
        r.DrawImage(owner.image, op)

        // Draw Pet
        op = &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(pet.x), float64(pet.y))

        // change pet image
        switch emotion {
        case 0:  // blissful
            r.DrawImage(blissfulImage, op)
        case 1: // happy
            r.DrawImage(happyImage, op)
        case 2: // concerned
            r.DrawImage(concernedImage, op)
        case 3: // woeful
            r.DrawImage(woefulImage, op)
        }


}

func (s *LevelScene) drawUI(r *ebiten.Image) {
        // Draw Pet Desire Alerts
        x, y := pet.posInt()
        if pet.resources[0] <= 0.5 {
	    message := fmt.Sprintf("I'm Hungry!")
	    drawTextWithShadow(r, message, x-20, y-10, 1, color.NRGBA{0x80, 0, 0, 0xff})
        }
        if pet.resources[1] <= 0.5 {
	    message := fmt.Sprintf("Love Me!")
	    drawTextWithShadow(r, message, x-20, y, 1, color.NRGBA{0x80, 0, 0, 0xff})
        }
        if pet.resources[2] <= 0.5 {
	    message := fmt.Sprintf("Play Ball!")
	    drawTextWithShadow(r, message, x-20, y+10, 1, color.NRGBA{0x80, 0, 0, 0xff})
        }

        uiSqr.Fill(color.Black)
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(0, 0)
        r.DrawImage(uiSqr, op)
        x, y = 5, 5
	message := fmt.Sprintf("H: %0.2f", pet.resources[0])
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})
	message = fmt.Sprintf("A: %0.2f", pet.resources[1])
	drawTextWithShadow(r, message, x, y+10, 1, color.NRGBA{0x80, 0, 0, 0xff})
	message = fmt.Sprintf("E: %0.2f", pet.resources[2])
	drawTextWithShadow(r, message, x, y+20, 1, color.NRGBA{0x80, 0, 0, 0xff})
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
        aBullet.count = 50     //set despawn counter
        aBullet.active = true
    case 2:
        eBullet.target[0] = float64(x)
        eBullet.target[1] = float64(y)
        eBullet.count = 100     //set despawn counter
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
