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
        "github.com/learnb/ld41"
        "github.com/hajimehoshi/ebiten"
)

const (
        ScreenWidth  = 480
        ScreenHeight = 320
)

type Game struct {
        sceneManager *SceneManager
        input        Input
}

func (g *Game) Update(r *ebiten.Image) error {
        if g.sceneManager == nil {
                g.sceneManager = &SceneManager{}
                g.sceneManager.GoTo(&TitleScene{})
        }

        g.input.Update()
        if err := g.sceneManager.Update(&g.input); err != nil {
                return err
        }
        if ebiten.IsRunningSlowly() {
                return nil
        }

        g.sceneManager.Draw(r)
        return nil
}
