package forms

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/umbrella-sh/um-common/logging/basic"

	gf "github.com/umbrella-sh/simply-dns-cli/internal/forms/generic_fields"
)

func RunDomainSelect(choices []string) (bool, string) {
	values := make([]any, 0)
	for _, choice := range choices {
		values = append(values, choice)
	}
	model := gf.InitGenericSelectModel(gf.GenericSelectModelInput{
		HeaderText:   fmt.Sprintf("%-*s", longestHeader, "Domain:"),
		Choices:      choices,
		Values:       values,
		InitialValue: 0,
	})
	p := tea.NewProgram(model)
	m, err := p.Run()
	if err != nil {
		log.Errorln("tea failed, ", err)
		os.Exit(1)
	}
	if m, ok := m.(gf.GenericSelectModel); ok && !m.InputCancelled() {
		return false, m.Values[m.SelectedIndex()].(string)
	}
	return true, ""
}
