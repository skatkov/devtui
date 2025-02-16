package main

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/lnquy/cron"
)

func main() {
	var cronExpression string

	// @see https://gist.github.com/Aterfax/401875eb3d45c9c114bbef69364dd045
	// @see https://regexr.com/4jp54
	cronRegex := `^((((\d+,)+\d+|(\d+(\/|-|#)\d+)|\d+L?|\*(\/\d+)?|L(-\d+)?|\?|[A-Z]{3}(-[A-Z]{3})?) ?){5,7})|(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)$`

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Cron Expression").
				Value(&cronExpression).
				Validate(func(str string) error {
					// First validate with regexp
					matched, err := regexp.MatchString(cronRegex, str)
					if err != nil {
						return fmt.Errorf("Validation error: %v", err)
					}
					if !matched {
						return fmt.Errorf("invalid cron expression format")
					}

					// Then validate with cron descriptor
					expr, err := cron.NewDescriptor()
					if err != nil {
						return fmt.Errorf("invalid cron expression: %v", err)
					}
					_, err = expr.ToDescription(str, cron.Locale_en)
					if err != nil {
						return fmt.Errorf("invalid cron expression: %v", err)
					}
					return nil
				}),
		),
	)

	form.Run()

	expr, err := cron.NewDescriptor(
		cron.Use24HourTimeFormat(true),
		cron.DayOfWeekStartsAtOne(true),
	)
	if err != nil {
		fmt.Printf("Error parsing cron expression: %v\n", err)
		return
	}
	desc, err := expr.ToDescription(cronExpression, cron.Locale_en)
	if err != nil {
		fmt.Printf("Error parsing cron expression: %v\n", err)
		return
	}

	fmt.Printf("%s => %s\n", cronExpression, desc)
}
