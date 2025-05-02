package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

func main() {
	// Initialize the termbox library for terminal graphics
	err := termbox.Init()
	if err != nil {
		log.Fatal(err) // Log and exit if initialization fails
	}
	defer termbox.Close() // Ensure termbox is closed when the program exits

	// Creates a new game instance
	game := NewGame()

	// Starts the game loop
	game.Run()
}
