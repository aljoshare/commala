package cli

import (
	"fmt"

	"github.com/aljoshare/commala/internal/validator"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/enescakir/emoji"
)

func PrintResultTable(l []*validator.ValidationResult) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#eda13e")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#eda13e")).Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("#eda13e"))
	itemStyleNeutral := lipgloss.NewStyle().Foreground(lipgloss.Color("#eda13e"))
	itemStyleGreen := lipgloss.NewStyle().Foreground(lipgloss.Color("#5fa29c"))
	itemStyleRed := lipgloss.NewStyle().Foreground(lipgloss.Color("#d25b5b"))
	t := tree.Root(fmt.Sprintf("%s %s %s", emoji.Rose, "commala - A commit linter with a lot of rice", emoji.RiceBall))
	for _, v := range l {
		vm := ""
		itemStyle := itemStyleGreen
		if v.Valid {
			vm = fmt.Sprintf("%s %s: %s", v.Validator, "✓", v.Summary)
			itemStyle = itemStyleGreen
		} else {
			vm = fmt.Sprintf("%s %s: %s", v.Validator, "✗", v.Summary)
			itemStyle = itemStyleRed
		}
		r := tree.Root(vm)
		for m := range v.Messages {
			r.Child(v.Messages[m].Message)
		}
		r.ItemStyle(itemStyle)
		t.Child(r)
	}
	t.Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyleNeutral)
	fmt.Println(t)
}
