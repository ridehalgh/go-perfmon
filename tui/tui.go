package tui

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	metrics "github.com/ridehalgh/go-perfmon/metrics"
	"github.com/ridehalgh/go-perfmon/utils"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type model struct {
	list list.Model

	monitor         metrics.SystemMonitor
	processes       []metrics.ProcessDetail
	currentCpuUsage float64
	currentMemUsage float64
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func initialModel() model {
	monitor := metrics.NewGopsutilMonitor(2 * time.Second)

	cpuUsage, err := monitor.GetCpuUsage()
	if err != nil {
		log.Printf("Error getting initial CPU usage: %v\n", err)
		cpuUsage = 0
	}

	processes, err := monitor.GetProcessDetails()
	if err != nil {
		log.Printf("Error getting processes: %v\n", err)
		processes = []metrics.ProcessDetail{}
	}

	//cpuUsageStr := fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage)

	items := make([]list.Item, len(processes))
	for i, process := range processes {
		items[i] = item{title: process.Name, description: fmt.Sprintf("PID: %d | CPU: %.2f%% | MEM: %s", process.PID, process.CPUPercent, utils.FormatBytesAuto(process.MemoryRSS))}
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = ""
	list.Styles.Title = titleStyle
	list.Styles.NoItems = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	return model{
		monitor:         monitor,
		processes:       processes,
		currentCpuUsage: cpuUsage,
		list:            list,
	}
}

func (m model) Init() tea.Cmd {
	tea.SetWindowTitle("System Performance Monitor")
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		utils.InfoLogger.Println("Tick received")
		m.currentCpuUsage, _ = m.monitor.GetCpuUsage()

		return m, tick()
	case tea.WindowSizeMsg:
		// h, v := appStyle.GetFrameSize()
		// m.list.SetSize(msg.Width-h, msg.Height-v)

		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			// if m.cursor > 0 {
			// 	m.cursor--
			// }
		case "down", "j":
			// if m.cursor < len(m.processes)-1 {
			// 	m.cursor++
			// }
		case "enter", " ":
			// _, ok := m.selected[m.cursor]
			// if ok {
			// 	delete(m.selected, m.cursor)
			// } else {
			// 	m.selected[m.cursor] = struct{}{}
			// }
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) headerView() string {
	vmStat, err := m.monitor.GetMemInfo()
	if err != nil {
		log.Printf("Error getting memory usage: %v\n", err)
	}

	memoryUsageString := fmt.Sprintf("Total: %s | Available: %s | Used: %s (%.2f%%) | Free: %s",
		utils.FormatBytesAuto(vmStat.Total),
		utils.FormatBytesAuto(vmStat.Available),
		utils.FormatBytesAuto(vmStat.Used),
		vmStat.UsedPercent,
		utils.FormatBytesAuto(vmStat.Free))

	title := titleStyle.Render(fmt.Sprintf("CPU Usage: %.2f%% | MEM Usage - %s", m.currentCpuUsage, memoryUsageString))
	line := strings.Repeat("─", max(0, m.list.Width()-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s", m.headerView(), m.list.View())
}

func InitTui() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
