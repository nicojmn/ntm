package main

import (
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	ifaces   []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	names := make([]string, len(ifaces))
	for i, iface := range ifaces {
		names[i] = iface.Name
	}
	return model{
		ifaces:   names,
		selected: make(map[int]struct{}),
		cursor:   0,
	}
}

func getTitle() string {
	return `                   ((( nn nnn  )))
                  (((  nnn  nn  )))
                  (((  nn   nn  )))
                  (((  nn   nn  )))
                   (((         )))
                     //      \\
                    //        \\
                   //          \\
                  //            \\
                 //              \\

     tt
 ((( tt    )))                        ((( mm mm mmmm  )))
(((  tttt   ))) _____  _____  _____  (((  mmm  mm  mm  )))
(((  tt     )))                      (((  mmm  mm  mm  )))
(((   tttt  )))                      (((  mmm  mm  mm  )))
 (((       )))                        (((             )))

`
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.ifaces)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	// FIXME: no trucation of the title
	s := getTitle() + "\nSelect interface(s) to use:\n\n"

	// Iterate over our choices
	for i, choice := range m.ifaces {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress space to validate interface selection\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
