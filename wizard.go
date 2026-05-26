package wizard

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Execute wraps a standard Cobra root command with an interactive fuzzy-finder fallback.
// If the user runs the CLI without a subcommand, it launches the interactive UI.
func Execute(root *cobra.Command) error {
	targetCmd, _, err := root.Find(os.Args[1:])

	if err == nil && targetCmd == root {
		selectedArgs, selectedCmd, fuzzyErr := triggerFuzzyMenu(root)

		if fuzzyErr != nil {
			if fuzzyErr == fuzzyfinder.ErrAbort {
				fmt.Println("Canceled.")
				return nil
			}
			return fuzzyErr
		}

		flagArgs, flagErr := promptForFlags(selectedCmd)
		if flagErr != nil {
			return fmt.Errorf("prompt canceled")
		}

		finalArgs := append(selectedArgs, flagArgs...)
		root.SetArgs(finalArgs)
	}

	return root.Execute()
}

// triggerFuzzyMenu handles the fzf UI (Internal)
func triggerFuzzyMenu(root *cobra.Command) ([]string, *cobra.Command, error) {
	availableCmds := getRunnableCommands(root)

	if len(availableCmds) == 0 {
		return nil, nil, fmt.Errorf("no subcommands available")
	}

	idx, err := fuzzyfinder.Find(
		availableCmds,
		func(i int) string { return availableCmds[i].CommandPath() },
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Command: %s\n\n%s", availableCmds[i].CommandPath(), availableCmds[i].Short)
		}),
	)
	if err != nil {
		return nil, nil, err
	}

	selectedCmd := availableCmds[idx]
	args := []string{}
	current := selectedCmd
	for current != root && current != nil {
		args = append([]string{current.Name()}, args...)
		current = current.Parent()
	}

	return args, selectedCmd, nil
}

// promptForFlags dynamically generates interactive prompts based on pflag types (Internal)
func promptForFlags(cmd *cobra.Command) ([]string, error) {
	var flagArgs []string

	if !cmd.HasLocalFlags() {
		return flagArgs, nil
	}

	fmt.Printf("⚙️  Configure options for '%s':\n", cmd.Name())

	var promptErr error
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if promptErr != nil {
			return
		}

		flagType := f.Value.Type()

		if flagType == "bool" {
			var val bool
			promptErr = huh.NewConfirm().
				Title(fmt.Sprintf("Enable --%s?", f.Name)).
				Description(f.Usage).
				Value(&val).
				Run()

			if promptErr == nil && val {
				flagArgs = append(flagArgs, fmt.Sprintf("--%s=true", f.Name))
			}
			return
		}

		var validator func(string) error
		if strings.Contains(flagType, "int") {
			validator = func(v string) error {
				if v == "" {
					return nil
				}
				_, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return fmt.Errorf("must be a valid integer")
				}
				return nil
			}
		} else if strings.Contains(flagType, "float") {
			validator = func(v string) error {
				if v == "" {
					return nil
				}
				_, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("must be a valid number")
				}
				return nil
			}
		}

		val := f.DefValue
		desc := fmt.Sprintf("%s\n[Type: %s]", f.Usage, flagType)

		isMultiline := false
		if ann, ok := f.Annotations["editor"]; ok && len(ann) > 0 && ann[0] == "multiline" {
			isMultiline = true
		}

		if isMultiline {
			promptErr = huh.NewText().
				Title(fmt.Sprintf("Set --%s", f.Name)).
				Description(desc + " (Press Esc then Enter to submit)").
				Lines(8).
				Value(&val).
				Run()
		} else {
			inputPrompt := huh.NewInput().
				Title(fmt.Sprintf("Set --%s", f.Name)).
				Description(desc).
				Value(&val)

			if validator != nil {
				inputPrompt.Validate(validator)
			}
			promptErr = inputPrompt.Run()
		}

		if promptErr == nil && val != "" {
			flagArgs = append(flagArgs, fmt.Sprintf("--%s=%s", f.Name, val))
		}
	})

	return flagArgs, promptErr
}

// getRunnableCommands recursively walks the Cobra tree (Internal)
func getRunnableCommands(cmd *cobra.Command) []*cobra.Command {
	var cmds []*cobra.Command
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.Name() == "help" {
			continue
		}
		if c.Runnable() {
			cmds = append(cmds, c)
		}
		if c.HasSubCommands() {
			cmds = append(cmds, getRunnableCommands(c)...)
		}
	}
	return cmds
}
