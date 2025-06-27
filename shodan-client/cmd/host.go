package cmd

import (
	"cmp"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"main/shodan"
	"main/ui"
	"slices"
	"strconv"
	"strings"
)

type vulnData struct {
	name string
	cvss float64
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

var hostCmd = &cobra.Command{
	Use:   "host [IP]",
	Short: "Get information about an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Searching: ", args[0], "...")
		host, err := client.IPHost(args[0])
		cobra.CheckErr(err)

		generalTable := table.New().
			Border(lipgloss.HiddenBorder()).
			BorderTop(false).
			BorderBottom(false).
			Rows(parseGeneralInfo(host)...).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row%2 == 0:
					return ui.TableEvenRowStyle
				default:
					return ui.TableOddRowStyle
				}
			})

		fmt.Println(ui.BoxStyle.BorderForeground(ui.Purple).Render(lipgloss.JoinVertical(
			0,
			ui.TitleStyle.Foreground(ui.Purple).Render("General"),
			generalTable.Render(),
		)))

		for _, hostData := range host.Data {
			portTable := table.New().
				Border(lipgloss.HiddenBorder()).
				BorderTop(false).
				BorderBottom(false).
				Rows(parseHostData(hostData)...).
				StyleFunc(func(row, col int) lipgloss.Style {
					switch {
					case row%2 == 0:
						return ui.TableEvenRowStyle
					default:
						return ui.TableOddRowStyle
					}
				})

			fmt.Println(ui.BoxStyle.BorderForeground(ui.Pink).Render(lipgloss.JoinVertical(
				0,
				ui.TitleStyle.Foreground(ui.Pink).Render(fmt.Sprintf("Port: %d", hostData.Port)),
				portTable.Render(),
			)))
		}
	},
}

func parseGeneralInfo(host shodan.IPHost) [][]string {
	return [][]string{
		{"Tags:", strings.Join(host.Tags, ", ")},
		{"Domains:", strings.Join(host.Domains, ", ")},
		{"Country:", host.CountryName},
		{"City:", host.City},
		{"Organization:", host.Org},
		{"ISP:", host.ISP},
		{"Operating System:", host.OS},
	}
}

func parseHostData(hostData shodan.IPHostData) [][]string {
	// iterator moment
	vulnList := func(yield func(vulnData) bool) {
		for key, value := range hostData.Vulns {
			if !yield(vulnData{name: key, cvss: value.CVSS}) {
				return
			}
		}
	}

	sortFunc := func(vuln1, vuln2 vulnData) int {
		return cmp.Compare(vuln2.cvss, vuln1.cvss)
	}
	sortedVulns := slices.SortedFunc(vulnList, sortFunc)

	var parsedVulns []string
	for _, vuln := range sortedVulns {
		var cvssStyle lipgloss.Style
		switch {
		case vuln.cvss >= 9:
			cvssStyle = ui.Base.Foreground(ui.Red)
		case vuln.cvss >= 7:
			cvssStyle = ui.Base.Foreground(ui.Orange)
		case vuln.cvss >= 4:
			cvssStyle = ui.Base.Foreground(ui.Blue)
		default:
			cvssStyle = ui.Base.Foreground(ui.Gray)
		}

		cvssString := cvssStyle.Render(strconv.FormatFloat(vuln.cvss, 'f', 1, 64))
		vulnString := fmt.Sprintf("%s (CVSS: %s)", vuln.name, cvssString)
		parsedVulns = append(parsedVulns, vulnString)
	}

	return [][]string{
		{"Transport:", hostData.Transport},
		{"Product:", hostData.Product},
		{"Version:", hostData.Version},
		{"Vulns:", strings.Join(parsedVulns, "\n")},
	}
}
