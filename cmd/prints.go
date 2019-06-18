package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func good(msg string) {
	fmt.Println(promptui.IconGood + " " + msg)
}

func warn(msg string) {
	fmt.Println(promptui.IconWarn + " " + msg)
}

func bad(msg string) {
	fmt.Println(promptui.IconBad + " " + msg)
}
