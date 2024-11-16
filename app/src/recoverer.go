package main

import "fmt"

// Attempts to recover a function that panicked.
// Set `maxPanics` to -1 to allow infinite number of panics.
//
// Also see: https://stackoverflow.com/a/41605875
func recoverer(maxPanics int, id int, fn func()) {
	defer func() {
		// this is called when `recoverer` exits

		if err := recover(); err != nil {
			fmt.Printf("[e] Program panicked (ID: %d): %v\n", id, err)

			if maxPanics == 0 {
				// we were at our last panic
				// abort the program
				panic("Too many panics")
			} else {
				// `newMaxPanics` should be one less than the previous value, unless
				// the previous value was -1
				newMaxPanics := maxPanics
				if maxPanics != -1 {
					newMaxPanics -= 1
				}

				// start the function in a new goroutine
				go recoverer(newMaxPanics, id, fn)
			}
		}
	}()

	// call our target function
	fn()
}
