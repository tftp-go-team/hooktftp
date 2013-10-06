package hooks

import (
	"github.com/epeli/dyntftp/config"
	"github.com/epeli/dyntftp/regexptransform"
	"io"
	"fmt"
)

var NO_MATCH = regexptransform.NO_MATCH

type HookComponents struct {
	Execute func(string) (io.Reader, error)
	escape regexptransform.Escape
}


type Hook func(string) (io.Reader, error)

func CompileHook(hookDef *config.HookDef) (Hook, error) {
	var template string
	var components HookComponents

	if hookDef.Regexp == "" {
		return nil, fmt.Errorf("Cannot find regexp from hook %v", hookDef)
	}

	if hookDef.FileTemplate != "" {
		template = hookDef.FileTemplate
		components = FileHook
	} else {
		return nil, fmt.Errorf("Cannot find template from hook %v", hookDef)
	}

	transform, err := regexptransform.NewRegexpTransform(
		hookDef.Regexp,
		template,
		components.escape,
	)
	if err != nil {
		return nil, err
	}

	return func(path string) (io.Reader, error) {
		newPath, err := transform(path)
		if err != nil {
			return nil, err
		}
		reader, err := components.Execute(newPath)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}, nil
}
