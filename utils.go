package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// var GRUV_BLACK string = "#282828"
// var GRUV_RED string = "#cc241d"
// var GRUV_GREEN string = "#98971a"
// var GRUV_YELLOW string = "#d78821"
// var GRUV_BLUE string = "#458588"
// var GRUV_PURPLE string = "#b16286"
// var GRUV_CYAN string = "#689d6a"
// var GRUV_GRAY string = "#a89984"

var BLACK string = "0"
var RED string = "1"
var GREEN string = "2"
var YELLOW string = "3"
var BLUE string = "4"
var PURPLE string = "5"
var CYAN string = "6"
var GRAY string = "7"

func Color(text string, color string) string {
	if color[0] == '#' {
		r := color[1:3]
		g := color[3:5]
		b := color[5:7]
		return fmt.Sprintf(
			"\x1b[38;2;%d;%d;%dm%s\x1b[0m",
			UnwrapOr(strconv.ParseUint(r, 16, 8))(0),
			UnwrapOr(strconv.ParseUint(g, 16, 8))(0),
			UnwrapOr(strconv.ParseUint(b, 16, 8))(0),
			text,
		)
	} else {
		return fmt.Sprintf("\x1b[3%sm%s\x1b[0m", color, text)
	}
}

func Italic(text ...string) string {
	return "\x1b[3m" + strings.Join(text, " \x1b[3m") + "\x1b[0m"
}

func Bold(text ...string) string {
	return "\x1b[1m" + strings.Join(text, " \x1b[1m") + "\x1b[0m"
}

func CreatePreferencesFile() {
	whoamiCmd := exec.Command("whoami")
	user := string(Unwrap(whoamiCmd.CombinedOutput()))
	user = strings.ReplaceAll(user, "\n", "")
	prefPath := fmt.Sprintf("/home/%s/.config/gopher", user)
	prefFile := prefPath + "/Preferences"
	mkdirCmd := exec.Command("mkdir", "-p", prefPath)
	Unwrap(mkdirCmd.CombinedOutput())
	f := Unwrap(os.Create(prefFile))
	defer f.Close()
	f.WriteString(DEFAULT_PREFERENCES)
	fmt.Printf("%s %s",
		Color(Bold("Preferences"), GRAY),
		Color("file created at ~/.config/gopher\n", GRAY),
	)
}

type Preference int

const (
	PrefPkgQueryLimit = iota
	PrefOpArchPairs
	PrefPrettyPrint
	PrefPrettyPrintPreviewLines
)

var preferenceName = map[Preference]string{}

func (p Preference) String() string {
	return preferenceName[p]
}

func GetPreference[T any](name Preference) (T, error) {
	var result any
	result = 0
	whoamiCmd := exec.Command("whoami")
	user := string(Unwrap(whoamiCmd.CombinedOutput()))
	user = strings.ReplaceAll(user, "\n", "")
	prefPath := fmt.Sprintf("/home/%s/.config/gopher", user)
	prefFile := prefPath + "/Preferences"
	fileContent := string(UnwrapOrElse(os.ReadFile(prefFile))(func() []byte {
		CreatePreferencesFile()
		return []byte(DEFAULT_PREFERENCES)
	}))
	prefLines := strings.Split(fileContent, "\n")
	// grab the file data from here
	prefMap := make(map[string]string)
	for _, prefLine := range prefLines {
		prefLine = strings.Trim(prefLine, " \n\t")
		if prefLine == "" {
			continue
		}
		kvPair := strings.Split(prefLine, "=")
		if len(kvPair) <= 1 {
			result = 0
			return result.(T), errors.New(
				"no value found for key " +
					Bold(kvPair[0]) +
					" in ~/.config/gopher/Preferences",
			)
		}
		kvPair[0] = strings.Trim(kvPair[0], " \n\t")
		kvPair[1] = strings.Trim(kvPair[1], " \n\t")
		prefMap[kvPair[0]] = kvPair[1]
	}

	switch name {
	case PrefPkgQueryLimit:
		result = 0
		if prefMap["PkgQueryLimit"] == "" {
			return result.(T), errors.New(
				"no value found for key " +
					Bold("PkgQueryLimit") +
					" in ~/.config/gopher/Preferences",
			)
		}
		r, err := strconv.Atoi(prefMap["PkgQueryLimit"])
		result = r
		if err != nil || result.(int) > 100 || result.(int) < 1 {
			return result.(T), errors.New(
				"non-numeric or integer value out side of range 1-100 found for key " +
					Bold("PkgQueryLimit") +
					" in ~/.config/gopher/Preferences",
			)
		}

		return result.(T), nil
	case PrefOpArchPairs:
		result = [][]string{}
		if prefMap["OpArchPairs"] == "" {
			return result.(T), errors.New(
				"no value found for key " +
					Bold("OpArchPairs") +
					" in ~/.config/gopher/Preferences",
			)
		}
		oparshArr := strings.Split(prefMap["OpArchPairs"], ",")
		for i := 1; i < len(oparshArr); i += 2 {
			op := strings.Trim(oparshArr[i-1], " \t\n")
			arch := strings.Trim(oparshArr[i], " \t\n")
			result = append(result.([][]string), []string{op, arch})
		}

		return result.(T), nil
	case PrefPrettyPrint:
		result = false
		if prefMap["PrettyPrint"] == "" {
			return result.(T), errors.New(
				"no value found for key " +
					Bold("PrettyPrint") +
					" in ~/.config/gopher/Preferences",
			)
		}
		if prefMap["PrettyPrint"] == "true" {
			result = true
			return result.(T), nil
		} else if prefMap["PrettyPrint"] == "false" {
			return result.(T), nil
		} else {
			return result.(T), errors.New(
				"non-boolean value found for key " +
					Bold("PrettyPrint") +
					" in ~/.config/gopher/Preferences",
			)
		}
	case PrefPrettyPrintPreviewLines:
		result = 0
		if prefMap["PrettyPrintPreviewLines"] == "" {
			return result.(T), errors.New(
				"no value found for key " +
					Bold("PrettyPrintPreviewLines") +
					" in ~/.config/gopher/Preferences",
			)
		}
		r, err := strconv.Atoi(prefMap["PrettyPrintPreviewLines"])
		result = r
		if err != nil || r < 0 {
			return result.(T), errors.New(
				"non-numeric or negative integer value found for key " +
					Bold("PrettyPrintPreviewLines") +
					" in ~/.config/gopher/Preferences",
			)
		}

		return result.(T), nil
	default:
		return result.(T), errors.New(fmt.Sprintf("preference key '%s' not found", name))
	}
}

func Unwrap[T any](val T, err error) T {
	if err != nil {
		handleErr(&CmdError{CmdInvalid, err.Error()})
	}

	return val
}

func UnwrapOr[T any](val T, err error) func(T) T {
	if err != nil {
		return func(d T) T {
			return d
		}
	} else {
		return func(_ T) T {
			return val
		}
	}
}

func UnwrapOrElse[T any](val T, err error) func(func() T) T {
	if err != nil {
		return func(fn func() T) T {
			return fn()
		}
	} else {
		return func(_ func() T) T {
			return val
		}
	}

}

func Expect(err error) {
	if err != nil {
		handleErr(&CmdError{CmdInvalid, err.Error()})
	}
}

func handleErr(err *CmdError) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s %s\n", Bold(Color("error:", RED)), err.Msg)
	if err.Type != CmdRun && err.Type != CmdBuild {
		fmt.Fprintf(os.Stderr, "\nrun %s for usage\n", Italic("gopher", Color("help", BLUE)))
	}
	os.Exit(1)
}
