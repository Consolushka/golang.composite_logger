package main

import (
	"fmt"
)

// TIP: To run your code, right-click the code and select Run. Alternatively, click the icon in the gutter and select the Run menu item from here.
func main() {
	// TIP: Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text to see how GoLand suggests fixing the warning. Alternatively, if available, click the lightbulb to view possible fixes.
	s := "gopher"

	// TIP: To run your code, right-click the code and select Run. Alternatively, click the icon in the gutter and select the Run menu item from here.
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		// TIP: To start your debugging session, right-click your code in the editor and select the Debug option. We have set one breakpoint for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.
		fmt.Println("i =", 100/i)
	}
}
