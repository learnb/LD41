package main

import (
    "image/color"

    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
    myR uint8
    myG uint8
    myB uint8
    myA uint8
)

func init() {
    myR = 0xff
    myG = 0x00
    myB = 0x00
    myA = 0xff
}

func update(screen *ebiten.Image) error {
    // Color!
    myR += 0x01
    myG -= 0x01
    myB += 0x04
    //myA += 0x01

    screen.Fill(color.NRGBA{myR, myG, myB, myA})

    // Text!
    ebitenutil.DebugPrint(screen, "Fuck Donald Trump!")

    return nil
}

func main() {
    if err := ebiten.Run(update, 480, 320, 2, "Whoahoho!"); err != nil {
        panic(err)
    }
}
