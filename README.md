# goactive-cli 🐍✨

**goactive-cli** is a drop-in middleware library for Go that instantly supercharges any existing [Cobra](https://github.com/spf13/cobra) CLI with an interactive, `fzf`-style menu and dynamic UI prompts.

Instead of throwing a help menu at your users when they forget a subcommand or flag, this library intercepts the execution and presents a beautiful, interactive wizard.

It requires **zero changes** to your existing Cobra commands or logic.

## ✨ Features

* **Drop-in Replacement:** Wrap your root command in a single line of code.
* **Fuzzy Command Finder:** Automatically traverses your entire Cobra command tree (including nested commands) and presents them in an `fzf`-style searchable UI.
* **Dynamic Flag Prompting:** Inspects the selected command's `pflag` definitions and automatically generates the correct UI prompts (Booleans become Yes/No toggles, Integers get automatic number validation, Strings become text inputs).
* **Multiline Editor Support:** Natively supports large, editable text areas for things like JSON request bodies via standard Cobra annotations.
* **Zero Pollution:** Your core CLI logic remains completely unaware of the UI. If a user runs `mycli deploy --region=us-east-1`, it executes instantly and bypasses the interactive prompts.

## 📦 Installation

```bash
go get [github.com/sorrowcode/goactive-cli](https://github.com/sorrowcode/goactive-cli)
