package hooks

import (
	"fmt"
	"github.com/epeli/hooktftp/regexptransform"
	"io"
)

var NO_MATCH = regexptransform.NO_MATCH

type HookComponents struct {
	Execute func(string) (io.Reader, error)
	escape  regexptransform.Escape
}

type iHookDef interface {
	GetName() string
	GetRegexp() string
	GetShellTemplate() string
	GetFileTemplate() string
}

type Hook func(string) (io.Reader, error)

func CompileHook(hookDef iHookDef) (Hook, error) {
	var template string
	var components HookComponents

	if hookDef.GetRegexp() == "" {
		return nil, fmt.Errorf("Cannot find regexp from hook %v", hookDef)
	}

	if t := hookDef.GetFileTemplate(); t != "" {
		template = t
		components = FileHook
	} else if t := hookDef.GetShellTemplate(); t != "" {
		template = t
		components = ShellHook
	} else {
		return nil, fmt.Errorf("Cannot find template from hook %v", hookDef)
	}

	transform, err := regexptransform.NewRegexpTransform(
		hookDef.GetRegexp(),
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

		fmt.Println("Executing hook", hookDef)
		reader, err := components.Execute(newPath)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}, nil
}
