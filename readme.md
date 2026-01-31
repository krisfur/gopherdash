# Gopher‑Dash 🐹

[![Go](https://github.com/krisfur/gopherdash/actions/workflows/go.yml/badge.svg)](https://github.com/krisfur/gopherdash/actions/workflows/go.yml)

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)    [![Go](https://img.shields.io/badge/Go-1.24.4-blue)](https://go.dev/)

> A tiny terminal endless‑runner written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss).
>
> Jump rocks, leap holes, chase the high score—all in your shell.

---

## Screencast

![screencast](screencast.gif)

---

## Installation

```bash
# Go ≥1.24.4
go install github.com/krisfur/gopherdash@latest
```

The binary ends up in `$GOBIN` (usually `~/go/bin`). Add that to your `$PATH` or run with a full path.

### From source

```bash
git clone https://github.com/krisfur/gopherdash.git
cd gopherdash
go run .
```

---

## Storing scores

Persistent high score is stored locally in `.gopherdash_highscore` in your executable's directory. 

If installing with go install it will be in your `$GOBIN` location, if compiling locally it will be right in that folder.

To reset your high score just remove the file.

---


## Controls

| Key            | Action                             |
| -------------- | ---------------------------------- |
| `Space` or `W` | Jump / **Restart** after game over |
| `Q`            | Quit immediately                   |

---

## License

MIT © 2025 [Krzysztof Furman](https://www.kfurman.dev)

---

## Expected look


```
╭─────────────────────────────────────────╮
│ Distance: 128                           │
╰─────────────────────────────────────────╯
╭─────────────────────────────────────────╮
│                                         
│            🐹                           
│🟫 🪨   🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫         
╰─────────────────────────────────────────╯
╭─────────────────────────────────────────╮
│ Space = jump   Q = quit                 │
╰─────────────────────────────────────────╯
```
