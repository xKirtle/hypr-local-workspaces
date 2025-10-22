# hypr-local-workspaces

[![AUR](https://img.shields.io/aur/version/hypr-local-workspaces)](https://aur.archlinux.org/packages/hypr-local-workspaces)
[![AUR votes](https://img.shields.io/aur/votes/hypr-local-workspaces)](https://aur.archlinux.org/packages/hypr-local-workspaces)
[![Release](https://img.shields.io/github/v/release/xKirtle/hypr-local-workspaces?sort=semver)](https://github.com/xKirtle/hypr-local-workspaces/releases)
[![codecov](https://codecov.io/github/xKirtle/hypr-local-workspaces/graph/badge.svg?token=A75GB31MAX)](https://codecov.io/github/xKirtle/hypr-local-workspaces)
[![Go Version](https://img.shields.io/github/go-mod/go-version/xKirtle/hypr-local-workspaces)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/xKirtle/hypr-local-workspaces)](https://goreportcard.com/report/github.com/xKirtle/hypr-local-workspaces)
[![License: MIT](https://img.shields.io/github/license/xKirtle/hypr-local-workspaces)](LICENSE)
[![AUR Publish](https://github.com/xKirtle/hypr-local-workspaces/actions/workflows/aur-publish.yml/badge.svg)](https://github.com/xKirtle/hypr-local-workspaces/actions/workflows/aur-publish.yml)

Make Hyprland workspaces local per monitor instead of global.

By default, Hyprland treats workspaces as global - meaning workspace `1` is shared across all monitors.
`hypr-local-workspaces` scopes workspaces to each monitor, so you get workspaces `1–9` per monitor instead of one global set.

This is achieved by using zero-width characters in workspace names and works seamlessly with your existing Hyprland keybinds.

## Features

- Localized workspaces per monitor (`1-9` available on each).
- Supports:
  - Switching to a local workspace.
  - Moving a window to a local workspace.
  - Moving all windows to a local workspace.
- Respects the **active monitor** (where the mouse is).
- Written in Go - fast and lightweight.

## Installation

Quick local install (recommended):

```bash
# From the repo root
./release-locally.sh
```

This script builds the binary, installs it to `/usr/bin/hypr-local-workspaces` (requires `sudo`), and removes the temporary build artifact.

Build from source manually (no sudo):

```bash
go build -o hypr-local-workspaces ./cmd/hypr-local-workspaces
install -Dm755 hypr-local-workspaces ~/.local/bin/hypr-local-workspaces
```

Arch Linux (AUR):

```bash
# With an AUR helper
paru -S hypr-local-workspaces
# or
yay -S hypr-local-workspaces

# Manually (without a helper)
git clone https://aur.archlinux.org/hypr-local-workspaces.git
cd hypr-local-workspaces
makepkg -si
```

Note: the `PKGBUILD` in this repository is updated dynamically during the release pipeline and is not intended for local `makepkg -si`. Use the AUR package above, or the local script/manual build.

## Usage

Just bind your workspace keys to the installed binary:

```bash
# Cycle existing workspaces on focused monitor
bind = $mainMod, Tab, exec, hypr-local-workspaces cycle next
bind = $mainMod SHIFT, Tab, exec, hypr-local-workspaces cycle prev

# Switch workspaces relative to active monitor
bind = $mainMod, 1, exec, hypr-local-workspaces goto 1
...

# Move active window to a workspace relative to active monitor
bind = $mainMod SHIFT, 1, exec, hypr-local-workspaces move 1
...

# Move all windows of active workspace to a workspace relative to active monitor
bind = $mainMod CTRL, 1, exec, hypr-local-workspaces move --all 1
...
```

### CLI reference

Commands and flags are structured as:

```text
hypr-local-workspaces goto  <1..9> [global flags]
hypr-local-workspaces move  <1..9> [--all] [global flags]
hypr-local-workspaces cycle <next|prev> [global flags]
```

Global flags must appear after the subcommand’s own args/flags.

- Global flags:
  - `--no-compact` - disable compact mode (enabled by default). When compact mode is enabled, the tool keeps local workspaces contiguous on each monitor by renaming zero-width workspace names as needed.

Examples:

```bash
# Switch to local workspace 3 on the focused monitor
hypr-local-workspaces goto 3

# Same, but without compacting workspaces
hypr-local-workspaces goto 3 --no-compact

# Move active window to local workspace 2
hypr-local-workspaces move 2

# Move ALL windows from the active workspace to local workspace 2
hypr-local-workspaces move --all 2

# ...and skip compaction for that move
hypr-local-workspaces move --all 2 --no-compact

# Cycle through existing local workspaces on the focused monitor
hypr-local-workspaces cycle next
hypr-local-workspaces cycle prev --no-compact
```

### What is “compaction”?

- Compaction keeps local workspaces contiguous on each monitor by renaming the internal zero‑width workspace names to remove gaps (e.g., when you close/move windows and leave empty slots in between).
- With compaction enabled (default), operations like `goto`, `move`, and `cycle` will normalize names so your local workspace sequence is effectively 1..N without holes.

Why care?

- Performance: on modern hardware this is negligible, but compaction does add a few extra Hyprland operations (fetching workspaces and occasionally renaming). If you want to minimize unnecessary work, you can disable it with `--no-compact`.
- UX: many bars and indicators expect clean, contiguous indices. Compaction keeps things tidy if you display workspace numbers.

When you might disable compaction

- If you’re using Waybar’s `hyprland/workspaces` module and your styling hides the workspace name/number, you may not care about tidy numbering. In that case, use `--no-compact` to skip renaming and reduce extra operations.
- If you prefer to keep existing names as-is and avoid any renames during navigation or moving windows.

## Development

Build:

```bash
go build -o hypr-local-workspaces ./cmd/hypr-local-workspaces
```

Run tests:

```bash
cd cmd/hypr-local-workspaces
go test ./...
```

## Contributing

Contributions are welcome! If you’d like to help:

- Open an issue to discuss bugs, improvements, or new features.
- Submit a pull request with focused changes and passing tests.
- Keep the CLI behavior and README in sync when you add/change flags.

## License

MIT - see [LICENSE](LICENSE) for details.
