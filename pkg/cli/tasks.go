package cli

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
)

// TaskOut provides a way to pretty print a long running function with
// colors and spinners.
func TaskOut(name string, fn func() error) error {
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spin.Prefix = fmt.Sprintf("%-40s    ", name)
	if err := spin.Color("yellow"); err != nil {
		return err
	}

	spin.Start()

	err := fn()

	result := aurora.Green("\u2713")
	if err != nil {
		result = aurora.Red("\u2718")
	}

	spin.FinalMSG = fmt.Sprintf("%-40s    %s\n", name, result)
	spin.Stop()

	if err != nil {
		fmt.Printf("%s\t%s\n", aurora.Red("\u21b3"), err)
		return err
	}

	return nil
}
