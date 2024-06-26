package semver

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var (
	ErrNoTags = errors.New("no suitable tags found")
)

type Bump struct {
	version *semver.Version
	input   string
}

func New(v string) *Bump {
	return &Bump{
		input:   v,
		version: semver.MustParse(v),
	}
}

func (b *Bump) Prefixed() bool {
	return strings.HasPrefix(b.input, "v")
}

func (b *Bump) StringFull() string {
	out := b.version.String()
	if b.Prefixed() {
		out = "v" + out
	}
	return out
}

func (b *Bump) IncMajor() {
	inc := b.version.IncMajor()
	b.version = &inc
}

func (b *Bump) IncMinor() {
	inc := b.version.IncMinor()
	b.version = &inc
}

func (b *Bump) IncPatch() {
	inc := b.version.IncPatch()
	b.version = &inc
}

func Latest(tags []string) (string, error) {
	vs := make(semver.Collection, 0, len(tags))

	for _, r := range tags {
		v, err := semver.NewVersion(r)
		if err != nil {
			continue
		}

		vs = append(vs, v)
	}

	if len(vs) < 1 {
		return "", ErrNoTags
	}

	latest := vs[0]

	for _, t := range vs {
		if t.GreaterThan(latest) {
			latest = t
		}
	}

	return latest.Original(), nil
}
