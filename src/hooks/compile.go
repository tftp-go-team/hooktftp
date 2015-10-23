package hooks

import (
	"fmt"
	"io"

	"github.com/tftp-go-team/hooktftp/src/logger"
	"github.com/tftp-go-team/hooktftp/src/regexptransform"
	"github.com/tftp-go-team/libgotftp/src"
)

var NO_MATCH = regexptransform.NO_MATCH

type HookResult struct {
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Length int
}

func newHookResult(stdout, stderr io.ReadCloser, length int) *HookResult {
	return &HookResult{
		stdout,
		stderr,
		length,
	}
}

type Hook func(path string, request tftp.Request) (*HookResult, error)

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

	return func(path string, request tftp.Request) (*HookResult, error) {
		newPath, err := transform(path)
		if err != nil {
			return nil, err
		}

		logger.Info("Executing hook: %s %s -> %s", hookDef, path, newPath)
		return components.Execute(newPath, request)
	}, nil
}
