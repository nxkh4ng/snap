# snap

> snap your commits into shape.

`snap` is a lightweight CLI tool that helps you write clean, consistent Git commits
following the [**Conventional Commits**](https://www.conventionalcommits.org/en/v1.0.0/) standard - without slowing you down.

![Demo](./demos/main.gif)

---

## Features

- Interactive commit flow (fast, minimal friction)
- Built-in validation for your commit types
- Enforces Conventional Commits format: `<type>(scope): <subject>`
- Optional fields via flags (description, footer, breaking change, ticket)
- Scope suggestions from your git history
- Zero config (works out of the box)
- Project-level config via `snap init` — creates `.snap.json` next to `.git`
- Global config via `snap init --global` — applies across all projects

---

## Installation

### Option 1: Using Go (recommended)

Requires **Go 1.20+**. Install Go from [golang.org](https://go.dev/dl/) if you haven't already

```bash
go install github.com/nxkh4ng/snap@latest
```

> [!NOTE]
>
> - This places the `snap` binary in your Go bin directory.
> - Make sure `$GOPATH/bin` (or `%USERPROFILE%\go\bin` on Windows) is in your `PATH`.

### Option 2: Download binary

Download the latest binary from: [Release Page](https://github.com/nxkh4ng/snap/releases)

#### MacOS / Linux

```bash
# Rename and move to PATH
sudo mv snap_ /usr/local/bin/snap

# Grant execute permission
chmod +x /usr/local/bin/snap
```

#### Windows

1. Download `snap_Windows_x86_64.exe` from the Releases page
2. Rename it to `snap.exe`
3. Move it to a folder that is in your `PATH`, or add its folder to `PATH` via **System Properties → Environment Variables**

---

**Verify the installation:**

```bash
snap --version
```

---

## Usage

### Basic

```bash
snap
```

Flow: type → scope (optional) → subject

![Basic](./demos/basic.gif)

### Flags

Extend your commit with optional parts:

| Flag               | Description                |
| ------------------ | -------------------------- |
| `-a`, `--all`      | enable all optional fields |
| `-d`, `--desc`     | add description (body)     |
| `-f`, `--footer`   | add footer                 |
| `-b`, `--breaking` | add breaking change        |
| `-t`, `--ticket`   | add ticket reference       |

#### Ticket

![Ticket](./demos/ticket.gif)

#### Description

![Description](./demos/description.gif)

#### Footer

![Footer](./demos/footer.gif)

#### Breaking Change

![Breaking Change](./demos/breaking-change.gif)

#### All-in-once

Flow:

1. type -> scope(optional) -> subject
2. desc(optional)
3. breaking change
4. footer(optional)
5. ticket

![All-in-once](./demos/all-in-once.gif)

---

## Configuration

`snap` works out of the box with no config required. When you need to customize behavior per project or globally, use `snap init`.

### Local config (per project)

```bash
snap init
```

Creates `.snap.json` next to your `.git` folder. Committed to the repo so the whole team shares the same config.

### Global config

```bash
snap init --global
```

Creates `~/.config/snap/config.json`. Applies to all projects that don't have a local `.snap.json`.

### Config options

| Field              | Type   | Default               | Description                                                   |
| ------------------ | ------ | --------------------- | ------------------------------------------------------------- |
| `types`            | array  | 10 built-in types     | commit types available in wizard                              |
| `scopes`           | array  | `[]`                  | predefined scopes — leave empty for free-text                 |
| `requireScope`     | bool   | `false`               | make scope mandatory                                          |
| `subjectCharLimit` | int    | `100`                 | max length of subject line                                    |
| `ticketKeyWords`   | array  | `Closes, Fixes, Refs` | keywords for ticket reference                                 |
| `theme`            | string | `base16`              | UI theme (`base`, `base16`, `catppuccin`, `dracula`, `charm`) |

Defaut config:

```json
{
  "types": [
    {
      "name": "feat",
      "description": "A new feature"
    },
    {
      "name": "fix",
      "description": "A bug fix"
    },
    {
      "name": "chore",
      "description": "Build process or auxiliary tool changes"
    },
    {
      "name": "docs",
      "description": "Documentation only changes"
    },
    {
      "name": "style",
      "description": "Markup, white-space, formatting, missing semi-colons..."
    },
    {
      "name": "refactor",
      "description": "A code change that neither fixes a bug nor adds a feature"
    },
    {
      "name": "perf",
      "description": "A code change that improves performance"
    },
    {
      "name": "test",
      "description": "Adding missing tests"
    },
    {
      "name": "build",
      "description": "Changes that affect the build system or external dependencies"
    },
    {
      "name": "ci",
      "description": "CI related changes"
    }
  ],
  "scopes": [],
  "requireScope": false,
  "subjectCharLimit": 100,
  "ticketKeyWords": [
    {
      "name": "Closes",
      "description": "Closes the issue when merged"
    },
    {
      "name": "Fixes",
      "description": "Fixes a bug and closes the issue"
    },
    {
      "name": "Refs",
      "description": "References without closing"
    }
  ],
  "theme": "base16"
}
```

---

## Commit Format

`snap` follows:

```
<type>(scope): <subject>

[optional body]

[optional footer]
```

Example:

```
feat(auth): add login with Google

support OAuth2 flow

BREAKING CHANGE: remove legacy login
Refs: #123
```

---

## Uninstall

### Via Go

```bash
# macOS / Linux
rm "$(go env GOPATH)/bin/snap"

# Windows (PowerShell)
Remove-Item "$(go env GOPATH)\bin\snap.exe"
```

### Via Binary

Remove the binary from wherever you placed it during installation.
If you're unsure of the location, run the following to find it:

```bash
# macOS / Linux
which snap

# Windows (PowerShell)
where snap
```

Then delete the file at the path shown

---

## License

[MIT](./LICENSE)
