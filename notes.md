# Virtual Pet Shooter

## Genre Descriptions

### Shooter
* Mechanics:
  * Perform actions at range
  * Control avatar movement
  * Manage resources (ammo)

* Goals:
  * Avoid death

### Virtual Pet
* Mechanics:
  * Interact with pet
  * Manage pet's happiness through mini-games

* Goals:
  * Keep pet alive
  * Keep pet happy as possible

### Virtual Pet Shooter
* Mechanics:
  * Control avatar movement
  * Interact with pet
  * Manage ammo
  * Manage pet's happiness through shooter actions

* Goals:
  * Keep pet alive
  * Keep pet happy as possible

## Game Description

The screen features an expressive Pet. The player controls the Trainer on the 
screen. The Trainer can interact with the Pet by launching objects over to it.

## Goal

Keep the Pet alive and as happy as possible.

## Entities

* Pet
  * Roams around the screen
  * Has several needs:
    * Hunger
    * Attention
    * Exercise

* Trainer
  * Collects Ammo of several types:
    * Hunger
    * Attention
    * Exercise
  * Manages Pet's needs by shooting it

* UI
  * Icon to show Pet's overall happiness
  * Show Trainer's Ammo levels
  * Timer (Pet lifetime)
  * Save/Load button

* Level
  * Large open screen in a top-down view.
  * Maze-like layout: Obstacles in the way of Pet.

## Actions

* Movement
  * Collect pickups (Ammo)
* Shoot
  * Hunger
    * Bullet Style: Line of sight
    * Target: Pet
  * Attention
    * Bullet Style: Area of effect
    * Target: Pet
  * Exercise
    * Bullet Style: Mortar
    * Target: Away from Pet


## Jam Plan

Day 1: Process theme and come up with ideas. Start planning phase.

Day 2: Complete planning phase (Morning). Implement framework of game without assests (Night).

Day 3: Build test level into final level: Gameplay test & tweak; Add in assests.
