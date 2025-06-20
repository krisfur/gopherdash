# Gopher‑Dash 🐹⛷️

[![Go](https://github.com/krisfur/gopherdash/actions/workflows/go.yml/badge.svg)](https://github.com/krisfur/gopherdash/actions/workflows/go.yml)

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)    [![Go](https://img.shields.io/badge/Go-1.24.4-blue)](https://go.dev/)

> A tiny terminal endless‑runner written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss).
>
> Jump rocks, leap holes, chase the high score—all in your shell.

---

## Features

* Emoji sprites (`🐹`, `🪨`, `🟫`) with double‑width handling
* Adaptive layout: resizes to any terminal window
* Gentle speed ramp with per‑run reset
* Persistent high score stored locally in `.gopherdash_highscore` in your executable's directory
* Game‑over cooldown & restart (`Space`)

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

## Controls

| Key            | Action                             |
| -------------- | ---------------------------------- |
| `Space` or `W` | Jump / **Restart** after game over |
| `Q`            | Quit immediately                   |

---

## How to Play

1. The hamster (`🐹`) stays in the centre; the world scrolls left.
2. Press **Space** / **W** to hop over rocks (`🪨`) or holes (`🟫`).
3. Distance increases every tick; speed **slowly** ramps up.
4. Collide once and it’s **Game Over**—your distance compares to the high score.
5. Wait the 2‑second countdown, then hit **Space** to dash again.

---

## High Score File

The game writes/reads a plain‑text integer from:

```
.gopherdash_highscore
```

It lives in whatever directory you launch the game from, so it vanishes if you move or delete the project folder. Feel free to add it to `.gitignore`.

---

## Contributing

PRs welcome! Bug fixes, difficulty tweaks, new themes—go for it.

1. Fork & clone
2. `git checkout -b feature/my‑thing`
3. Hack away, keep the `go test` green
4. Open a pull request

---

## License

MIT © 2025 [Krzysztof Furman](https://www.kfurman.dev)

---

## Expected look


```
╭─────────────────────────────────────────╮
│ Distance: 128                          │
╰─────────────────────────────────────────╯
╭─────────────────────────────────────────╮
│                                         │
│            🐹                           │
│🟫 🪨   🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫              │
╰─────────────────────────────────────────╯
╭─────────────────────────────────────────╮
│ Space = jump   Q = quit                │
╰─────────────────────────────────────────╯
```
