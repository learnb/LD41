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
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	rshootepet "github.com/learnb/ld41/resources/images/shootepet"
)

var gameoverImageBackground *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(rshootepet.Title_png))
	if err != nil {
		panic(err)
	}
	gameoverImageBackground, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type GameOverScene struct {
	count int
}

func (s *GameOverScene) Update(state *GameState) error {
        s.count++
        if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
                state.SceneManager.GoTo(&TitleScene{})
                return nil
        }
        if anyGamepadAbstractButtonPressed(state.Input) {
                state.SceneManager.GoTo(&TitleScene{})
                return nil
        }
	return nil
}

func (s *GameOverScene) Draw(r *ebiten.Image) {
        ebitenutil.DebugPrint(r, "\nYou let your adorable pet die!\n\nPress FIRE to restart")
	message := "~ Game Over ~"
	x := 0
	y := ScreenHeight - 48
	drawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, ScreenWidth)
}

