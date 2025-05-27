# Game Server SSH Example

This example demonstrates an SSH server that hosts multiple text-based games, providing an interactive gaming experience over SSH.

## Features

- Multiple games in one server
- Score tracking and statistics
- Interactive game menus
- Number guessing game
- Rock Paper Scissors
- Trivia quiz
- User-friendly prompts and feedback

## Setup

1. Generate SSH keys:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect and start playing:
   ```bash
   ssh -p 2227 player@localhost
   ```

## Available Games

### 1. Number Guessing Game (`guess`)
- Guess a number between 1 and 100
- 7 attempts to find the correct number
- Hints available (even/odd, range)
- Score based on attempts used

**Commands:**
- `<number>` - Make a guess
- `hint` - Get a helpful hint
- `give up` - Surrender and see the answer

### 2. Rock Paper Scissors (`rps`)
- Classic rock-paper-scissors against computer
- Win statistics tracking
- Emoji-enhanced display
- Continuous play

**Commands:**
- `rock`, `r` - Play rock
- `paper`, `p` - Play paper
- `scissors`, `s` - Play scissors
- `stats` - Show win/loss statistics

### 3. Trivia Quiz (`quiz`)
- Multiple-choice knowledge questions
- Score tracking
- Immediate feedback
- Grade calculation

**Commands:**
- Type your answer to each question
- Questions cover various topics

### 4. Text Adventure (`adventure`)
- Coming soon!
- Will feature story-based gameplay

## Global Commands

Available in any game:
- `menu` - Return to main menu
- `score` - Show your total score and stats
- `help` - Show help information
- `quit` - Exit the game server

## Example Gaming Session

```
$ ssh -p 2227 player@localhost
🎮 Welcome to the Game Server, player! 🎮

🎮 GAME SERVER MAIN MENU 🎮

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

Choose a game by typing its name!

🎮 game> guess
🔢 NUMBER GUESSING GAME 🔢

I'm thinking of a number between 1 and 100.
You have 7 attempts to guess it!

Commands:
- <number>  Make a guess
- hint      Get a hint
- give up   Surrender
- menu      Return to main menu

Make your first guess!

🔢 guess> 50
📉 Too high! You have 6 attempts left.
🔢 guess> 25
📈 Too low! You have 5 attempts left.
🔢 guess> hint
💡 Hint: The number is odd. It's in the range 26-50.
🔢 guess> 35
📉 Too high! You have 4 attempts left.
🔢 guess> 30
📈 Too low! You have 3 attempts left.
🔢 guess> 33
🎉 Congratulations! You guessed it in 6 attempts!
You earned 20 points!
Type 'menu' to play another game.

🔢 guess> menu
🎮 GAME SERVER MAIN MENU 🎮
[... menu repeats ...]

🎮 game> rps
✂️ ROCK PAPER SCISSORS ✂️

Commands:
- rock, r     Play rock
- paper, p    Play paper
- scissors, s Play scissors
- stats       Show game statistics
- menu        Return to main menu

Best of luck! Make your move:

✂️ rps> rock
You: 🗿 rock
Computer: ✂️ scissors
You win!

Play again or type 'menu' to return.

✂️ rps> score
🏆 Your Stats:
Total Score: 25 points
Games Won: 2
Username: player
```

## Scoring System

### Number Guessing Game
- Points = (Max Attempts - Used Attempts + 1) × 10
- Maximum: 70 points (guess in 1 try)
- Minimum: 10 points (guess in 7 tries)

### Rock Paper Scissors
- 5 points per win
- No points for ties or losses

### Trivia Quiz
- 10 points per correct answer
- Bonus for high percentage scores

## Game Features

### Interactive Prompts
- Each game has custom prompts
- Visual indicators (emojis)
- Clear command instructions

### State Management
- Games maintain their own state
- Seamless switching between games
- Progress preservation within games

### Statistics Tracking
- Total score across all games
- Games won counter
- Per-game statistics

## Customization

You can extend the game server by:

1. **Adding new games:**
   ```go
   case "newgame":
       return h.startNewGame(), 0
   ```

2. **Modifying scoring:**
   ```go
   h.score += customPoints
   ```

3. **Adding difficulty levels:**
   ```go
   h.gameState.Data["difficulty"] = "hard"
   ```

4. **Implementing save/load:**
   - Persistent user profiles
   - Game state saving
   - Leaderboards

## Use Cases

- Entertainment servers
- Educational games
- Team building activities
- Programming challenges
- Interactive tutorials
- Retro gaming experiences
