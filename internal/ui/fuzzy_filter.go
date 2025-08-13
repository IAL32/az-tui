package ui

import (
	"fmt"
	"strings"

	"github.com/IAL32/az-tui/internal/ui/components/filter"
	"github.com/evertras/bubble-table/table"
)

// NewFuzzyFilter returns a filterFunc that performs case-insensitive fuzzy
// matching (subsequence) over the concatenation of all filterable column values.
// This is compatible with the new filter component system.
// Example wiring:
//
//	m.filterFunc = NewFuzzyFilter(m.columns)
func NewFuzzyFilter(columns []table.Column) func(table.Row, string) bool {
	return func(row table.Row, filterText string) bool {
		filterText = strings.TrimSpace(filterText)
		if filterText == "" {
			return true
		}

		// Concatenate all filterable values for this row into one string
		var b strings.Builder
		for _, col := range columns {
			if !col.Filterable() {
				continue
			}
			if v, ok := row.Data[col.Key()]; ok {
				// Unwrap StyledCell if present
				switch vv := v.(type) {
				case table.StyledCell:
					v = vv.Data
				}

				switch vv := v.(type) {
				case string:
					b.WriteString(vv)
				case fmt.Stringer:
					b.WriteString(vv.String())
				default:
					b.WriteString(fmt.Sprintf("%v", v))
				}
				b.WriteByte(' ')
			}
		}

		haystack := strings.ToLower(b.String())
		if haystack == "" {
			return false
		}

		// Support multi-token filters: "acme stl" must fuzzy-match both tokens
		for _, token := range strings.Fields(strings.ToLower(filterText)) {
			if !fuzzySubsequenceMatch(haystack, token) {
				return false
			}
		}
		return true
	}
}

// CreateFilterInput creates a new filter input for the given placeholder
// This provides integration with the new filter component system
func CreateFilterInput(placeholder string) filter.FilterInput {
	return filter.NewFilterInput(placeholder)
}

// CreateFilterState creates a new filter state with all inputs initialized
// This provides integration with the new filter component system
func CreateFilterState() filter.FilterState {
	return filter.NewFilterState()
}

// fuzzySubsequenceMatch returns true if all runes in needle appear in order
// within haystack (not necessarily contiguously). Case must be normalized by caller.
func fuzzySubsequenceMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	hi, ni := 0, 0
	hr := []rune(haystack)
	nr := []rune(needle)

	for hi < len(hr) && ni < len(nr) {
		if hr[hi] == nr[ni] {
			ni++
		}
		hi++
	}
	return ni == len(nr)
}
