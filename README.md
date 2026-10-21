# Gototp

A terminal application for generating time-based one-time passwords (TOTP)
for two-factor authentication. Accounts are stored in a local encrypted
database unlocked with a passphrase.

## Features

- View codes for all accounts or periodically refresh a selected account's code.
- Add, edit, and delete accounts.
- Configurable periods, 6- or 8-digit codes, and SHA1, SHA256, SHA512, or MD5.
- Generate codes with a time offset based on a specified date and time.
- Import and export individual accounts using `.gototp` files.
- Erase temporary passphrase copies on supported platforms.

## Installation

Download an archive from the [latest release](https://github.com/psauxwwf/gototp/releases/latest):

| Platform | Archive |
| --- | --- |
| Linux amd64 | [gototp_linux_amd64.tar.gz](https://github.com/psauxwwf/gototp/releases/latest/download/gototp_linux_amd64.tar.gz) |
| Windows amd64 | [gototp_windows_amd64.tar.gz](https://github.com/psauxwwf/gototp/releases/latest/download/gototp_windows_amd64.tar.gz) |

Extract the archive and place `gototp` or `gototp.exe` in a directory on your
`PATH`. Go is not required to run a prebuilt binary.

For example, on Linux, from the directory containing the downloaded archive:

```sh
tar -xzf gototp_linux_amd64.tar.gz
install -Dm755 gototp "$HOME/.local/bin/gototp"
```

Make sure `$HOME/.local/bin` is on your `PATH`.
Separate source-build tasks are available for macOS and Android/Termux.

## Quick start

```sh
gototp
```

1. Enter a passphrase at the `enter your secret` prompt and press `Enter`.
   If the database does not exist, it is created using this passphrase.
   Use the same passphrase to unlock an existing database on subsequent runs.
2. Select `add` and enter the account details.
3. Select `show` to view current codes for all accounts, or `code` to
   periodically display a selected account's code.

Use the arrow keys to navigate options, `Tab` / `Shift+Tab` to move between
fields, and `Enter` to confirm. In fields with suggestions, `Ctrl+E` accepts
the suggested value. Press `Ctrl+C` in the main menu to exit.

> The database passphrase and a service's TOTP secret are different things.
> Keep your passphrase somewhere safe: there is no passphrase recovery mechanism.

### Adding an account

Select `add` and fill in the following fields:

| Field | Value |
| --- | --- |
| `name` | Service or account name, such as `GitHub — personal` |
| `update period` | Period in seconds, usually `30`; enter it or accept the suggestion |
| `set digits` | `6` or `8` digits, as required by the service |
| `algorithm` | The service's algorithm, usually `SHA1` |
| `secret` | The Base32 TOTP secret provided by the service during 2FA setup |

The account is saved after you complete the form. Enter the TOTP secret,
not a current one-time code. All parameters must match the service's settings.

### Main menu

| Action | Description |
| --- | --- |
| `show` | Display accounts, their parameters, and current OTP codes in a table |
| `code` | Select an account and display codes at its configured update interval |
| `add` | Add an account |
| `change` | Edit an account's parameters, not the database passphrase |
| `delete` | Delete an account |
| `export` | Export an account to `<name>.gototp` in the current directory |
| `import` | Import an account from a `.gototp` file in the current directory |
| `code for time` | Refresh codes using a time offset based on the entered date and time |

For `code for time`, enter a date and time in your local time zone using
`DD.MM.YYYY HH:MM:SS`, for example `21.10.2026 12:30:00`.
Normal TOTP generation requires an accurate system clock.

> **Exports are not encrypted.** A `.gototp` file contains the TOTP secret
> and allows code generation without the database passphrase. Treat it as
> sensitive data when storing or sharing it. Only import trusted files.

## Storage and configuration

The application reads an optional `config.yaml` and stores the database,
`gototp.enc`, in the same directory by default:

| Platform / environment | Directory |
| --- | --- |
| Linux, macOS, Termux without `XDG_CONFIG_HOME` | `~/.config/gototp/` |
| `XDG_CONFIG_HOME` set on these platforms | Directly in `$XDG_CONFIG_HOME/`, without a `gototp` subdirectory |
| Windows | `%USERPROFILE%\AppData\Roaming\gototp\` |

To use a different database location, create `config.yaml`:

```yaml
database: "~/private/gototp.enc"
```

The database's parent directory must exist. The `~/` prefix is supported;
relative paths are resolved against the application's working directory.
Pointing to a nonexistent database creates a new, empty database.

To back up your data, close the application and copy `gototp.enc` to a safe
location. Restoring it requires both the database file and its passphrase.
Individual account exports are not a substitute for an encrypted backup.

## Security

### Database encryption

The database is encrypted with AES-256-GCM. The current format derives the
key as `SHA-256(passphrase)`, without a salt or a password-stretching KDF.
Weak passphrases are therefore vulnerable to offline guessing if the database
is stolen. Choose a long, random passphrase and restrict access to the file.

### Passphrase memory protection

Password entry and database unlocking run inside experimental
[`runtime/secret.Do`](https://pkg.go.dev/runtime/secret). It erases temporary
stacks and registers on return. A forced garbage collection then erases
unreachable heap allocations made in secret mode, including temporary copies
created by the input form. The passphrase does not escape into the main menu.

In Go 1.27, this erasure is supported **only on Linux amd64/arm64**.
On Windows, macOS, and Android/Termux, `secret.Do` runs the function without
this protection; there is no passphrase-erasure guarantee on those platforms.

On Linux, before password entry, the application also sets `PR_SET_DUMPABLE=0`
and zero `RLIMIT_CORE` limits. This disables ordinary process dumps and
restricts unprivileged access to process memory. If these restrictions cannot
be applied, the application exits before requesting the passphrase.

**Limitations:**

- This is not absolute protection against root, a privileged debugger, or a
  compromised system.
- Erasure after unlocking does not protect the passphrase during input or
  decryption, or remove existing copies from swap.
- The derived encryption key and decrypted OTP data remain in memory while
  the application is running so it can generate codes and save changes.
- Exported `.gototp` files are not encrypted, and displayed codes may remain
  in terminal scrollback.

## Building from source

Install Go **1.27.3 or later** and [Task](https://taskfile.dev/installation/).
The commands below assume a POSIX environment. `task build` and `task release`
also require `strip` and `objcopy` with support for the target binary formats.

```sh
git clone https://github.com/psauxwwf/gototp.git
cd gototp
task --list
task build
```

| Command | Result |
| --- | --- |
| `task build` | Linux and Windows amd64: `bin/gototp`, `bin/gototp.exe` |
| `task build:linux` | Linux amd64 only |
| `task build:windows` | Windows amd64 only |
| `task build:darwin` | macOS amd64: `bin/gototp_darwin` |
| `task build:termux` | Android arm64: `bin/gototp_termux`; requires Android NDK r27 in `android-ndk-r27/` |
| `task zip` | Two `.tar.gz` archives from existing Linux/Windows binaries |
| `task build zip` | Build and package without publishing |
| `task dev` | Rebuild when Go files change; requires `entr` |

Despite the task name `zip`, archives use **`.tar.gz`**, not `.zip`.
macOS and Termux builds are not included in automatic packaging or publishing.

The Taskfile enables `GOEXPERIMENT=runtimesecret`. When invoking Go directly
in a POSIX shell, set it explicitly:

```sh
GOEXPERIMENT=runtimesecret go build -o gototp ./cmd/gototp
GOEXPERIMENT=runtimesecret go test ./...
GOEXPERIMENT=runtimesecret go test -race ./internal/view ./internal/security
```

Race tests require a platform supported by Go's race detector and an available
C compiler. Being able to build with the experimental package does not imply
that memory erasure is supported on the target platform.

## Publishing a release

Install [GitHub CLI](https://cli.github.com/), push the intended changes to
GitHub, and authenticate with permission to create releases:

```sh
gh auth login
task release TAG=v1.2.3
```

You can set `GH_TOKEN` instead of authenticating interactively.
The task checks authentication, builds and packages the binaries, and creates
a published GitHub Release with automatically generated release notes and
two assets:

- `gototp_linux_amd64.tar.gz`
- `gototp_windows_amd64.tar.gz`

`TAG` is required. If the tag does not exist, GitHub creates it from the
current state of the default branch, not from local uncommitted changes.
Creating a release that already exists fails: the task does not automatically
replace its assets.
