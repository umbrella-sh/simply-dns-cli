package forms

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/umbrella-sh/um-common/logging/basic"

	gf "github.com/umbrella-sh/simply-dns-cli/internal/forms/generic_fields"
)

func RunCommentInput() (bool, string) {
	p := tea.NewProgram(gf.InitGenericInputModel(gf.GenericInputModelInput{
		HeaderText:      fmt.Sprintf("%-*s", longestHeader, "Comment:"),
		PlaceHolderText: "",
		ValueCharLimit:  255,
		IsRequired:      false,
		InputValidator:  validateDataInput,
		InputConverter:  nil,
	}))
	m, err := p.Run()
	if err != nil {
		log.Errorln("tea failed, ", err)
		os.Exit(1)
	}
	if m, ok := m.(gf.GenericInputModel); ok && !m.InputCancelled() {
		return false, m.GetValue()
	}
	return true, ""
}

func validateCommentInput(val string, required bool, valueConverter gf.GenericInputConverter) (ok bool, msg string) {
	if !required && val == "" {
		return true, ""
	}

	if len(val) > 255 {
		return false, "Comment cannot be longer than 255 chars"
	}

	return true, ""
}
