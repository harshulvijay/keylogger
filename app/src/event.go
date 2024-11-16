package main

import hook "github.com/robotn/gohook"

// Returns a list of all modifier keys pressed, determined by `event.Mask`
func getModifiers(event hook.Event) []string {
	var modifiers []string = []string{}
	shift, ctrl, win, alt, altGr := RAWCODE_MAP_FULL[0x10],
		RAWCODE_MAP_FULL[0x11], "WIN", RAWCODE_MAP_FULL[0x12], "ALT_GR"

	// figure out the modifier keys pressed from `event.Mask`
	//
	// basically, `AND` of `event.Mask` and the mask of the respective modifier
	// key isn't zero if
	if event.Mask != 0 {
		// shift
		if event.Mask&MOD_MASKS[shift] != 0 {
			modifiers = append(modifiers, shift)
		}

		// ctrl
		if event.Mask&MOD_MASKS[ctrl] != 0 {
			modifiers = append(modifiers, ctrl)
		}

		// windows
		if event.Mask&MOD_MASKS[win] != 0 {
			modifiers = append(modifiers, win)
		}

		// alt
		if event.Mask&MOD_MASKS[alt] != 0 {
			modifiers = append(modifiers, alt)
		}

		// alt-gr
		if event.Mask&MOD_MASKS[altGr] != 0 {
			modifiers = append(modifiers, altGr)
		}
	}

	return modifiers
}
