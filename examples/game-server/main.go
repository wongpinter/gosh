package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"repo.nusatek.id/sugeng/gosh"
)

// GameState represents the current game state
type GameState struct {
	CurrentGame string
	Data        map[string]interface{}
}

// GameHandler implements games over SSH
type GameHandler struct {
	username  string
	gameState *GameState
	score     int
	gamesWon  int
}

// NewGameHandler creates a new game handler
func NewGameHandler(username string) *GameHandler {
	return &GameHandler{
		username: username,
		gameState: &GameState{
			CurrentGame: "menu",
			Data:        make(map[string]interface{}),
		},
		score:    0,
		gamesWon: 0,
	}
}

// Execute implements the CommandHandler interface
func (h *GameHandler) Execute(cmd string) (string, uint32) {
	cmd = strings.TrimSpace(cmd)
	
	if cmd == "" {
		return "", 0
	}
	
	parts := strings.Fields(cmd)
	command := parts[0]
	
	// Global commands available in any game
	switch command {
	case "menu", "main":
		return h.showMainMenu(), 0
	case "help":
		return h.getHelp(), 0
	case "score":
		return h.showScore(), 0
	case "quit", "exit":
		return "Thanks for playing! Goodbye!", 0
	}
	
	// Game-specific commands
	switch h.gameState.CurrentGame {
	case "menu":
		return h.handleMenuCommand(command, parts[1:])
	case "guess":
		return h.handleGuessCommand(command, parts[1:])
	case "rps":
		return h.handleRPSCommand(command, parts[1:])
	case "quiz":
		return h.handleQuizCommand(command, parts[1:])
	case "adventure":
		return h.handleAdventureCommand(command, parts[1:])
	default:
		return "Unknown game state. Type 'menu' to return to main menu.", 1
	}
}

func (h *GameHandler) showMainMenu() string {
	h.gameState.CurrentGame = "menu"
	h.gameState.Data = make(map[string]interface{})
	
	return `🎮 GAME SERVER MAIN MENU 🎮

Available Games:
1. guess    - Number Guessing Game
2. rps      - Rock Paper Scissors
3. quiz     - Trivia Quiz
4. adventure - Text Adventure

Commands:
- <game>    Start a game
- score     Show your score
- help      Show help
- quit      Exit

Choose a game by typing its name!`
}

func (h *GameHandler) handleMenuCommand(command string, args []string) (string, uint32) {
	switch command {
	case "guess", "1":
		return h.startGuessGame(), 0
	case "rps", "2":
		return h.startRPSGame(), 0
	case "quiz", "3":
		return h.startQuizGame(), 0
	case "adventure", "4":
		return h.startAdventureGame(), 0
	default:
		return fmt.Sprintf("Unknown game: %s\nType 'help' to see available games", command), 1
	}
}

// Number Guessing Game
func (h *GameHandler) startGuessGame() string {
	h.gameState.CurrentGame = "guess"
	h.gameState.Data["number"] = rand.Intn(100) + 1
	h.gameState.Data["attempts"] = 0
	h.gameState.Data["maxAttempts"] = 7
	
	return `🔢 NUMBER GUESSING GAME 🔢

I'm thinking of a number between 1 and 100.
You have 7 attempts to guess it!

Commands:
- <number>  Make a guess
- hint      Get a hint
- give up   Surrender
- menu      Return to main menu

Make your first guess!`
}

func (h *GameHandler) handleGuessCommand(command string, args []string) (string, uint32) {
	switch command {
	case "hint":
		return h.getGuessHint(), 0
	case "give", "up", "surrender":
		if len(args) > 0 && args[0] == "up" || command == "surrender" {
			number := h.gameState.Data["number"].(int)
			return fmt.Sprintf("The number was %d. Better luck next time!\nType 'menu' to return to main menu.", number), 0
		}
	}
	
	// Try to parse as number
	guess, err := strconv.Atoi(command)
	if err != nil {
		return "Please enter a number between 1 and 100, or type 'hint' for help.", 1
	}
	
	if guess < 1 || guess > 100 {
		return "Please enter a number between 1 and 100.", 1
	}
	
	return h.processGuess(guess), 0
}

func (h *GameHandler) processGuess(guess int) string {
	number := h.gameState.Data["number"].(int)
	attempts := h.gameState.Data["attempts"].(int) + 1
	maxAttempts := h.gameState.Data["maxAttempts"].(int)
	
	h.gameState.Data["attempts"] = attempts
	
	if guess == number {
		h.score += (maxAttempts - attempts + 1) * 10
		h.gamesWon++
		return fmt.Sprintf("🎉 Congratulations! You guessed it in %d attempts!\n"+
			"You earned %d points!\n"+
			"Type 'menu' to play another game.", 
			attempts, (maxAttempts-attempts+1)*10)
	}
	
	if attempts >= maxAttempts {
		return fmt.Sprintf("💀 Game Over! The number was %d.\n"+
			"Type 'menu' to try again.", number)
	}
	
	var hint string
	if guess < number {
		hint = "📈 Too low!"
	} else {
		hint = "📉 Too high!"
	}
	
	return fmt.Sprintf("%s You have %d attempts left.", hint, maxAttempts-attempts)
}

func (h *GameHandler) getGuessHint() string {
	number := h.gameState.Data["number"].(int)
	var hint string
	
	if number%2 == 0 {
		hint = "The number is even."
	} else {
		hint = "The number is odd."
	}
	
	if number <= 25 {
		hint += " It's in the range 1-25."
	} else if number <= 50 {
		hint += " It's in the range 26-50."
	} else if number <= 75 {
		hint += " It's in the range 51-75."
	} else {
		hint += " It's in the range 76-100."
	}
	
	return "💡 Hint: " + hint
}

// Rock Paper Scissors Game
func (h *GameHandler) startRPSGame() string {
	h.gameState.CurrentGame = "rps"
	h.gameState.Data["wins"] = 0
	h.gameState.Data["losses"] = 0
	h.gameState.Data["ties"] = 0
	
	return `✂️ ROCK PAPER SCISSORS ✂️

Commands:
- rock, r     Play rock
- paper, p    Play paper
- scissors, s Play scissors
- stats       Show game statistics
- menu        Return to main menu

Best of luck! Make your move:`
}

func (h *GameHandler) handleRPSCommand(command string, args []string) (string, uint32) {
	switch command {
	case "stats":
		return h.getRPSStats(), 0
	case "rock", "r":
		return h.playRPS("rock"), 0
	case "paper", "p":
		return h.playRPS("paper"), 0
	case "scissors", "s":
		return h.playRPS("scissors"), 0
	default:
		return "Choose: rock (r), paper (p), or scissors (s)", 1
	}
}

func (h *GameHandler) playRPS(playerMove string) string {
	moves := []string{"rock", "paper", "scissors"}
	computerMove := moves[rand.Intn(3)]
	
	var result string
	var outcome string
	
	if playerMove == computerMove {
		result = "It's a tie!"
		outcome = "tie"
		h.gameState.Data["ties"] = h.gameState.Data["ties"].(int) + 1
	} else if (playerMove == "rock" && computerMove == "scissors") ||
		(playerMove == "paper" && computerMove == "rock") ||
		(playerMove == "scissors" && computerMove == "paper") {
		result = "You win!"
		outcome = "win"
		h.gameState.Data["wins"] = h.gameState.Data["wins"].(int) + 1
		h.score += 5
	} else {
		result = "You lose!"
		outcome = "loss"
		h.gameState.Data["losses"] = h.gameState.Data["losses"].(int) + 1
	}
	
	emoji := map[string]string{
		"rock":     "🗿",
		"paper":    "📄",
		"scissors": "✂️",
	}

	log.Println(outcome)
	
	return fmt.Sprintf("You: %s %s\nComputer: %s %s\n%s\n\nPlay again or type 'menu' to return.",
		emoji[playerMove], playerMove,
		emoji[computerMove], computerMove,
		result)
}

func (h *GameHandler) getRPSStats() string {
	wins := h.gameState.Data["wins"].(int)
	losses := h.gameState.Data["losses"].(int)
	ties := h.gameState.Data["ties"].(int)
	total := wins + losses + ties
	
	if total == 0 {
		return "No games played yet!"
	}
	
	winRate := float64(wins) / float64(total) * 100
	
	return fmt.Sprintf("📊 RPS Statistics:\n"+
		"Wins: %d\n"+
		"Losses: %d\n"+
		"Ties: %d\n"+
		"Total: %d\n"+
		"Win Rate: %.1f%%",
		wins, losses, ties, total, winRate)
}

// Simple Quiz Game
func (h *GameHandler) startQuizGame() string {
	h.gameState.CurrentGame = "quiz"
	h.gameState.Data["currentQuestion"] = 0
	h.gameState.Data["correctAnswers"] = 0
	
	return h.getNextQuestion()
}

func (h *GameHandler) handleQuizCommand(command string, args []string) (string, uint32) {
	questions := h.getQuizQuestions()
	currentQ := h.gameState.Data["currentQuestion"].(int)
	
	if currentQ >= len(questions) {
		return h.finishQuiz(), 0
	}
	
	question := questions[currentQ]
	
	// Check answer
	answer := strings.ToLower(strings.TrimSpace(command))
	correctAnswer := strings.ToLower(question.Answer)
	
	var result string
	if answer == correctAnswer {
		result = "✅ Correct!"
		h.gameState.Data["correctAnswers"] = h.gameState.Data["correctAnswers"].(int) + 1
		h.score += 10
	} else {
		result = fmt.Sprintf("❌ Wrong! The correct answer was: %s", question.Answer)
	}
	
	h.gameState.Data["currentQuestion"] = currentQ + 1
	
	if currentQ+1 >= len(questions) {
		return result + "\n\n" + h.finishQuiz(), 0
	}
	
	return result + "\n\n" + h.getNextQuestion(), 0
}

type QuizQuestion struct {
	Question string
	Answer   string
}

func (h *GameHandler) getQuizQuestions() []QuizQuestion {
	return []QuizQuestion{
		{"What is the capital of France?", "Paris"},
		{"What is 2 + 2?", "4"},
		{"What programming language is this server written in?", "Go"},
		{"What year was the first iPhone released?", "2007"},
		{"What is the largest planet in our solar system?", "Jupiter"},
	}
}

func (h *GameHandler) getNextQuestion() string {
	questions := h.getQuizQuestions()
	currentQ := h.gameState.Data["currentQuestion"].(int)
	
	if currentQ >= len(questions) {
		return h.finishQuiz()
	}
	
	question := questions[currentQ]
	return fmt.Sprintf("🧠 QUIZ - Question %d/%d\n\n%s\n\nYour answer:",
		currentQ+1, len(questions), question.Question)
}

func (h *GameHandler) finishQuiz() string {
	questions := h.getQuizQuestions()
	correct := h.gameState.Data["correctAnswers"].(int)
	total := len(questions)
	percentage := float64(correct) / float64(total) * 100
	
	var grade string
	if percentage >= 80 {
		grade = "Excellent! 🌟"
		h.gamesWon++
	} else if percentage >= 60 {
		grade = "Good job! 👍"
	} else {
		grade = "Keep studying! 📚"
	}
	
	return fmt.Sprintf("🎓 Quiz Complete!\n\n"+
		"Score: %d/%d (%.1f%%)\n"+
		"%s\n\n"+
		"Type 'menu' to play another game.",
		correct, total, percentage, grade)
}

func (h *GameHandler) showScore() string {
	return fmt.Sprintf("🏆 Your Stats:\n"+
		"Total Score: %d points\n"+
		"Games Won: %d\n"+
		"Username: %s",
		h.score, h.gamesWon, h.username)
}

func (h *GameHandler) getHelp() string {
	return `🎮 GAME SERVER HELP 🎮

Global Commands:
- menu      Return to main menu
- score     Show your score and stats
- help      Show this help
- quit      Exit the game server

Available Games:
- guess     Number guessing game (1-100)
- rps       Rock Paper Scissors
- quiz      Trivia questions
- adventure Text-based adventure (coming soon)

Each game has its own commands. Type the game name to start!`
}

// Placeholder for adventure game
func (h *GameHandler) startAdventureGame() string {
	return "🏰 Adventure game coming soon!\nType 'menu' to try other games."
}

func (h *GameHandler) handleAdventureCommand(command string, args []string) (string, uint32) {
	return "Adventure game not implemented yet. Type 'menu' to return.", 1
}

// GetPrompt implements the CommandHandler interface
func (h *GameHandler) GetPrompt() string {
	switch h.gameState.CurrentGame {
	case "menu":
		return "🎮 game> "
	case "guess":
		return "🔢 guess> "
	case "rps":
		return "✂️ rps> "
	case "quiz":
		return "🧠 quiz> "
	default:
		return fmt.Sprintf("🎮 %s> ", h.gameState.CurrentGame)
	}
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *GameHandler) GetWelcomeMessage() string {
	return fmt.Sprintf("🎮 Welcome to the Game Server, %s! 🎮\n\n%s",
		h.username, h.showMainMenu())
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2227"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "game_server.log"

	// Create game handler (username will be set per connection)
	handler := NewGameHandler("player")

	// Create and start server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Game Server started on port 2227!")
	log.Println("Connect with: ssh -p 2227 <username>@localhost")
	log.Println("Available games: guess, rps, quiz")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down game server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
