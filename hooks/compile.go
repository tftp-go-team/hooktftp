package hooks

import (
	"fmt"
	"github.com/epeli/hooktftp/regexptransform"
	"io"
)

var NO_MATCH = regexptransform.NO_MATCH

type HookComponents struct {
	Execute func(string) (io.ReadCloser, error)
	escape  regexptransform.Escape
}

type iHookDef interface {
	GetType() string
	GetDescription() string
	GetRegexp() string
	GetTemplate() string
}

type Hook func(string) (io.ReadCloser, error)

var hookMap = map[string]HookComponents{
	"file":  FileHook,
	"http":  HTTPHook,
	"shell": ShellHook,
}

func CompileHook(hookDef iHookDef) (Hook, error) {
	var ok bool
	var components HookComponents

	if hookDef.GetRegexp() == "" {
		return nil, fmt.Errorf("Cannot find regexp from hook %v", hookDef)
	}

	if components, ok = hookMap[hookDef.GetType()]; !ok {
		return nil, fmt.Errorf("Cannot find template from hook %v", hookDef)
	}

	transform, err := regexptransform.NewRegexpTransform(
		hookDef.GetRegexp(),
		hookDef.GetTemplate(),
		components.escape,
	)
	if err != nil {
		return nil, err
	}

	return func(path string) (io.ReadCloser, error) {
		newPath, err := transform(path)
		if err != nil {
			return nil, err
		}

		fmt.Println("Executing hook:", hookDef, path, "->", newPath)
		reader, err := components.Execute(newPath)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}, nil
}
