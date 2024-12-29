package org

import (
  "fmt"
  "strings"
)

// Headings represent logical separations between sections
// within a document and are components of an outline (or node)
// tree within a document. A heading occuring at a "lower" (higher
// int value) level than the preceeding heading are considered to
// be children of the preceeding heading. A child heading will always
// inherit tags from its parents as well as any set with the `FILETAGS`
// property. Other heritable items (e.g., properties) are handled at the
// node level.
type Heading struct {
  Text        string
  Priority    HeadingPriority
  IsComment   bool
  Tags        []string
  Level       int
  Node        *Node
}

func (h Heading) Kind() ElementKind {
  return ELEMENT_HEADING
}

func (h Heading) IsGreaterElement() bool {
  return true
}

// GetPriority returns the value held by Heading.Priority, or returns
// PRIORITY_DEFAULT if none is defined. 
func (h *Heading) GetPriority() HeadingPriority {
  if h.Priority != nil {
    return h.Priority
  }

  return PriorityExtrema("B")
}

type InvalidHeadingLevelError struct {
  Lvl int
}

func (ihle InvalidHeadingLevelError) Error() string {
  return fmt.Sprintf("Heading level of %d is invalid, must be greater than 0", ihle.Lvl)
}

func NewInvalidHeadingLevelError(l int) *InvalidHeadingLevelError {
  return &InvalidHeadingLevelError{Lvl: l}
}

type HeadingPriority interface {
  String()  string
  Kind()    HeadingPriorityKind
}

type HeadingPriorityKind int

const (
  HEADING_PRIORITY_INT HeadingPriorityKind = iota
  HEADING_PRIORITY_ALPHA
  HEADING_PRIORITY_EXTREMA
)

// Type for handling integer-based heading priorities
type IntHeadingPriority int

// Returns true if this heading is higher significance (lower number)
// than the provided priority. Used for sorting.
func (ihp IntHeadingPriority) Higher(p IntHeadingPriority) bool {
  return int(ihp) < int(p)
}

// Returns true if the provided priority is of the same significance
func (ihp IntHeadingPriority) Equal(p IntHeadingPriority) bool {
  return int(ihp) == int(p)
}

// Returns HEADING_PRIORITY_INT for type assertion purposes
func (ihp IntHeadingPriority) Kind() HeadingPriorityKind {
  return HEADING_PRIORITY_INT
}

// Returns a stringification of the provided value
func (ihp IntHeadingPriority) String() string {
  return fmt.Sprintf("%d", int(ihp))
}

// Type for handling alpha-based heading priorities
type AlphaHeadingPriority string

// Returns true if this heading has a higher significance (earlier
// alpha character, E.G., A being 'higher' than B) than p
func (ahp AlphaHeadingPriority) Higher(p AlphaHeadingPriority) bool {
  alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  return strings.Index(alphabet, string(ahp)) < strings.Index(alphabet, string(p))
}

// Returns true if this heading's priority is of the same significance as p
func (ahp AlphaHeadingPriority) Equal(p AlphaHeadingPriority) bool {
  return string(ahp) == string(p)
}

// Returns HEADING_PRIORITY_ALPHA for type assertion purposes
func (ahp AlphaHeadingPriority) Kind() HeadingPriorityKind {
  return HEADING_PRIORITY_ALPHA
}

// Returns the string representation of the priority's value
func (ahp AlphaHeadingPriority) String() string {
  return string(ahp)
}

// The priority extrema type is defined to provide context-unaware handling
// of priority operations, such that a heading element extant outside the
// context of a tree (for instance, generated programatically and not yet
// inserted into the tree) can still fulfill required behaviors (eg, sort).
type PriorityExtrema string

const (
  PRIORITY_HIGHEST PriorityExtrema = "A"
  PRIORITY_LOWEST                  = "C"
  PRIORITY_DEFAULT                 = "B"
)

func (pe PriorityExtrema) Kind() HeadingPriorityKind {
  return HEADING_PRIORITY_EXTREMA
}

func (pe PriorityExtrema) String() string {
  return string(pe)
}
