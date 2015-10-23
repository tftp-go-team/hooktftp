package hooks

import (
	"fmt"
	"io"

	"github.com/tftp-go-team/hooktftp/src/logger"
	"github.com/tftp-go-team/hooktftp/src/regexptransform"
	"github.com/tftp-go-team/libgotftp/src"
)

var NO_MATCH = regexptransform.NO_MATCH

type Hook func(string, tftp.Request) (io.ReadCloser, io.ReadCloser, int, error)

type HookComponents struct {
	Execute Hook
	escape  regexptransform.Escape
}

type iHookDef interface {
	GetType() string
	GetDescription() string
	GetRegexp() string
	GetTemplate() string
}

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

	return func(path string, request tftp.Request) (io.ReadCloser, io.ReadCloser, int, error) {
		newPath, err := transform(path)
		if err != nil {
			return nil, nil, -1, err
		}

		logger.Info("Executing hook: %s %s -> %s", hookDef, path, newPath)
		outReader, errReader, length, err := components.Execute(newPath, request)
		if err != nil {
			return nil, nil, length, err
		}
		return outReader, errReader, length, nil
	}, nil
}
