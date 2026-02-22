package app

import (
	"math/rand/v2"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

const (
	confettiParticles = 35
	confettiFrames    = 18
	confettiInterval  = 50 * time.Millisecond
)

var confettiGlyphs = []string{"●", "◆", "⚜", "✦", "✧", "★", "♦"}
var confettiColors = []lipgloss.Color{
	ui.Purple, ui.Gold, ui.Green,
	ui.BrightPurple, ui.BrightGold, ui.BrightGreen,
}

type particle struct {
	x, y   float64
	vx, vy float64
	glyph  string
	color  lipgloss.Color
}

// Confetti is a particle animation triggered on issue close.
type Confetti struct {
	particles []particle
	frame     int
	width     int
	height    int
	active    bool
}

// confettiTickMsg advances the animation one frame.
type confettiTickMsg struct{}

// NewConfetti creates a confetti animation centered on the screen.
func NewConfetti(width, height int) Confetti {
	particles := make([]particle, confettiParticles)
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	for i := range particles {
		particles[i] = particle{
			x:     centerX,
			y:     centerY,
			vx:    (rand.Float64() - 0.5) * 6,
			vy:    (rand.Float64() - 0.8) * 5, // bias upward
			glyph: confettiGlyphs[rand.IntN(len(confettiGlyphs))],
			color: confettiColors[rand.IntN(len(confettiColors))],
		}
	}

	return Confetti{
		particles: particles,
		frame:     0,
		width:     width,
		height:    height,
		active:    true,
	}
}

// Tick returns a command to advance the animation.
func (c Confetti) Tick() tea.Cmd {
	if !c.active {
		return nil
	}
	return tea.Tick(confettiInterval, func(time.Time) tea.Msg {
		return confettiTickMsg{}
	})
}

// Update advances particle positions by one frame.
func (c *Confetti) Update() {
	if !c.active {
		return
	}
	c.frame++
	if c.frame >= confettiFrames {
		c.active = false
		return
	}

	gravity := 0.3
	for i := range c.particles {
		c.particles[i].x += c.particles[i].vx
		c.particles[i].y += c.particles[i].vy
		c.particles[i].vy += gravity // gravity pulls down
		// Slow horizontal movement
		c.particles[i].vx *= 0.95
	}
}

// View renders the confetti overlay. Returns empty string if not active.
func (c Confetti) View() string {
	if !c.active || c.width == 0 || c.height == 0 {
		return ""
	}

	// Build a character grid
	grid := make([][]rune, c.height)
	colors := make([][]lipgloss.Color, c.height)
	for y := range grid {
		grid[y] = make([]rune, c.width)
		colors[y] = make([]lipgloss.Color, c.width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	// Place particles
	for _, p := range c.particles {
		px := int(p.x)
		py := int(p.y)
		if px >= 0 && px < c.width && py >= 0 && py < c.height {
			runes := []rune(p.glyph)
			if len(runes) > 0 {
				grid[py][px] = runes[0]
				colors[py][px] = p.color
			}
		}
	}

	// Render grid
	var lines []string
	for y := range grid {
		var line strings.Builder
		for x := range grid[y] {
			ch := string(grid[y][x])
			if grid[y][x] != ' ' {
				style := lipgloss.NewStyle().Foreground(colors[y][x])
				line.WriteString(style.Render(ch))
			} else {
				line.WriteString(ch)
			}
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// Active returns whether the animation is still running.
func (c Confetti) Active() bool {
	return c.active
}
