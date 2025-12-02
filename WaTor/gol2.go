/**
 * @file gol2.go
 * @brief Wa-Tor my solution
 * @author Temur Rustamov
 * @date 2025-11-28
 * @license GNU General Public License v3.0
 */

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

/**
 * @brief pixel scale
 */
const scale = 2

/**
 * @brief pixel size
 */
const width = 300

/**
 * @brief grid cell
 */
const height = 300

/**
 * @brief grid cell
 */
const fishBreedTime = 3
const sharkBreedTime = 8
const sharkStarve = 3

/**
 * @brief 0 to 100 percent
 */
const sharkPercent = 10
const fishPercent = 90

type Cell struct {
	Type       uint8
	BreedTime  int
	StarveTime int
}

const (
	Empty uint8 = iota
	Fish
	Shark
)

var grid [width][height]Cell
var count int
var buffer [width][height]Cell

func directions() [][2]int {
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := len(directions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		directions[i], directions[j] = directions[j], directions[i]
	}
	return directions
}

func initWorld() {
	rand.Seed(1)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := rand.Intn(100)
			c := Cell{}
			if r < fishPercent {
				c.Type = Fish
				c.BreedTime = fishBreedTime
			} else if r < fishPercent+sharkPercent {
				c.Type = Shark
				c.BreedTime = sharkBreedTime
				c.StarveTime = sharkStarve
			} else {
				c.Type = Empty
			}
			grid[x][y] = c
		}
	}
}

func wrap(x, max int) int {
	if x < 0 {
		return max - 1
	}
	if x >= max {
		return 0
	}
	return x
}

func update() {
	for x := 0; x < width; x++ {
		/**
		 * @brief clear buffer
		 */
		for y := 0; y < height; y++ {
			buffer[x][y] = Cell{}
		}
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			cell := grid[x][y]
			if cell.Type != Shark {
				continue
			}
			directions := directions()
			fishTargets := make([][2]int, 0, 4)
			emptyTargets := make([][2]int, 0, 4)
			for _, d := range directions {
				/**
				 * @brief look at neighbours in random order
				 */
				nx := wrap(x+d[0], width)
				ny := wrap(y+d[1], height)
				n := grid[nx][ny]
				if n.Type == Fish {
					fishTargets = append(fishTargets, [2]int{nx, ny})
				} else if n.Type == Empty {
					emptyTargets = append(emptyTargets, [2]int{nx, ny})
				}
			}
			moved := false
			c := cell
			/**
			 * @brief prefer only fish
			 */
			if len(fishTargets) > 0 {
				target := fishTargets[rand.Intn(len(fishTargets))]
				tx, ty := target[0], target[1]
				if buffer[tx][ty].Type == Empty {
					c.BreedTime--
					c.StarveTime = sharkStarve // reset
					/**
					 * @brief reproduction
					 */
					if c.BreedTime <= 0 {
						/**
						 * @brief new shark leave behind
						 */
						if buffer[x][y].Type == Empty {
							buffer[x][y] = Cell{
								Type:       Shark,
								BreedTime:  sharkBreedTime,
								StarveTime: sharkStarve,
							}
						}
						c.BreedTime = sharkBreedTime
					}

					buffer[tx][ty] = c
					moved = true
				}
			}
			/**
			 * @brief no fish is eaten ok so try to move to different cell
			 */
			if !moved && len(emptyTargets) > 0 {
				target := emptyTargets[rand.Intn(len(emptyTargets))]
				tx, ty := target[0], target[1]
				if buffer[tx][ty].Type == Empty {
					c.BreedTime--
					c.StarveTime--
					/**
					 * @brief dead of starvation no food etc
					 */
					if c.StarveTime <= 0 {
						/**
						 * @brief sharky disappears ciao cacao
						 */
						moved = true
					} else {
						if c.BreedTime <= 0 {
							if buffer[x][y].Type == Empty {
								buffer[x][y] = Cell{
									Type:       Shark,
									BreedTime:  sharkBreedTime,
									StarveTime: sharkStarve,
								}
							}
							c.BreedTime = sharkBreedTime
						}
						buffer[tx][ty] = c
						moved = true
					}
				}
			}
			/**
			 * @brief stay in one place
			 */
			if !moved {
				c.BreedTime--
				c.StarveTime--
				if c.StarveTime > 0 {
					if c.BreedTime <= 0 {
						c.BreedTime = sharkBreedTime
					}
					if buffer[x][y].Type == Empty {
						buffer[x][y] = c
					}
				}
				/**
				 * @brief here could be dead
				 */
			}
		}
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			cell := grid[x][y]
			if cell.Type != Fish {
				continue
			}
			/**
			 * @brief if shark or something already placed here in my buffer, fish is eaten
			 */
			if buffer[x][y].Type != Empty && buffer[x][y].Type != Fish {
				continue
			}

			directions := directions()
			emptyTargets := make([][2]int, 0, 4)
			for _, d := range directions {
				nx := wrap(x+d[0], width)
				ny := wrap(y+d[1], height)
				/**
				 * @brief can move only into places that were empty before and are still empty in my buffer
				 */
				if grid[nx][ny].Type == Empty && buffer[nx][ny].Type == Empty {
					emptyTargets = append(emptyTargets, [2]int{nx, ny})
				}
			}
			c := cell
			moved := false
			if len(emptyTargets) > 0 {
				target := emptyTargets[rand.Intn(len(emptyTargets))]
				tx, ty := target[0], target[1]
				c.BreedTime--
				if c.BreedTime <= 0 {
					/**
					 * @brief reproduction
					 */
					if buffer[x][y].Type == Empty {
						buffer[x][y] = Cell{
							Type:      Fish,
							BreedTime: fishBreedTime,
						}
					}
					c.BreedTime = fishBreedTime
				} else {
					/**
					 * @brief current position is untouched
					 */
					fmt.Println("nvm")
				}

				if buffer[tx][ty].Type == Empty {
					buffer[tx][ty] = c
				}
				moved = true
			}

			if !moved {
				/**
				 * @brief if nth else put here stay
				 */
				if buffer[x][y].Type == Empty {
					buffer[x][y] = c
				}
			}
		}
	}
	temp := grid
	grid = buffer
	buffer = temp
}

/**
 * @brief draw pixels
 */

func display(window *ebiten.Image) {
	window.Fill(color.RGBA{0, 0, 255, 1})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := grid[x][y]
			var col color.Color
			switch c.Type {
			case Fish:
				col = color.RGBA{245, 40, 145, 1}
			case Shark:
				col = color.RGBA{0, 0, 0, 0}
			default:
				continue
			}
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					window.Set(x*scale+i, y*scale+j, col)
				}
			}
		}
	}
}

/**
 * @brief working on window
 */
func frame(window *ebiten.Image) error {
	count++
	var err error = nil
	if count%2 == 0 {
		update()
	}
	if !ebiten.IsDrawingSkipped() {
		display(window)
	}
	return err
}

func main() {
	initWorld()
	if err := ebiten.Run(frame, width, height, scale, "The best game that humanity has ever seen in their life"); err != nil {
		log.Fatal(err)
	}

}
