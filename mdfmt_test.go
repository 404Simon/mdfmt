package main

import (
	"testing"
)

func TestBlankLineAfterHeadingRule(t *testing.T) {
	rule := NewBlankLineAfterHeadingRule()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "adds blank after heading",
			input:    "# Heading 1\nText under heading\n## Heading 2\n\nText2",
			expected: "# Heading 1\n\nText under heading\n## Heading 2\n\nText2",
		},
		{
			name:     "does not add extra blank if already present",
			input:    "# Heading 1\n\nText under heading",
			expected: "# Heading 1\n\nText under heading",
		},
		{
			name:     "heading at end of file",
			input:    "# Heading 1",
			expected: "# Heading 1\n",
		},
		{
			name:     "heading followed by blank and then text",
			input:    "# Heading 1\n\nText",
			expected: "# Heading 1\n\nText",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Apply(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, got)
			}
		})
	}
}

func TestInlineMathRule(t *testing.T) {
	rule := NewInlineMathReplaceRule()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple inline math",
			input:    "Here is math: \\( x + y \\)",
			expected: "Here is math: $x + y$",
		},
		{
			name:     "math with spaces",
			input:    "Here is math: \\(   x + y   \\)",
			expected: "Here is math: $x + y$",
		},
		{
			name:     "multiple formulas",
			input:    "First: \\( a \\), second: \\( b \\)",
			expected: "First: $a$, second: $b$",
		},
		{
			name:     "no inline math",
			input:    "No formulas here.",
			expected: "No formulas here.",
		},
		{
			name:     "math with LaTeX command",
			input:    "Formula: \\( \\text{likes(Andrew, Jane)} \\)",
			expected: "Formula: $\\text{likes(Andrew, Jane)}$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Apply(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, got)
			}
		})
	}
}

func TestReplacementRule(t *testing.T) {
	// Rule replaces smart quotes with ASCII quotes.
	rule := NewReplacementRule("SmartQuotesToAscii", map[string]string{
		"„": `"`,
		"“": `"`,
	})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "replace both smart quotes",
			input:    `He said: „Hello, world!“`,
			expected: `He said: "Hello, world!"`,
		},
		{
			name:     "no smart quotes",
			input:    `Just "normal" quotes.`,
			expected: `Just "normal" quotes.`,
		},
		{
			name:     "mixed quotes",
			input:    `„Hello“ and "hi"`,
			expected: `"Hello" and "hi"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Apply(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, got)
			}
		})
	}
}
