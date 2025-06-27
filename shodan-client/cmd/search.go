package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"main/ui"
	"strconv"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for something",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostSearch, err := client.HostSearch(args[0])
		cobra.CheckErr(err)

		var rows [][]string
		for _, host := range hostSearch.Matches {
			row := []string{
				host.IPString,
				strconv.Itoa(host.Port),
				host.Org,
			}
			rows = append(rows, row)
		}

		resultTable := table.New().
			Border(lipgloss.HiddenBorder()).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == table.HeaderRow:
					return ui.TableHeaderStyle
				case row%2 == 0:
					return ui.TableEvenRowStyle
				default:
					return ui.TableOddRowStyle
				}
			}).
			Headers("IP", "Port", "Org").
			Rows(rows...)

		fmt.Println(resultTable)
	},
}
