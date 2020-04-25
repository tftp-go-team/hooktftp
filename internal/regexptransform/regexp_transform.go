package regexptransform

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/tftp-go-team/hooktftp/internal/config"
)

var NO_MATCH = errors.New("No match")
var BAD_GROUPS = errors.New("Regexp has too few groups")

type Escape func(string, config.HookExtraArgs) (string, error)
type Transform func(string) (string, error)

var fieldPat = regexp.MustCompile("\\$([0-9]+)")

func NewRegexpTransform(regexpStr, template string, escape Escape, extraArgs config.HookExtraArgs) (Transform, error) {
	pat, err := regexp.Compile(regexpStr)
	if err != nil {
		return nil, err
	}

	return func(input string) (string, error) {
		match := pat.FindAllStringSubmatch(input, -1)
		if len(match) == 0 {
			return "", NO_MATCH
		}
		fields := match[0]

		var err error

		output := fieldPat.ReplaceAllStringFunc(template, func(f string) string {
			i, _ := strconv.Atoi(fieldPat.FindAllStringSubmatch(f, -1)[0][1])
			if len(fields)-1 < i {
				err = BAD_GROUPS
				return ""
			}
			res, escape_err := escape(fields[i], extraArgs)

			if escape_err != nil {
				err = escape_err
			}

			return res
		})

		if err != nil {
			return "", err
		}

		return output, err
	}, nil
}
