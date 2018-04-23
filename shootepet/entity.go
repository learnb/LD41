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
        "math"
	"github.com/hajimehoshi/ebiten"
)


type Entity struct {
    image *ebiten.Image
    w, h int
    x, y float64
    speed float64
    resoures [3]float64
}

func (e *Entity) setSizeByImage() {
    rect := e.image.Bounds()
    e.w = rect.Dx()
    e.h = rect.Dy()
}

func (e *Entity) size() (int, int) {
    return e.w, e.h
}

func (e *Entity) posInt() (int, int) {
    return int(e.x), int(e.y)
}

func (e *Entity) pos() (float64, float64) {
    return e.x, e.y
}

func (e *Entity) centerPos() (float64, float64) {
    return e.x+float64(e.w)/2, e.y+float64(e.h)/2
}

func (e *Entity) centerPosInt() (int, int) {
    x := int(e.x)
    y := int(e.y)
    w := int(e.w)
    h := int(e.h)
    return x+w/2, y+h/2
}

func (e *Entity) moveTowardCell(cX, cY int) {
    x := float64(cX * tileSize)
    y := float64(cY * tileSize)
    e.moveTowardPoint(x,y)
}

func (e *Entity) moveTowardPoint(cX, cY float64) {
    dx := cX - e.x
    dy := cY - e.y

    d := math.Sqrt( math.Pow((dx),2) + math.Pow((dy),2) )
    if d < 0.2 {
        d = 0.2
    }
    normedX := dx / d
    normedY := dy / d

    e.x += normedX * e.speed //time delta?
    e.y += normedY * e.speed
}

func (e *Entity) getVecComponents(tX, tY float64) (float64, float64) {
    dx := tX - e.x
    dy := tY - e.y

    d := math.Sqrt( math.Pow((dx),2) + math.Pow((dy),2) )
    if d < 0.02 {
        d = 0.02
    }
    normedX := dx / d
    normedY := dy / d

    return normedX * e.speed, normedY * e.speed
}

func (e *Entity) moveByVecComponents(cX, cY float64) {
    e.x += cX
    e.y += cY
}

func (e *Entity) isAtPoint(cX, cY float64) bool {
    d := dist(e.x, e.y, cX, cY)
    fmt.Printf("d: %0.2f\n", d)
    if d < 10.0 {
        return true
    }
    return false
}

func (e *Entity) isAtCell(cX, cY int) bool {
    d := dist(e.x, e.y, float64(cX * tileSize), float64(cY * tileSize))
    if d < 2.0 {
        return true
    }
    return false
}

func (e *Entity) isOnScreen() bool {
    if e.x >= 0.0 && e.x < ScreenWidth {
        if e.y >= 0.0 && e.y < ScreenHeight{
            return true
        }
    }
    return false
}

func (e *Entity) doesCollideWith(a *Entity) bool {
    x1, y1 := e.centerPos()
    x2, y2 := a.centerPos()
    if dist(x1, y1, x2, y2) <= float64(e.w){
        return true
    }
    return false
}

func dist(x1, y1, x2, y2 float64) float64 {
    return math.Sqrt( math.Pow((x1-x2),2) + math.Pow((y1-y2),2) )
}

func (e *Entity) distanceTo(a *Entity) float64 {
    x1, y1 := e.centerPos()
    x2, y2 := a.centerPos()
    return dist(x1, y1, x2, y2)
}
