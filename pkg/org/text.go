package org

type Paragraph struct {
  Raw string
}

func (p Paragraph) Kind() ElementKind {
  return ELEMENT_PARAGRAPH
}

func (p Paragraph) IsGreaterElement() bool {
  return false
}
