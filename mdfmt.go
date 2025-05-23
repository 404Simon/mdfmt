package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Rule is any transformation over the whole document.
type Rule interface {
	// Name is used in error messages.
	Name() string
	// Apply transforms the entire document.
	Apply(content string) (string, error)
}

// Formatter applies a sequence of Rules in order.
type Formatter struct {
	rules []Rule
}

func NewFormatter(rules ...Rule) *Formatter {
	return &Formatter{rules: rules}
}

func (f *Formatter) Format(content string) (string, error) {
	var err error
	for _, r := range f.rules {
		content, err = r.Apply(content)
		if err != nil {
			return "", fmt.Errorf("rule %q failed: %w", r.Name(), err)
		}
	}
	return content, nil
}

// ----------------------------------------------------------------
// Rule 1: ensure exactly one blank line after each ATX heading
// ----------------------------------------------------------------

type BlankLineAfterHeadingRule struct{}

func NewBlankLineAfterHeadingRule() Rule { return BlankLineAfterHeadingRule{} }

func (BlankLineAfterHeadingRule) Name() string {
	return "BlankLineAfterHeading"
}

func (BlankLineAfterHeadingRule) Apply(content string) (string, error) {
	var outLines []string
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		outLines = append(outLines, line)
		if isATXHeading(line) {
			// look ahead: if next line is non‐blank or EOF, insert one blank
			if i+1 >= len(lines) || strings.TrimSpace(lines[i+1]) != "" {
				outLines = append(outLines, "")
			}
		}
	}
	return strings.Join(outLines, "\n"), nil
}

func isATXHeading(line string) bool {
	// trim leading space/tabs
	t := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(t, "#") {
		return false
	}
	// count # up to 6
	count := 0
	for _, ch := range t {
		if ch == '#' && count < 6 {
			count++
		} else {
			break
		}
	}
	if count == 0 || len(t) <= count {
		return false
	}
	// must have space or tab after the hashes
	return t[count] == ' ' || t[count] == '\t'
}

// ----------------------------------------------------------------
// Rule 2: replace \(...\) with $...$
// ----------------------------------------------------------------

type InlineMathRule struct {
	// matches literal `\(`, optional spaces, capture anything non‐greedy,
	// optional spaces, then literal `\)`
	re *regexp.Regexp
}

func NewInlineMathReplaceRule() Rule {
	// `\\\(\s*(.*?)\s*\\\)` in Go literal:
	//   \\$ → literal `\$` in replacement; here we just compile the pattern
	return InlineMathRule{
		re: regexp.MustCompile(`\\\(\s*(.*?)\s*\\\)`),
	}
}

func (InlineMathRule) Name() string {
	return "InlineMathToDollar"
}

func (r InlineMathRule) Apply(content string) (string, error) {
	// replace each `\(...\)` with `$...$`
	// replacement string: "\\$$1\\$" →
	//   \\$  → regex engine sees `\$` → emits literal `$`
	//   $1   → emits group 1
	//   \\$  → emits literal `$`
	return r.re.ReplaceAllString(content, "$$$1$"), nil
}

// ----------------------------------------------------------------
// Rule 3: Replace characters with other ones
// ----------------------------------------------------------------

type ReplacementRule struct {
	// replacements maps each unwanted string to its replacement.
	replacements map[string]string
	// name is used for identification and error messages.
	name string
}

// NewReplacementRule constructs a ReplacementRule with a name and a map of replacements.
func NewReplacementRule(name string, replacements map[string]string) Rule {
	return &ReplacementRule{name: name, replacements: replacements}
}

func (r *ReplacementRule) Name() string {
	return r.name
}

func (r *ReplacementRule) Apply(content string) (string, error) {
	// For each unwanted string, replace all its occurrences with the replacement.
	for old, new := range r.replacements {
		content = strings.ReplaceAll(content, old, new)
	}
	return content, nil
}

// ----------------------------------------------------------------

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin:", err)
		os.Exit(1)
	}

	fmter := NewFormatter(
		NewBlankLineAfterHeadingRule(),
		NewInlineMathReplaceRule(),
		NewReplacementRule("SmartQuotesToAscii", map[string]string{
			"„": `"`,
			"“": `"`,
		}),
	)

	out, err := fmter.Format(string(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// ensure trailing newline
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	fmt.Print(out)
}
