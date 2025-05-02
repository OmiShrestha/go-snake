# Snake Game

This is a terminal-based Snake game implemented in Go using the `termbox-go` library.

## Prerequisites

To run this game, ensure you have the following installed on your system:

1. **Go Programming Language**: Install Go from [golang.org](https://golang.org/dl/).
2. **Git**: Ensure Git is installed to clone the repository and manage dependencies.

## How to Run

Follow these steps to run the Snake game:

1. **Clone the Repository** (if applicable):
   ```bash
   git clone <repository-url>
   cd gosnake
   ```

2. **Install Dependencies**:
   Run the following command to download the required dependencies:
   ```bash
   go mod tidy
   ```

3. **Build the Game**:
   Compile the game by running:
   ```bash
   go build -o gosnake main.go snake.go
   ```

4. **Run the Game**:
   Execute the compiled binary:
   ```bash
   ./gosnake
   ```

## Controls

- **Arrow Keys**: Control the direction of the snake.
- **ESC**: Exit the game.

## How the Game Works

- The snake moves continuously in the current direction.
- Use the arrow keys to change the direction of the snake.
- The game ends if the snake collides with the walls or itself.
- The goal is to eat the food (`*`) to grow the snake.

## Troubleshooting

1. **Go Not Installed**:
   If the `go` command is not recognized, install Go and ensure it is added to your system's PATH.

2. **Missing Dependencies**:
   If you encounter errors about missing dependencies, run:
   ```bash
   go mod tidy
   ```

3. **Terminal Issues**:
   Ensure your terminal supports the `termbox-go` library. If the game does not render correctly, try a different terminal emulator.

## License

This project is open-source and available under the MIT License.