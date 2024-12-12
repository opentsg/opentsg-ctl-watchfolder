//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
//
// pretry printing and other functions for semver

package semver

import (
	_ "embed"
	"log/slog"
	"runtime"
	"unicode"

	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/crayon"
)

// iterate through a string and highlight it for display on a TTY.
//
// capital letters at the start of words use the capitals `c` highlighter,
// everything else uses the body `b` highlighter.
func highlightTitleCase(str string) string {
	var crayon = crayon.Color()

	c := crayon.Success
	b := crayon.Info
	res := ""
	skipped := ""

	for _, ch := range str {
		if unicode.IsUpper(ch) && unicode.IsLetter(ch) {
			if len(skipped) > 0 {
				res += b(skipped)
			}
			res += c(string(ch))
			skipped = ""
		} else {
			skipped = skipped + string(ch)
		}
	}
	if len(skipped) > 0 {
		res += b(skipped)
	}
	return res
}

func init() {
	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
