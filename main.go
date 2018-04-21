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

package main

import (
        "flag"
        "log"
        "os"
        "runtime/pprof"

        "github.com/hajimehoshi/ebiten"
        "github.com/learnb/ld41/shootepet"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
        flag.Parse()
        if *cpuProfile != "" {
                f, err := os.Create(*cpuProfile)
                if err != nil {
                        log.Fatal(err)
                }
                if err := pprof.StartCPUProfile(f); err != nil {
                        log.Fatal(err)
                }
                defer pprof.StopCPUProfile()
        }

        game := &shootepet.Game{}
        update := game.Update
        if err := ebiten.Run(update, shootepet.ScreenWidth, shootepet.ScreenHeight, 2, "Shoot'E'Pet | Ludum Dare 41"); err != nil {
                log.Fatal(err)
        }
}