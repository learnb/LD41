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

    dude *ebiten.Image
)

func init() {
    myR = 0xff
    myG = 0x00
    myB = 0x00
    myA = 0xff

    dude, _ = ebiten.NewImage(16,16, ebiten.FilterNearest)
}

func update(screen *ebiten.Image) error {

    // Color!
    myR += 0x01
    myG -= 0x01
    myB += 0x04
    screen.Fill(color.NRGBA{myR, myG, myB, myA})

    // Dude
    dude.Fill(color.Black)

    // Text!
    ebitenutil.DebugPrint(screen, "Shoot to Pet")

    // Set image draw options
    opts := &ebiten.DrawImageOptions{}
    opts.GeoM.Translate(64,64)

    // Draw image to screen
    screen.DrawImage(dude, opts)

    /* User Input */
    if (ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton10) || ebiten.IsKeyPressed(ebiten.KeyUp)){
        ebitenutil.DebugPrint(screen, "\nUP")
    }
    if (ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton12) || ebiten.IsKeyPressed(ebiten.KeyDown)){
        ebitenutil.DebugPrint(screen, "\n\nDOWN")
    }
    if (ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton13) || ebiten.IsKeyPressed(ebiten.KeyLeft)){
        ebitenutil.DebugPrint(screen, "\n\n\nLEFT")
    }
    if (ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton11) || ebiten.IsKeyPressed(ebiten.KeyRight)){
        ebitenutil.DebugPrint(screen, "\n\n\n\nRIGHT")
    }
    /* END - User Input */


    return nil
}

func main() {
    if err := ebiten.Run(update, 480, 320, 2, "Ludum Dare 41 | Shoot to Pet"); err != nil {
        panic(err)
    }
}
