package hooks

import (
	"github.com/epeli/hooktftp/regexptransform"
	"io"
	"os"
)

type FileHookCore struct {
	transform regexptransform.Transform
}

func (h *FileHookCore) Transform(s string) (string, error) {
	return h.transform(s)
}

func (h *FileHookCore) NewReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func FileEscape(s string) string {
	return s
}

func NewFileHookCore(hookDef iHookDef) (HookCore, error) {
	var template string = hookDef.GetFileTemplate()
	var regexp string = hookDef.GetRegexp()
	if template == "" || regexp == "" {
		return nil, INVALID_HOOK
	}

	tranform, err := regexptransform.NewRegexpTransform(
		hookDef.GetRegexp(),
		template,
		FileEscape,
	)
	if err != nil {
		return nil, err
	}

	h := &FileHookCore{tranform}

	return h, nil
}
