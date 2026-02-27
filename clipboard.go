package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type clipboardCmd struct {
	bin  string
	args []string
}

var clipboardCandidates = []clipboardCmd{
	{bin: "pbcopy"},
	{bin: "wl-copy"},
	{bin: "xclip", args: []string{"-selection", "clipboard"}},
	{bin: "xsel", args: []string{"--clipboard", "--input"}},
}

func CopyToClipboard(text string) error {
	for _, c := range clipboardCandidates {
		path, err := exec.LookPath(c.bin)
		if err != nil {
			continue
		}

		cmd := exec.Command(path, c.args...)
		cmd.Stdin = strings.NewReader(text)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("clipboard command %s failed: %w", c.bin, err)
		}

		return nil
	}

	return fmt.Errorf("no clipboard tool found. Install one of: pbcopy, wl-copy, xclip, xsel")
}
