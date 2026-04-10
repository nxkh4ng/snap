# snap

> snap your commits into shape.

`snap` is a lightweight CLI tool that helps you write clean, consistent Git commits 
following the **Conventional Commits** standard - without slowing you down.

![Demo](./demos/main.gif)

---

## Features

- Interactive commit flow (fast, minimal friction)
- Built-in validation for your commit types
- Enforces Conventional Commits format: `<type>(scope): <subject>`
- Optional fields via flags (description, footer, breaking change, ticket)
- Zero config (works out of the box)

---

## Installation

### Using Go (recommended)

```bash
go install github.com/nxkh4ng/snap@latest
```

> Requires Go 1.20+

After installation, make sure your `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

### From GitHub Releases

Download the latest binary from:
https://github.com/nxkh4ng/snap/releases

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

| Flag | Description |
| --- | --- |
| `-a`, `--all` | enable all optional fields |
| `-d`, `--desc` | add description (body) |
| `-f`, `--footer` | add footer |
| `-b`, `--breaking` | add breaking change |
| `-t`, `--ticket` | add ticket reference |

#### Ticket

![Ticket](./demos/ticket.gif)

#### Description

![Description](./demos/description.gif)

#### Footer

![Footer](./demos/footer.gif)

#### Breaking Change

![Breaking Change](./demos/breaking-change.gif)

#### All-in-once

![Breaking Change](./demos/all-in-once.gif)

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

## License

[MIT](./LICENSE)
