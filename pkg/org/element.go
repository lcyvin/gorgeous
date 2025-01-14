package org

type Element interface {
  Kind() ElementKind
  IsGreaterElement() bool
  String() string
  Strings() []string
}

type ElementKind int

const (
  // reserve enum 0 for bad assignment
  ELEMENT_INVALID ElementKind = iota
  // greater elements
  ELEMENT_HEADING
  ELEMENT_GREATER_BLOCK
  ELEMENT_DRAWER
  ELEMENT_DYNAMIC_BLOCK //TODO
  ELEMENT_FOOTNOTE_DEF
  ELEMENT_INLINE_TASK //TODO
  ELEMENT_ITEM
  ELEMENT_LIST
  ELEMENT_PROPERTY_DRAWER
  ELEMENT_TABLE //TODO
  // lesser elements
  ELEMENT_BLOCK
  ELEMENT_CLOCK
  ELEMENT_PLANNING
  ELEMENT_COMMENT
  ELEMENT_FIXED_WIDTH
  ELEMENT_HORIZONTAL_RULE
  ELEMENT_KEYWORD
  ELEMENT_LATEX //TODO
  ELEMENT_NODE_PROPERTY
  ELEMENT_PARAGRAPH
  ELEMENT_TABLE_ROW
  ELEMENT_DIARY_SEXP //TODO
)

// Legible strings for error and debug output purposes
func (ek ElementKind) String() string {
  elemStringMap := map[ElementKind]string{
    ELEMENT_HEADING: "Heading",
    ELEMENT_GREATER_BLOCK: "Greater Block",
    ELEMENT_DRAWER: "Drawer",
    ELEMENT_DYNAMIC_BLOCK: "Dynamic Block",
    ELEMENT_FOOTNOTE_DEF: "Footnote Definition",
    ELEMENT_INLINE_TASK: "Inline Task",
    ELEMENT_ITEM: "Item",
    ELEMENT_LIST: "List",
    ELEMENT_PROPERTY_DRAWER: "Property Drawer",
    ELEMENT_TABLE: "Table",
    ELEMENT_BLOCK: "Block",
    ELEMENT_CLOCK: "Clock",
    ELEMENT_PLANNING: "Planning",
    ELEMENT_COMMENT: "Comment",
    ELEMENT_FIXED_WIDTH: "Fixed Width Area",
    ELEMENT_HORIZONTAL_RULE: "Horizontal Rule",
    ELEMENT_KEYWORD: "Keyword",
    ELEMENT_LATEX: "Latex Environment",
    ELEMENT_NODE_PROPERTY: "Node Property",
    ELEMENT_PARAGRAPH: "Paragraph",
    ELEMENT_TABLE_ROW: "Table Row",
  }

  o, ok := elemStringMap[ek]
  if !ok {
    o = "Invalid"
  }

  return o
}

// Standard convention, this may not get used
func (ek ElementKind) EnumIndex() int {
  return int(ek)
}
