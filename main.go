package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	formTemplate = `
------------------
Enter your details
------------------
{{ range $item := .FormLines }}
{{$item}}
{{end}}

Up/Down to move, Ctrl+C to quit
`
)

type formModel struct {
	// labels represent the labels of the form inputs
	labels []string
	// inputs represent the current values for the corresponding labels
	inputs []string
	// index represents the current input position
	index int
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			// remove the last character from the current input
			input := m.inputs[m.index]
			if len(input) > 0 {
				m.inputs[m.index] = input[:len(input)-1]
			}
		case "up":
			// move to the previous input
			if m.index > 0 {
				m.index--
			}
		case "down", "enter":
			// move to the next input
			if m.index < len(m.labels)-1 {
				m.index++
			}
		default:
			// append the character to the current input
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
				input := m.inputs[m.index] + msg.String()
				m.inputs[m.index] = input
			}
		}
	}

	return m, nil
}

func (m formModel) View() string {
	formLines := make([]string, 0)

	for i := range m.inputs {
		line := m.labels[i] + ": " + m.inputs[i]

		if m.index == i {
			line += "_"
		}

		formLines = append(formLines, line)
	}

	templateData := map[string]interface{}{
		"FormLines": formLines,
	}

	t, err := template.New("form").Parse(formTemplate)
	if err != nil {
		return err.Error()
	}

	buff := bytes.NewBuffer([]byte{})

	err = t.Execute(buff, templateData)
	if err != nil {
		return err.Error()
	}

	return buff.String()
}

func initializeModel() formModel {
	return formModel{
		labels: []string{"First Name", "Second Name", "Email"},
		inputs: []string{"", "", ""},
		index:  0,
	}
}

func main() {
	p := tea.NewProgram(initializeModel())
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println()

	form, ok := m.(formModel)
	if ok {
		for i := range form.inputs {
			fmt.Printf("%s: %s \n", form.labels[i], form.inputs[i])
		}
	}
}
