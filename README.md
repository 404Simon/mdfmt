# mdfmt

My markdown formatter.

## Build & Install

```bash
# Build

go build -o ~/.local/bin/mdfmt mdfmt.go

# Ensure ~/.local/bin is on your PATH

```

## Usage

```bash
# Format stdin→stdout

cat in.md | mdfmt > out.md
```

## Neovim Integration (conform.nvim)

In your Neovim Lua config:

```lua
require("conform").setup({
  formatters_by_ft = {
    markdown = { "mdfmt" },
  },
  formatters = {
    mdfmt = {
      command = "mdfmt",
      args    = {},
      stdin   = true,
    },
  },
})
```

Now `:lua require("conform").format()` or your on-save hook will run `mdfmt` on Markdown files.
