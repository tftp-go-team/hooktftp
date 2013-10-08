package hooks

import (
	"errors"
	"fmt"
	"github.com/epeli/hooktftp/regexptransform"
	"io"
)

var NO_MATCH = regexptransform.NO_MATCH
var INVALID_HOOK = errors.New("invalid hook")

type HookComponents struct {
	Execute func(string) (io.Reader, error)
	escape  regexptransform.Escape
}

type iHookDef interface {
	GetRegexp() string
	GetShellTemplate() string
	GetFileTemplate() string
}

type HookCore interface {
	NewReader(string) (io.Reader, error)
	Transform(string) (string, error)
}

type Hook struct {
	core      HookCore
}

type NewHookCore func(iHookDef) (HookCore, error)

func (h *Hook) Execute(path string) (io.Reader, error) {
	transformedPath, err := h.core.Transform(path)
	if err != nil {
		return nil, err
	}
	return h.core.NewReader(transformedPath)
}

func CompileHook(hookDef iHookDef) (*Hook, error) {
	hookCoreCreators := []NewHookCore{
		NewFileHookCore,
	}

	var hook *Hook
	for _, newCore := range hookCoreCreators {
		core, err := newCore(hookDef)
		if err == INVALID_HOOK {
			continue
		}
		if err != nil {
			return hook, err
		}

		hook = &Hook{core}
	}

	if hook == nil {
		return hook, fmt.Errorf("Failed to compile hook from %v", hookDef)
	}

	return hook, nil
}
