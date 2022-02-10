# golang-events

This Go package allows you to set handler functions that run when named events occur.

## Example

```golang
package main

import (
	"fmt"
	"time"

	events "github.com/theTardigrade/golang-events"
)

func main() {
	events.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EVENT ONE RUNS")
		},
		Name:               "one",
		ShouldWaitTillDone: true,
	})

	events.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EVENT TWO RUNS")
		},
		Name:               "two",
		ShouldWaitTillDone: true,
	})

	events.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EITHER EVENT ONE OR EVENT TWO RUNS")
		},
		Names:              []string{"one", "two"},
		ShouldWaitTillDone: true,
	})

	fmt.Println("***")
	events.Run("one")
	fmt.Println("***")
	events.Run("two")
	fmt.Println("***")
	events.Run("one", "two")
	fmt.Println("***")
	events.RunAll()
	fmt.Println("***")
}
```
