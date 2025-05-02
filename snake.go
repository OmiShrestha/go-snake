package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Replace global random generator with a local one
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Represents a point in 2D space
// Used for the snake's body, food, obstacles, and portals
type Point struct {
	X, Y int
}

// Entry and exit points for a portal
type Portal struct {
	Entry Point
	Exit  Point
}

// Game state
type Game struct {
	snake     []Point // The snake's body
	direction Point   // Current direction of the snake
	food      Point   // Position of the food
	width     int     // Width of the game board
	height    int     // Height of the game board
	gameOver  bool    // Indicates if the game is over
	score     int     // Current score
	paused    bool    // Indicates if the game is paused
	obstacles []Point // Positions of obstacles
	portal    Portal  // Portal entry and exit points
}

// Initializes a new game instance
func NewGame() *Game {
	w, h := termbox.Size()                  // Get terminal size
	s := []Point{{w / 2, h / 2}}            // Initialize the snake at the center
	obstacles := generateObstacles(w, h, s) // Generate obstacles

	portal := Portal{
		Entry: randomPoint(w, h, append(s, obstacles...)),
		Exit:  randomPoint(w, h, append(s, obstacles...)),
	}

	return &Game{
		snake:     s,
		direction: Point{1, 0}, // Initial direction: right
		food:      randomFood(w, h, append(s, obstacles...)),
		width:     w,
		height:    h,
		obstacles: obstacles,
		portal:    portal,
	}
}

// Generates random obstacles on the board
func generateObstacles(w, h int, snake []Point) []Point {
	obstacles := []Point{}
	for i := 0; i < 10; i++ { // Generate 10 obstacles
		for {
			p := Point{rng.Intn(w), rng.Intn(h)}
			conflict := false
			// Ensure obstacles do not overlap with the snake or other obstacles
			for _, s := range snake {
				if s == p {
					conflict = true
					break
				}
			}
			for _, o := range obstacles {
				if o == p {
					conflict = true
					break
				}
			}
			if !conflict {
				obstacles = append(obstacles, p)
				break
			}
		}
	}
	return obstacles
}

// Main game loop
func (g *Game) Run() {
	ticker := time.NewTicker(120 * time.Millisecond) // Game tick interval
	defer ticker.Stop()

	go g.pollEvents() // Start listening for user input

	for !g.gameOver {
		<-ticker.C // Wait for the next tick
		g.update() // Update game state
		g.draw()   // Render the game
	}
	g.drawGameOver()            // Display game over screen
	time.Sleep(2 * time.Second) // Pause before exiting
}

// Updates the game state
func (g *Game) update() {
	if g.paused {
		return // Do nothing if the game is paused
	}

	head := g.snake[0]                                               // Get the current head of the snake
	newHead := Point{head.X + g.direction.X, head.Y + g.direction.Y} // Calculate new head position

	// Check for collisions
	if newHead.X < 0 || newHead.Y < 0 || newHead.X >= g.width || newHead.Y >= g.height {
		g.gameOver = true // Collision with walls
		return
	}
	for _, p := range g.snake {
		if p == newHead {
			g.gameOver = true // Collision with itself
			return
		}
	}
	for _, o := range g.obstacles {
		if o == newHead {
			g.gameOver = true // Collision with obstacles
			return
		}
	}

	// Handle portal logic (make entry and exit interchangeable)
	if newHead == g.portal.Entry {
		newHead = g.portal.Exit
	} else if newHead == g.portal.Exit {
		newHead = g.portal.Entry
	}

	// Move the snake
	g.snake = append([]Point{newHead}, g.snake...)
	if newHead == g.food {
		g.food = randomFood(g.width, g.height, append(g.snake, g.obstacles...)) // Generate new food
		g.score++                                                               // Increase score
	} else {
		g.snake = g.snake[:len(g.snake)-1] // Remove the tail
	}
}

// Draws the game border
func (g *Game) drawBorder() {
	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '─', termbox.ColorWhite, termbox.ColorDefault)          // Top border
		termbox.SetCell(x, g.height-1, '─', termbox.ColorWhite, termbox.ColorDefault) // Bottom border
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '│', termbox.ColorWhite, termbox.ColorDefault)         // Left border
		termbox.SetCell(g.width-1, y, '│', termbox.ColorWhite, termbox.ColorDefault) // Right border
	}
	// Draw corners
	termbox.SetCell(0, 0, '┌', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, 0, '┐', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.height-1, '└', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, g.height-1, '┘', termbox.ColorWhite, termbox.ColorDefault)
}

// Renders the game state
func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault) // Clear the screen
	g.drawBorder()                                            // Draw the border

	// Draw the snake
	for i, p := range g.snake {
		if i == 0 {
			termbox.SetCell(p.X, p.Y, '@', termbox.ColorGreen, termbox.ColorDefault) // Head
		} else {
			termbox.SetCell(p.X, p.Y, 'o', termbox.ColorGreen, termbox.ColorDefault) // Body
		}
	}

	// Draw the food
	termbox.SetCell(g.food.X, g.food.Y, '*', termbox.ColorRed, termbox.ColorDefault)

	// Draw obstacles
	for _, o := range g.obstacles {
		termbox.SetCell(o.X, o.Y, '?', termbox.ColorMagenta, termbox.ColorDefault)
	}

	// Draw portal
	termbox.SetCell(g.portal.Entry.X, g.portal.Entry.Y, 'O', termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(g.portal.Exit.X, g.portal.Exit.Y, 'O', termbox.ColorBlue, termbox.ColorDefault)

	// Display score
	scoreStr := fmt.Sprintf("Score: %d", g.score)
	for i, c := range scoreStr {
		termbox.SetCell(i+1, 1, c, termbox.ColorYellow, termbox.ColorDefault)
	}

	termbox.Flush() // Render the changes
}

// Displays the game over screen
func (g *Game) drawGameOver() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault) // Clear the screen
	msg := "GAME OVER"
	for i, c := range msg {
		termbox.SetCell(g.width/2-len(msg)/2+i, g.height/2, c, termbox.ColorRed, termbox.ColorDefault) // Center the message
	}
	termbox.Flush() // Render the changes
}

// Handles user input (arrow keys for movement, space to pause, escape to exit)
func (g *Game) pollEvents() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if g.direction.Y != 1 {
					g.direction = Point{0, -1} // Move up
				}
			case termbox.KeyArrowDown:
				if g.direction.Y != -1 {
					g.direction = Point{0, 1} // Move down
				}
			case termbox.KeyArrowLeft:
				if g.direction.X != 1 {
					g.direction = Point{-1, 0} // Move left
				}
			case termbox.KeyArrowRight:
				if g.direction.X != -1 {
					g.direction = Point{1, 0} // Move right
				}
			case termbox.KeyEsc:
				g.gameOver = true // Exit the game
				return
			case termbox.KeySpace:
				g.paused = !g.paused // Toggle pause
			}
		}
	}
}

// Generates a random position for the food
func randomFood(w, h int, snake []Point) Point {
	for {
		p := Point{rng.Intn(w), rng.Intn(h)}
		conflict := false
		// Ensure the food does not overlap with the snake
		for _, s := range snake {
			if s == p {
				conflict = true
				break
			}
		}
		if !conflict {
			return p
		}
	}
}

// Generates a random positions for the portal
func randomPoint(w, h int, occupied []Point) Point {
	for {
		p := Point{rng.Intn(w), rng.Intn(h)}
		conflict := false
		// Ensure the portal does not overlap with occupied points
		for _, o := range occupied {
			if o == p {
				conflict = true
				break
			}
		}
		if !conflict {
			return p
		}
	}
}
