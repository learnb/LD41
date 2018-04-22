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
        "math"
	"github.com/hajimehoshi/ebiten"
)


type Entity struct {
    image *ebiten.Image
    x, y int
    w, h int
    resoures [3]float32
}

func (e *Entity) setSizeByImage() {
    rect := e.image.Bounds()
    e.w = rect.Dx()
    e.h = rect.Dy()
}

func (e *Entity) size() (int, int) {
    return e.w, e.h
}

func (e *Entity) pos() (int, int) {
    return e.x, e.y
}

func (e *Entity) posf() (float64, float64) {
    return float64(e.x), float64(e.y)
}

func (e *Entity) centerPos() (int, int) {
    return e.x+e.w/2, e.y+e.h/2
}

func (e *Entity) centerPosf() (float64, float64) {
    x := float64(e.x)
    y := float64(e.y)
    w := float64(e.w)
    h := float64(e.h)
    return x+w/2, y+h/2
}

func (e *Entity) moveToCell(cX, cY int) {
    e.x = cX * tileSize
    e.y = cY * tileSize

}

func (e *Entity) doesCollideWith(a *Entity) bool {
    x1, y1 := e.centerPosf()
    x2, y2 := a.centerPosf()
    if dist(x1, y1, x2, y2) <= float64(e.w){
        return true
    }
    return false
}

func dist(x1, y1, x2, y2 float64) float64 {
    return math.Sqrt( math.Pow((x1-x2),2) + math.Pow((y1-y2),2) )
}

func (e *Entity) distanceTo(a *Entity) float64 {
    x1, y1 := e.centerPosf()
    x2, y2 := a.centerPosf()
    return dist(x1, y1, x2, y2)
}
