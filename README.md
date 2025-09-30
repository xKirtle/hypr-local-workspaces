# hypr-local-workspaces

Make Hyprland workspaces **local per monitor** instead of global.

By default, Hyprland treats workspaces as global - meaning workspace `1` is shared across all monitors.  
`hypr-local-workspaces` provides a way to scope workspaces to each monitor, so that you get **workspaces `1–9` per monitor** instead of one global set.

This is achieved through some clever tricks (like using zero-width characters in workspace names) and works seamlessly with your existing Hyprland keybinds.

---

## Features

- Localized workspaces per monitor (`1–9` available on each).
- Supports:
    - Switching to a local workspace.
    - Moving a window to a local workspace.
    - Moving all windows to a local workspace.
- Respects the **active monitor** (where the mouse is).
- Written in **Go** - fast and lightweight.

---