// Copyright 2014 Hajime Hoshi
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
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	rshootepet "github.com/learnb/ld41/resources/images/shootepet"
)

var tutorialImageBackground *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(rshootepet.Title_png))
	if err != nil {
		panic(err)
	}
	tutorialImageBackground, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func (s *TutorialScene) Init() {
    //init()
}

type TutorialScene struct {
	count int
}

func (s *TutorialScene) Update(state *GameState) error {
	s.count++
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoTo(NewGameScene())
		return nil
	}
	if anyGamepadAbstractButtonPressed(state.Input) {
		state.SceneManager.GoTo(NewGameScene())
		return nil
	}

	// If 'abstract' gamepad buttons are not set and any gamepad buttons are pressed,
	// go to the gamepad configuration scene.
//	if state.Input.IsAnyGamepadButtonPressed() {
//		state.SceneManager.GoTo(&GamepadScene{})
//		return nil
//	}
	return nil
}

func (s *TutorialScene) Draw(r *ebiten.Image) {
	s.drawBackground(r, s.count)

	message := "how to play"
	x := 0
	y := 10
	drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)

        // Movement
	message = "Movement"
        x = 25
	y = 20
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "W"
        x = 50
	y = 45
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "A"
        x = 40
	y = 55
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "S"
        x = 50
	y = 55
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "D"
        x = 60
	y = 55
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})


        // Shooting
	message = "Shooting"
        x = 25
	y = 100
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Aim:   mouse cursor"
        x = 25
	y = 120
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Shoot: left mouse button"
	x = 25
	y = 140
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Change Weapons: E"
	x = 25
	y = 160
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "    - Food Blaster, Love Bomb, Ball Cannon"
	x = 25
	y = 175
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

        // Goal
	message = "Goal"
        x = 25
	y = 200
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Keep your pet happy by shooting it in the face!"
        x = 25
	y = 220
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Hungry? Shoot 'em with the Food Blaster"
	x = 25
	y = 240
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Lonely? Get close and use your Love Bomb"
	x = 25
	y = 260
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})

	message = "Playful? Use the Ball Cannon to play fetch"
	x = 25
	y = 280
	drawTextWithShadow(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff})


}

func (s *TutorialScene) drawBackground(r *ebiten.Image, c int) {
	w, h := tutorialImageBackground.Size()
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < (ScreenWidth/w+1)*(ScreenHeight/h+2); i++ {
		op.GeoM.Reset()
		dx := -(c / 4) % w
		dy := (c / 4) % h
		dstX := (i%(ScreenWidth/w+1))*w + dx
		dstY := (i/(ScreenWidth/w+1)-1)*h + dy
		op.GeoM.Translate(float64(dstX), float64(dstY))
		r.DrawImage(tutorialImageBackground, op)
	}
}

