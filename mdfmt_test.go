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

func TestBlankLineBeforeTableRule_Apply(t *testing.T) {
	rule := NewBlankLineBeforeTableRule()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "table at start",
			input: "| A | B |\n| - | - |\nContent",
			want:  "\n| A | B |\n| - | - |\nContent",
		},
		{
			name:  "table after paragraph",
			input: "Paragraph.\n| A |\n|--|\n",
			want:  "Paragraph.\n\n| A |\n|--|\n",
		},
		{
			name:  "table after blank",
			input: "Paragraph.\n\n| X |\n|---|\nEnd",
			want:  "Paragraph.\n\n| X |\n|---|\nEnd",
		},
		{
			name:  "no table",
			input: "Line1\nLine2\n",
			want:  "Line1\nLine2\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := rule.Apply(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Apply(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSingleSpaceAfterListItemRule(t *testing.T) {
	rule := NewSingleSpaceAfterListItemRule()
	cases := []struct {
		name, input, want string
	}{
		{
			name:  "dash, double space",
			input: "-  two spaces",
			want:  "- two spaces",
		},
		{
			name:  "dash, single space",
			input: "- one space",
			want:  "- one space",
		},
		{
			name:  "star, double space→dash",
			input: "*  star",
			want:  "- star",
		},
		{
			name:  "star, single space→dash",
			input: "* star",
			want:  "- star",
		},
		{
			name:  "indented dash multi-space",
			input: "   -   lots",
			want:  "   - lots",
		},
		{
			name:  "indented star multi-space",
			input: "  *    spaced",
			want:  "  - spaced",
		},
		{
			name:  "tabs + star",
			input: "\t* \t\ttabbed",
			want:  "\t- tabbed",
		},
		{
			name:  "not a list",
			input: "foo*  bar",
			want:  "foo*  bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := rule.Apply(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out != tc.want {
				t.Errorf("got %q, want %q", out, tc.want)
			}
		})
	}
}
