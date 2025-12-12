package cmd

import (
	"fmt"
	"os"

	"github.com/0xjuanma/golazo/internal/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golazo",
	Short: "Live football match updates in your terminal",
	Long:  `A minimalistic TUI application that displays RSS feeds for football matches in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(app.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running application: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here
}
