package netlifyblob

import (
	storeerrors "github.com/Ankumeah/JSBEE/backend/internal/objectstore/errors"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

type stderrJSON struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func scriptErr(op string, args []string, stderr []byte, err error) error {
	label := op
	if len(args) > 0 {
		label += " " + strings.Join(args, " ")
	}
	var ej stderrJSON
	if json.Unmarshal(stderr, &ej) == nil && ej.Code != "" {
		msg := ej.Message
		if msg == "" {
			msg = ej.Code
		}
		if ej.Code == "not_found" {
			return fmt.Errorf(
				"netlifyblob: %s: %w: %s", label, storeerrors.ErrNotFound, msg,
			)
		}
		return fmt.Errorf("netlifyblob: %s failed: %s", label, msg)
	}
	trimmed := strings.TrimSpace(string(stderr))
	if trimmed == "" {
		return fmt.Errorf("netlifyblob: %s failed: %v", label, err)
	}
	return fmt.Errorf("netlifyblob: %s failed: %s", label, trimmed)
}

func (s *Store) run(
	ctx context.Context,
	op string,
	args []string,
	stdin io.Reader,
) ([]byte, error) {
	argv := make([]string, 0, len(args)+2)
	argv = append(argv, filepath.Join(s.scriptsDir, helperScript), op)
	argv = append(argv, args...)
	cmd := exec.CommandContext(ctx, s.bin, argv...)
	cmd.Env = s.env
	if stdin != nil {
		cmd.Stdin = stdin
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &limitedWriter{W: &stderr, N: 4 << 10}

	if err := cmd.Run(); err != nil {
		return nil, scriptErr(op, args, stderr.Bytes(), err)
	}
	return stdout.Bytes(), nil
}

type limitedWriter struct {
	W *bytes.Buffer
	N int
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	remaining := l.N - l.W.Len()
	if remaining <= 0 {
		return len(p), nil
	}
	if len(p) > remaining {
		p = p[:remaining]
	}
	return l.W.Write(p)
}
