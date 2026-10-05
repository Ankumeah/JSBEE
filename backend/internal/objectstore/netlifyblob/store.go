package netlifyblob

import (
	storeerrors "github.com/Ankumeah/JSBEE/backend/internal/objectstore/errors"

	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"uuid"
)

const privatePrefix = "private/"
const publicPrefix  = "public/"
const backupPrefix  = "backup/"
const backupBase    = "jsbee.sql.bak."

const bunBinary = "bun"

// Store implements the object store over
// Netlify Blobs by shelling out to a bun run helper
type Store struct {
	cfg        Config
	bin        string
	scriptsDir string
}


func resolveScriptsDir() (string, error) {
	const script = "blob.mjs"
	relForms := []string{
		"backend/internal/objectstore/netlifyblob/scripts",
		"internal/objectstore/netlifyblob/scripts",
	}
	var candidates []string
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for {
			for _, rel := range relForms {
				candidates = append(candidates, filepath.Join(dir, rel))
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(dir, "backend/internal/objectstore/netlifyblob/scripts"),
			filepath.Join(dir, "internal/objectstore/netlifyblob/scripts"),
			filepath.Join(dir, "netlifyblob-scripts"),
			filepath.Join(dir, "scripts"),
		)
	}
	var tried []string
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		tried = append(tried, abs)
		if info, err := os.Stat(filepath.Join(abs, script)); err == nil && !info.IsDir() {
			return abs, nil
		}
	}
	return "", fmt.Errorf(
		"netlifyblob: helper scripts not found (tried %s)", strings.Join(tried, ", "),
	)
}

func New(ctx context.Context, cfg Config) (*Store, error) {
	bin, err := exec.LookPath(bunBinary)
	if err != nil {
		return nil, fmt.Errorf(
			"netlifyblob: %s not on PATH: %w", bunBinary, err,
		)
	}
	scriptsDir, err := resolveScriptsDir()
	if err != nil {
		return nil, err
	}
	return &Store{cfg: cfg, bin: bin, scriptsDir: scriptsDir}, nil
}

func newWithScriptsDir(
	ctx context.Context,
	cfg Config,
	bin, scriptsDir string,
) (*Store, error) {
	resolved, err := exec.LookPath(bin)
	if err != nil {
		return nil, fmt.Errorf(
			"netlifyblob: runtime %q not found: %w", bin, err,
		)
	}
	if info, err := os.Stat(filepath.Join(scriptsDir, "blob.mjs")); err != nil || info.IsDir() {
		return nil, fmt.Errorf("netlifyblob: helper not found in %s", scriptsDir)
	}
	return &Store{cfg: cfg, bin: resolved, scriptsDir: scriptsDir}, nil
}

func (s *Store) withTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, s.cfg.Timeout)
}

func (s *Store) Init(ctx context.Context) error {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	_, err := s.run(ctx, "ping", nil, nil)
	return err
}

func (s *Store) AddFile(
	ctx context.Context,
	filename string,
	content io.Reader,
	size int64,
) error {
	if content == nil {
		return fmt.Errorf("netlifyblob: nil content for %q", filename)
	}
	if size < 0 {
		return fmt.Errorf("netlifyblob: negative size for %q", filename)
	}
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	_, err := s.run(ctx, "set", []string{privatePrefix + filename}, content)
	return err
}

func (s *Store) GetFile(
	ctx context.Context,
	filename string,
) (io.ReadCloser, error) {
	key := publicPrefix + filename

	mctx, mcancel := s.withTimeout(ctx)
	_, err := s.meta(mctx, key)
	mcancel()
	if err != nil {
		return nil, err
	}

	sctx, scancel := s.withTimeout(ctx)
	return s.stream(sctx, scancel, "get", []string{key})
}

func (s *Store) PublicFile(
	ctx context.Context,
	filename string,
) error {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	_, err := s.run(ctx, "move",
		[]string{privatePrefix + filename, publicPrefix + filename}, nil)
	return err
}

func (s *Store) PublicBaseURL() string {
	return ""
}

func (s *Store) DeleteFile(
	ctx context.Context,
	filename string,
) error {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	for _, key := range []string{publicPrefix + filename, privatePrefix + filename} {
		if _, err := s.run(ctx, "del", []string{key}, nil); err != nil {
			if errors.Is(err, storeerrors.ErrNotFound) {
				continue
			}
			return err
		}
	}
	return nil
}

func (s *Store) PresignedUploadURL(
	ctx context.Context,
	filename string,
	expiry time.Duration,
) (string, error) {
	return "", storeerrors.ErrDirectUpload
}

type metaJSON struct {
	Size int64  `json:"size"`
	ETag string `json:"etag"`
}

func (s *Store) meta(ctx context.Context, key string) (metaJSON, error) {
	var m metaJSON
	out, err := s.run(ctx, "meta", []string{key}, nil)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(out, &m); err != nil {
		return m, fmt.Errorf("netlifyblob: bad meta response: %w", err)
	}
	if m.Size < 0 {
		return m, fmt.Errorf("netlifyblob: bad meta response")
	}
	return m, nil
}

func (s *Store) PrivateFileSize(
	ctx context.Context,
	filename string,
) (int64, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	m, err := s.meta(ctx, privatePrefix+filename)
	if err != nil {
		return 0, err
	}
	return m.Size, nil
}

func (s *Store) StoreDBBackup(
	ctx context.Context,
	backupPath string,
) error {
	if s.cfg.MaxDBSnapshots < 1 {
		return nil
	}
	f, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("netlifyblob: cannot open backup: %w", err)
	}
	defer f.Close()

	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	key := backupPrefix + backupBase +
		fmt.Sprintf("%020d", time.Now().Unix()) + "." + uuid.New().String()
	if _, err := s.run(ctx, "set", []string{key}, f); err != nil {
		return err
	}

	out, err := s.run(ctx, "list", []string{backupPrefix + backupBase}, nil)
	if err != nil {
		return err
	}
	var listed struct {
		Keys []string `json:"keys"`
	}
	if err := json.Unmarshal(out, &listed); err != nil {
		return fmt.Errorf("netlifyblob: bad list response: %w", err)
	}
	sort.Strings(listed.Keys)
	for len(listed.Keys) > int(s.cfg.MaxDBSnapshots) {
		oldest := listed.Keys[0]
		listed.Keys = listed.Keys[1:]
		if _, err := s.run(ctx, "del", []string{oldest}, nil); err != nil {
			if errors.Is(err, storeerrors.ErrNotFound) {
				continue
			}
			return err
		}
	}
	return nil
}

func (s *Store) stream(
	ctx context.Context,
	cancel context.CancelFunc,
	op string,
	args []string,
) (io.ReadCloser, error) {
	argv := append(
		[]string{filepath.Join(s.scriptsDir, "blob.mjs"), op}, args...,
	)
	cmd := exec.CommandContext(ctx, s.bin, argv...)
	cmd.Env = helperEnv(os.Environ(), s.cfg)
	cmd.Stdin = nil
	var stderr bytes.Buffer
	cmd.Stderr = &limitedWriter{W: &stderr, N: 4 << 10}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("netlifyblob: %s failed: %w", op, err)
	}
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, scriptErr(op, args, nil, err)
	}
	return &streamBody{
		op:     op,
		args:   args,
		pipe:   stdout,
		cmd:    cmd,
		stderr: &stderr,
		cancel: cancel,
	}, nil
}

type streamBody struct {
	op     string
	args   []string
	pipe   io.ReadCloser
	cmd    *exec.Cmd
	stderr *bytes.Buffer
	cancel context.CancelFunc
}

func (b *streamBody) Read(p []byte) (int, error) {
	return b.pipe.Read(p)
}

func (b *streamBody) Close() error {
	readErr := b.pipe.Close()
	b.cancel()
	waitErr := b.cmd.Wait()
	if waitErr != nil {
		return scriptErr(b.op, b.args, b.stderr.Bytes(), waitErr)
	}
	if !isClosedPipeError(readErr) {
		return readErr
	}
	return nil
}

func isClosedPipeError(err error) bool {
	return err == nil || strings.Contains(err.Error(), "closed pipe")
}
