# Wa-Tor Simulation 

Author: **Temur Rustamov**  
License: **GNU General Public License v3.0**

## Overview

This project implements the classic **Wa-Tor** predator–prey simulation on a toroidal grid, inspired by:

- A.K. Dewdney, *“Computer Recreations; Sharks and Fish wage an ecological war on the toroidal planet of Wa-Tor”*, Scientific American.
- Online descriptions such as the Wa-Tor article on Wikipedia.

The world is a 2D grid wrapped like a torus. Each cell may contain:

- a **fish**
- a **shark**
- or be **empty**

Time advances in discrete steps called **chronons**. At each chronon, fish and sharks move, reproduce, and (for sharks) starve or eat according to the rules described below.

The simulation is written in **Go**, uses **Ebiten** for simple graphical output, and is documented with **Doxygen**.

---

## Implementation

- Language: **Go**
- OS target: **Linux** (but also runs on macOS/Windows if Go + Ebiten installed)
- Graphics: **Ebiten** (pixel grid showing fish and sharks)
- Documentation: **Doxygen** (this README + Doxyfile + mainpage.dox)

Current code is implemented in a single file:

- `gol2.go` – main Wa-Tor simulation and Ebiten loop.

Many parameters (breed time, starve time, grid size) are currently **hard-coded** as Go constants. They could be extended to be read from command-line arguments to exactly match the assignment specification.

---

## Wa-Tor Rules (Assignment Spec)

### Parameters

The simulation conceptually takes seven parameters:

- **NumShark**: Starting population of sharks
- **NumFish**: Starting population of fish
- **FishBreed**: Number of chronons before a fish can reproduce
- **SharkBreed**: Number of chronons before a shark can reproduce
- **Starve**: Number of chronons a shark can survive without eating before it dies
- **GridSize**: Dimensions of the world (e.g. 300×300)
- **Threads**: Number of threads to use (for a parallel implementation and speedup experiments)

In this implementation, typical values are:

- `fishBreedTime = 3`
- `sharkBreedTime = 8`
- `sharkStarve = 3`
- `width = 300`
- `height = 300`
- Initial fish/shark distribution is random, based on `fishPercent` and `sharkPercent`.

### Movement and Chronons

Time passes in discrete **chronons**. In each chronon:

- A fish or shark may move **north, east, south or west** to an adjacent cell.
- Movement is controlled by a **random-number generator**.
- They **cannot** move into a cell already occupied by the same species.

### Fish Rules

- At each chronon, a fish:
    - Chooses one of the adjacent **unoccupied** cells at random and moves there.
    - If all four neighbors are occupied, it does **not** move.
- Reproduction:
    - After surviving `FishBreed` chronons, a fish can reproduce.
    - When it moves to a new cell, it **leaves a new fish behind** in its old cell.
    - Its reproduction counter is reset to zero.

### Shark Rules

- Priority: **hunting fish** has priority over simple movement.
- At each chronon, a shark:
    - Looks at adjacent cells:
        - If one or more contain **fish**, it chooses one at random, moves there, and **eats the fish**.
        - If no fish are adjacent, it moves like a fish: to a random adjacent **empty** cell, if possible.
        - If all neighbors are sharks or blocked, the shark does **not** move.
- Energy / starvation:
    - Each chronon, a shark loses one unit of **energy**.
    - If its energy reaches zero, the shark **dies**.
    - When a shark eats a fish, it **gains energy** (reset up to `Starve` in this implementation).
- Reproduction:
    - After surviving `SharkBreed` chronons, a shark can reproduce.
    - Reproduction happens in the **same way as fish**: when it moves, it leaves a newborn shark in its previous cell and resets its reproduction counter.

---

## Building and Running

### Requirements

- Go (e.g. Go 1.20+)
- Ebiten library:
  ```bash
  go get github.com/hajimehoshi/ebiten
  
use the command
```bash
sudo apt install doxygen graphviz

