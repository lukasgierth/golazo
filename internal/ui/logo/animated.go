// Package logo renders a GOLAZO wordmark in a stylized way.
package logo

import (
	"strings"
)

// AnimatedLogo wraps a rendered logo and reveals it progressively using a wave pattern.
// It uses logo.Render() internally and does not modify the render logic.
type AnimatedLogo struct {
	fullContent  string   // Complete output from logo.Render()
	lines        []string // Split lines for wave reveal
	revealedCols []int    // Per-line reveal progress (wave offset per line)
	charsPerTick int      // Characters to reveal per tick (derived from duration)
	waveOffset   int      // Stagger between lines starting reveal (in chars)
	totalTicks   int      // Total ticks for animation
	currentTick  int      // Current tick count
	complete     bool     // Animation finished flag
	playCount    int      // Number of times animation has played
	maxPlays     int      // Max animations (0 = infinite, 1 = once)
	maxLineWidth int      // Width of the longest line
}

// NewAnimatedLogo creates a new animated logo that wraps logo.Render().
// Parameters:
//   - version: version string to display
//   - compact: whether to render in compact mode
//   - opts: logo rendering options
//   - durationMs: total animation duration in milliseconds (e.g., 1000 for 1 second)
//   - maxPlays: number of times to play animation (1 = once, 0 = infinite)
func NewAnimatedLogo(version string, compact bool, opts Opts, durationMs int, maxPlays int) *AnimatedLogo {
	// Render the full logo once
	fullContent := Render(version, compact, opts)

	// Split into lines
	lines := strings.Split(fullContent, "\n")

	// Find max line width (for wave calculation)
	maxWidth := 0
	for _, line := range lines {
		// Count visible characters (excluding ANSI codes)
		visibleWidth := visibleLength(line)
		if visibleWidth > maxWidth {
			maxWidth = visibleWidth
		}
	}

	// Calculate animation parameters
	// Tick interval is 70ms, so ticks in duration = durationMs / 70
	const tickIntervalMs = 70
	totalTicks := durationMs / tickIntervalMs
	if totalTicks < 1 {
		totalTicks = 1
	}

	// Chars per tick = total width / total ticks
	charsPerTick := maxWidth / totalTicks
	if charsPerTick < 1 {
		charsPerTick = 1
	}

	// Wave offset: how many chars delay between each line starting
	// A smaller value creates a tighter wave, larger creates more stagger
	waveOffset := 1

	// Initialize reveal progress for each line
	revealedCols := make([]int, len(lines))

	return &AnimatedLogo{
		fullContent:  fullContent,
		lines:        lines,
		revealedCols: revealedCols,
		charsPerTick: charsPerTick,
		waveOffset:   waveOffset,
		totalTicks:   totalTicks,
		currentTick:  0,
		complete:     false,
		playCount:    0,
		maxPlays:     maxPlays,
		maxLineWidth: maxWidth,
	}
}

// Tick advances the wave reveal animation by one frame.
// Each line reveals with a slight delay from the previous line, creating a wave effect.
func (a *AnimatedLogo) Tick() {
	if a.complete {
		return
	}

	a.currentTick++

	// Calculate reveal progress for each line with wave offset
	allComplete := true
	for i := range a.lines {
		// Wave delay: line i starts after i * waveOffset chars have been revealed on line 0
		lineDelay := i * a.waveOffset
		effectiveTick := a.currentTick - (lineDelay / a.charsPerTick)

		if effectiveTick > 0 {
			// This line should be revealing
			targetChars := effectiveTick * a.charsPerTick
			lineWidth := visibleLength(a.lines[i])

			if targetChars >= lineWidth {
				a.revealedCols[i] = lineWidth
			} else {
				a.revealedCols[i] = targetChars
				allComplete = false
			}
		} else {
			// This line hasn't started yet
			a.revealedCols[i] = 0
			allComplete = false
		}
	}

	// Check if animation is complete
	if allComplete {
		a.complete = true
		a.playCount++
	}
}

// View returns the current animation frame.
// When complete, returns the full logo content.
func (a *AnimatedLogo) View() string {
	if a.complete {
		return a.fullContent
	}

	var result strings.Builder
	for i, line := range a.lines {
		if i > 0 {
			result.WriteString("\n")
		}

		revealed := a.revealedCols[i]
		if revealed <= 0 {
			// Line not started yet - output empty space to maintain layout
			result.WriteString(strings.Repeat(" ", visibleLength(line)))
			continue
		}

		// Truncate the line to revealed chars (handling ANSI codes)
		truncated := truncateToVisible(line, revealed)
		remaining := visibleLength(line) - revealed
		if remaining > 0 {
			// Pad with spaces to maintain layout
			truncated += strings.Repeat(" ", remaining)
		}
		result.WriteString(truncated)
	}

	return result.String()
}

// IsComplete returns whether the animation has finished.
func (a *AnimatedLogo) IsComplete() bool {
	return a.complete
}

// Reset resets the animation state for potential replay.
func (a *AnimatedLogo) Reset() {
	a.currentTick = 0
	a.complete = false
	for i := range a.revealedCols {
		a.revealedCols[i] = 0
	}
}

// visibleLength returns the number of visible characters in a string,
// excluding ANSI escape codes.
func visibleLength(s string) int {
	length := 0
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		length++
	}

	return length
}

// truncateToVisible truncates a string to n visible characters,
// preserving ANSI escape codes.
func truncateToVisible(s string, n int) string {
	if n <= 0 {
		return ""
	}

	var result strings.Builder
	visibleCount := 0
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			result.WriteRune(r)
			continue
		}
		if inEscape {
			result.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}

		if visibleCount >= n {
			break
		}
		result.WriteRune(r)
		visibleCount++
	}

	// Reset ANSI at the end to prevent color bleeding
	if visibleCount > 0 {
		result.WriteString("\x1b[0m")
	}

	return result.String()
}
