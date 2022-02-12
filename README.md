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
	manager := events.NewManager()

	manager.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EVENT ONE RUNS")
		},
		Name:               "one",
		ShouldWaitTillDone: true,
	})

	manager.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EVENT TWO RUNS")
		},
		Name:               "two",
		ShouldWaitTillDone: true,
	})

	manager.Add(events.AddOptions{
		Handler: func() {
			fmt.Println("THIS HANDLER IS CALLED WHEN EITHER EVENT ONE OR EVENT TWO RUNS")
		},
		Names:              []string{"one", "two"},
		ShouldWaitTillDone: true,
	})

	fmt.Println("***")
	manager.Run("one")
	fmt.Println("***")
	manager.Run("two")
	fmt.Println("***")
	manager.Run("one", "two")
	fmt.Println("***")
	manager.RunAll()
	fmt.Println("***")
}
```
