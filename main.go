package main

import (
    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
    ebitenutil.DebugPrint(screen, "Fuck Donald Trump!")
    return nil
}

func main() {
    if err := ebiten.Run(update, 480, 320, 2, "Whoahoho!"); err != nil {
        panic(err)
    }
}
