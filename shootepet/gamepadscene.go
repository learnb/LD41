// Copyright 2015 Hajime Hoshi
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
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type GamepadScene struct {
	currentIndex      int
	countAfterSetting int
	buttonStates      []string
}

func (s *GamepadScene) Init() {
        //
}


func (s *GamepadScene) Update(state *GameState) error {
	if s.currentIndex == 0 {
		state.Input.gamepadConfig.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		state.Input.gamepadConfig.Reset()
		state.SceneManager.GoTo(&TitleScene{})
	}

	if s.buttonStates == nil {
		s.buttonStates = make([]string, len(virtualGamepadButtons))
	}
	for i, b := range virtualGamepadButtons {
		if i < s.currentIndex {
			s.buttonStates[i] = strings.ToUpper(state.Input.gamepadConfig.ButtonName(b))
			continue
		}
		if s.currentIndex == i {
			s.buttonStates[i] = "_"
			continue
		}
		s.buttonStates[i] = ""
	}

	if 0 < s.countAfterSetting {
		s.countAfterSetting--
		if s.countAfterSetting <= 0 {
			state.SceneManager.GoTo(&TitleScene{})
		}
		return nil
	}

	b := virtualGamepadButtons[s.currentIndex]
	const gamepadID = 0
	if state.Input.gamepadConfig.Scan(gamepadID, b) {
		s.currentIndex++
		if s.currentIndex == len(virtualGamepadButtons) {
			s.countAfterSetting = ebiten.FPS
		}
	}
	return nil
}

func (s *GamepadScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	if s.buttonStates == nil {
		return
	}

	f := `GAMEPAD CONFIGURATION
(PRESS ESC TO CANCEL)

* Joysticks don't work well for movement right now :(
  Please use the D-pad or non-analog buttons
  Sorry!!


Move Left:     %s

Move Right:    %s

Move Up:       %s

Move Down:     %s

Fire Weapon:   %s

Change Weapon: %s



%s`
	msg := ""
	if s.currentIndex == len(virtualGamepadButtons) {
		msg = "OK!"
	}
	str := fmt.Sprintf(f, s.buttonStates[0], s.buttonStates[1], s.buttonStates[2], s.buttonStates[3], s.buttonStates[4], s.buttonStates[5], msg)
	drawTextWithShadow(screen, str, 16, 16, 1, color.White)
}
