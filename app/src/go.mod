module github.com/harshulvijay/keylogger/app/src/program

go 1.23.0

require github.com/robotn/gohook v0.41.0

replace github.com/robotn/gohook => ../../gohook

require (
	github.com/tobychui/goHidden v0.0.0-20210912041315-888f0999f674
	github.com/vcaesar/keycode v0.10.1 // indirect
)
