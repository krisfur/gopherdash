package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
   Gopher‑Dash (emoji edition + highscores + cooldown)
   --------------------------------------------------
   Endless‑runner mini-game built with Bubble Tea + Lip Gloss.

   ✦ Emoji sprites (🐹 jump‑gopher, 🪨 rock, 🟫 ground)
   ✦ Persistent high‑score stored in CWD (./.gopherdash_highscore)
   ✦ Mild speed‑up that resets every run
   ✦ Game‑over screen with 2‑second cooldown & countdown; <Q> quits anytime
   ✦ Middle pane shrinks during game‑over for a compact layout
   ✦ Controls: <W> or <Space> to jump, <Q> to quit
*/

// ----------------------------------------------------------------------------
// CONSTANTS
// ----------------------------------------------------------------------------
const (
	// timing
	startFrame      = 45 * time.Millisecond // initial ~22 FPS
	accelFactor     = 0.998                 // gentle speed‑up per tick
	cooldownSeconds = 2                     // restart delay on game‑over
	gameOverTick    = 250 * time.Millisecond

	// physics
	gravity = 1
	jumpVel = -4

	// sprites (each emoji is width‑2)
	playerChar = "🐹"
	groundChar = "🟫"
	rockChar   = "🪨"

	// gameplay
	minGapCells   = 4 // logical cells between hazards
	highScoreFile = ".gopherdash_highscore"

	// UI strings
	controlsRunning  = "W/Space = jump   Q = quit"
	controlsGameOver = "Q = quit"
)

// ----------------------------------------------------------------------------
// TYPES & GLOBALS
// ----------------------------------------------------------------------------

// scoped RNG (avoids deprecated package‑level rand)
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// tick message tagged with the run generation
type tickMsg struct{ gen int }

// obstacle in the world grid
type obstacle struct {
	x   int    // horizontal logical cell (emoji = 2 columns)
	typ string // "hole" or "rock"
}

// model holds the complete program state
type model struct {
	// terminal size
	w, h int

	// derived grid size
	gameRows int
	gameCols int

	// timing
	frameDur time.Duration
	tickGen  int // generation id; increments on every restart

	// gameplay
	dist      int
	playerY   int
	velY      int
	obstacles []obstacle

	// meta
	highScore int
	gameOver  bool
	restartAt time.Time // earliest time a restart is allowed
}

// ----------------------------------------------------------------------------
// ENTRY POINT & INITIALISATION
// ----------------------------------------------------------------------------

func initialModel() model {
	return model{
		frameDur:  startFrame,
		highScore: loadHighScore(),
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("error:", err)
	}
}

// ----------------------------------------------------------------------------
// HIGH‑SCORE PERSISTENCE
// ----------------------------------------------------------------------------

func scoreFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return highScoreFile
	}
	return filepath.Join(wd, highScoreFile)
}

func loadHighScore() int {
	data, err := os.ReadFile(scoreFilePath())
	if err != nil {
		return 0
	}
	s, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || s < 0 {
		return 0
	}
	return s
}

func saveHighScore(score int) {
	_ = os.WriteFile(scoreFilePath(), []byte(strconv.Itoa(score)), 0o644)
}

// ----------------------------------------------------------------------------
// TEA HELPERS
// ----------------------------------------------------------------------------

func tickAfter(d time.Duration, gen int) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return tickMsg{gen} })
}

// recompute grid on resize
func (m *model) recalcSizes() {
	topRows, bottomRows := 1, 1 // inner heights for HUD & control bars
	borders := 2 * 3            // three boxes, two border rows each
	m.gameRows = m.h - topRows - bottomRows - borders
	if m.gameRows < 5 {
		m.gameRows = 5
	}

	m.gameCols = (m.w - 2) / 2 // logical cells (emoji width‑2)
	if m.gameCols < 10 {
		m.gameCols = 10
	}

	m.playerY = m.gameRows - 2 // one row above ground
}

// restart a new run
func (m *model) restart() tea.Cmd {
	m.dist = 0
	m.playerY = m.gameRows - 2
	m.velY = 0
	m.obstacles = nil
	m.frameDur = startFrame
	m.gameOver = false
	m.tickGen++ // invalidate all pending ticks from previous run
	return tickAfter(m.frameDur, m.tickGen)
}

// ----------------------------------------------------------------------------
// TEA IMPLEMENTATION
// ----------------------------------------------------------------------------

func (m model) Init() tea.Cmd { return tickAfter(m.frameDur, m.tickGen) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.recalcSizes()
		// no new command
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ", "w":
			if m.gameOver {
				if time.Now().After(m.restartAt) {
					return m, m.restart()
				}
				return m, nil
			}
			if m.playerY == m.gameRows-2 {
				m.velY = jumpVel
			}
		}

	case tickMsg:
		// ignore stale ticks from previous generations
		if msg.gen != m.tickGen {
			return m, nil
		}

		if m.gameOver {
			// refresh countdown every gameOverTick
			return m, tickAfter(gameOverTick, m.tickGen)
		}
		if m.gameRows == 0 || m.gameCols == 0 {
			return m, tickAfter(m.frameDur, m.tickGen)
		}

		// --- gameplay step ---
		m.dist++

		// physics
		m.velY += gravity
		m.playerY += m.velY
		if m.playerY >= m.gameRows-2 {
			m.playerY = m.gameRows - 2
			m.velY = 0
		}

		// shift obstacles
		kept := m.obstacles[:0]
		for _, ob := range m.obstacles {
			ob.x--
			if ob.x >= -1 {
				kept = append(kept, ob)
			}
		}
		m.obstacles = kept

		// spawn new obstacle if last is far enough
		furthest := -1
		for _, ob := range m.obstacles {
			if ob.x > furthest {
				furthest = ob.x
			}
		}
		if furthest < m.gameCols-minGapCells-1 && rng.Float64() < 0.12 {
			kind := "hole"
			if rng.Float64() < 0.5 {
				kind = "rock"
			}
			spawn := m.gameCols + rng.Intn(4)
			m.obstacles = append(m.obstacles, obstacle{spawn, kind})
		}

		// collision
		for _, ob := range m.obstacles {
			if ob.x == 2 {
				switch ob.typ {
				case "hole":
					if m.playerY >= m.gameRows-2 {
						m.setGameOver()
					}
				case "rock":
					if m.playerY == m.gameRows-2 {
						m.setGameOver()
					}
				}
			}
		}

		// accelerate
		m.frameDur = time.Duration(float64(m.frameDur) * accelFactor)
		return m, tickAfter(m.frameDur, m.tickGen)
	}
	return m, nil
}

func (m *model) setGameOver() {
	m.gameOver = true
	m.restartAt = time.Now().Add(cooldownSeconds * time.Second)
	if m.dist > m.highScore {
		m.highScore = m.dist
		saveHighScore(m.highScore)
	}
}

// ----------------------------------------------------------------------------
// RENDER HELPERS
// ----------------------------------------------------------------------------

// pad right to n runes (assumes width‑1 runes)
func pad(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return string(r[:n])
	}
	return s + strings.Repeat(" ", n-len(r))
}

// build grid when game is running
func (m model) renderGame() string {
	if m.gameRows == 0 || m.gameCols == 0 {
		return ""
	}
	blank := "  "
	rows := make([][]string, m.gameRows)
	for i := range rows {
		rows[i] = make([]string, m.gameCols)
		for j := range rows[i] {
			rows[i][j] = blank
		}
	}

	groundY := m.gameRows - 1
	for x := 0; x < m.gameCols; x++ {
		rows[groundY][x] = groundChar
	}
	for _, ob := range m.obstacles {
		if ob.x < 0 || ob.x >= m.gameCols {
			continue
		}
		switch ob.typ {
		case "hole":
			rows[groundY][ob.x] = blank
		case "rock":
			if groundY-1 >= 0 {
				rows[groundY-1][ob.x] = rockChar
			}
		}
	}

	px, py := 2, m.playerY
	if py >= 0 && py < m.gameRows && px < m.gameCols {
		rows[py][px] = playerChar
	}

	lines := make([]string, m.gameRows)
	for i, cells := range rows {
		var b strings.Builder
		for _, c := range cells {
			b.WriteString(c)
		}
		lines[i] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ----------------------------------------------------------------------------
// VIEW
// ----------------------------------------------------------------------------

func (m model) View() string {
	if m.w < 4 || m.h < 4 {
		return "Resizing…"
	}

	border := lipgloss.NormalBorder()

	// top HUD
	hud := lipgloss.NewStyle().Border(border).Width(m.w).
		Align(lipgloss.Left).Render(pad(fmt.Sprintf("Distance: %d", m.dist), m.w-2))

	var centerPane, ctrl string

	if m.gameOver {
		// remaining cooldown seconds (ceil)
		countdown := max(int(math.Ceil(time.Until(m.restartAt).Seconds())), 0)

		lines := []string{
			"Game over!",
			fmt.Sprintf("Distance: %d", m.dist),
			fmt.Sprintf("High score: %d", m.highScore),
		}
		if countdown > 0 {
			lines = append(lines, fmt.Sprintf("You can go again in %d…", countdown))
		} else {
			lines = append(lines, "Press Space to go again")
		}
		msg := strings.Join(lines, "\n")

		inner := lipgloss.NewStyle().Align(lipgloss.Center).
			Height(7).Width(m.w - 2).Render(msg)
		centerPane = lipgloss.NewStyle().Border(border).Width(m.w).Render(inner)

		ctrl = lipgloss.NewStyle().Border(border).Width(m.w).
			Align(lipgloss.Left).Render(pad(controlsGameOver, m.w-2))
	} else {
		centerPane = lipgloss.NewStyle().Border(border).Width(m.w).
			Render(m.renderGame())
		ctrl = lipgloss.NewStyle().Border(border).Width(m.w).
			Align(lipgloss.Left).Render(pad(controlsRunning, m.w-2))
	}

	return strings.Join([]string{hud, centerPane, ctrl}, "\n")
}
