package utils

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
)

type PromptContent struct {
	ErrorMsg     string
	Label        string
	ValidateFunc func(input string) error
}

func PromptGetInput(pc PromptContent) string {
	if pc.ValidateFunc == nil {
		pc.ValidateFunc = func(input string) error {
			if len(input) <= 0 {
				return errors.New(pc.ErrorMsg)
			}
			return nil
		}
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label,
		Templates: templates,
		Validate:  pc.ValidateFunc,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func PromptGetSelect(pc PromptContent, items []string) string {
	var (
		result string
		err    error
	)

	prompt := promptui.Select{
		Label: pc.Label,
		Items: items,
	}
	_, result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}
