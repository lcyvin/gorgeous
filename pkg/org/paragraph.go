package org

import (
  "strings"
)

type Paragraph struct {
  Lines []string
  Raw string
}

func (p *Paragraph) String() string {
  return strings.Join(p.Lines, "\n")
}

func (p *Paragraph) Strings() []string {
  return p.Lines
}

func (p Paragraph) Kind() ElementKind {
  return ELEMENT_PARAGRAPH
}

func (p Paragraph) IsGreaterElement() bool {
  return false
}
