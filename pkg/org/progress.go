package org

import (
	"fmt"
	"strings"
)

type ProgressCookie struct {
  Tree []Element
  Kind ProgressKind
}

func (pc *ProgressCookie) String() string {
  switch pc.Kind {
  case PROGRESS_KIND_PERCENT:
    return pc.PercentString()
  case PROGRESS_KIND_FRACTION:
    fallthrough
  default:
    return pc.FractionString()
  }
}

func (pc *ProgressCookie) FractionString() string {
  return fmt.Sprintf("[%d/%d]", pc.Done(), pc.Total())
}

func (pc *ProgressCookie) PercentString() string {
  div := float64(pc.Done())/float64(pc.Total())
  return fmt.Sprintf("[%.0f%%]", div*100)
}

func (pc *ProgressCookie) Done() int {
  // TODO

  return 0
}

func (pc *ProgressCookie) Total() int {
  //TODO

  return 0
}

// Returns a new pointer to a ProgressCookie with the `kind` set.
func ProgressCookieFromString(s string) *ProgressCookie {
  out := &ProgressCookie{}

  find := []ProgressKind{PROGRESS_KIND_PERCENT, PROGRESS_KIND_FRACTION}
  for _, kind := range find {
    if test := strings.Index(s, kind.String()); test > -1 {
      out.Kind = kind
      return out
    }
  }

  return nil
}

type ProgressKind string

const (
  PROGRESS_KIND_FRACTION ProgressKind = "/"
  PROGRESS_KIND_PERCENT ProgressKind = "%"
)

func (pk ProgressKind) String() string {
  return string(pk)
}
