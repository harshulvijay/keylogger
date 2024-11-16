// taken from https://github.com/MarinX/keylogger/blob/master/keymapper.go
// (with modifications)
package main

// Keyboard mappings taken from
// https://github.com/boppreh/keyboard/blob/master/keyboard/_winkeyboard.py
//
// TODO: make sure these work on Linux as well; otherwise, create mappings for
// Linux
var RAWCODE_MAP_FULL = map[uint16]string{
	0x03: "CONTROL_BREAK_PROCESSING",
	0x08: "BS",
	0x09: "TAB",
	0x0c: "CLEAR",
	0x0d: "ENTER",
	0x10: "SHIFT",
	0x11: "CTRL",
	0x12: "ALT",
	0x13: "PAUSE",
	0x14: "CAPS_LOCK",
	0x15: "IME_KANA_MODE",
	0x17: "IME_JUNJA_MODE",
	0x18: "IME_FINAL_MODE",
	0x19: "IME_KANJI_MODE",
	0x1b: "ESC",
	0x1c: "IME_CONVERT",
	0x1d: "IME_NONCONVERT",
	0x1e: "IME_ACCEPT",
	0x1f: "IME_MODE_CHANGE_REQUEST",
	0x20: "SPACEBAR",
	0x21: "PGUP",
	0x22: "PGDOWN",
	0x23: "END",
	0x24: "HOME",
	0x25: "LEFT",
	0x26: "UP",
	0x27: "RIGHT",
	0x28: "DOWN",
	0x29: "SELECT",
	0x2a: "PRINT",
	0x2b: "EXECUTE",
	0x2c: "PRT_SCR",
	0x2d: "INSERT",
	0x2e: "DELETE",
	0x2f: "HELP",
	0x5b: "L_WIN",
	0x5c: "R_WIN",
	0x5d: "APPLICATIONS",
	0x5f: "SLEEP",
	0x6c: "SEPARATOR",
	0x6e: "DECIMAL",
	0x70: "F1",
	0x71: "F2",
	0x72: "F3",
	0x73: "F4",
	0x74: "F5",
	0x75: "F6",
	0x76: "F7",
	0x77: "F8",
	0x78: "F9",
	0x79: "F10",
	0x7a: "F11",
	0x7b: "F12",
	0x7c: "F13",
	0x7d: "F14",
	0x7e: "F15",
	0x7f: "F16",
	0x80: "F17",
	0x81: "F18",
	0x82: "F19",
	0x83: "F20",
	0x84: "F21",
	0x85: "F22",
	0x86: "F23",
	0x87: "F24",
	0x90: "NUM_LOCK",
	0x91: "SCROLL_LOCK",
	0xa0: "L_SHIFT",
	0xa1: "R_SHIFT",
	0xa2: "L_CTRL",
	0xa3: "R_CTRL",
	// idk what to do about menu keys
	// (these two were originally `L_MENU` and `R_MENU`)
	0xa4: "L_ALT",
	0xa5: "R_ALT",
	0xa6: "BROWSER_BACK",
	0xa7: "BROWSER_FORWARD",
	0xa8: "BROWSER_REFRESH",
	0xa9: "BROWSER_STOP",
	0xaa: "BROWSER_SEARCH_KEY",
	0xab: "BROWSER_FAVORITES",
	0xac: "BROWSER_START_AND_HOME",
	0xad: "VOL_MUTE",
	0xae: "VOL_DOWN",
	0xaf: "VOL_UP",
	0xb0: "NEXT_TRACK",
	0xb1: "PREVIOUS_TRACK",
	0xb2: "STOP_MEDIA",
	0xb3: "PLAY/PAUSE_MEDIA",
	0xb4: "START_MAIL",
	0xb5: "SELECT_MEDIA",
	0xb6: "START_APPLICATION 1",
	0xb7: "START_APPLICATION 2",
	0xe5: "IME_PROCESS",
	0xf6: "ATTN",
	0xf7: "CRSEL",
	0xf8: "EXSEL",
	0xf9: "ERASE_EOF",
	0xfa: "PLAY",
	0xfb: "ZOOM",
	0xfc: "RESERVED ",
	0xfd: "PA1",
	0xfe: "CLEAR",
}

// Masks which, when `&` with `event.Mask`, return a non-zero value for the
// respective modifier key which is pressed
var MOD_MASKS map[string]uint16 = map[string]uint16{
	"SHIFT":  1 << 0,
	"CTRL":   1 << 1,
	"WIN":    1 << 2,
	"ALT":    1 << 3,
	"ALT_GR": 1 << 7,
}

// `RAWCODE_MAP` is `RAWCODE_MAP_FULL`, but it excludes keys that can be
// detected using the `keydown` event
//
// This is initialized after `initializeMaps` is called
var RAWCODE_MAP = map[uint16]string{}

// Copies `src` to `dest`
func copyMap[T comparable, U any](dest map[T]U, src map[T]U) {
	for key, value := range src {
		dest[key] = value
	}
}

// Copies `RAWCODE_MAP_FULL` to `RAWCODE_MAP` and removes unintended keys
func initializeRawcodeMap() {
	copyMap(RAWCODE_MAP, RAWCODE_MAP_FULL)

	delete(RAWCODE_MAP, 0x08) // BS
	delete(RAWCODE_MAP, 0x09) // TAB
	delete(RAWCODE_MAP, 0x0d) // ENTER
	delete(RAWCODE_MAP, 0x1b) // ESC
	delete(RAWCODE_MAP, 0x20) // SPACEBAR
	delete(RAWCODE_MAP, 0x6e) // DECIMAL
}
