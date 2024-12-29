package org

type Document struct {
  // A list containing each node in order of parsing where the int value
  // refers to the list index of the node's parent. For a first-level Node,
  // this will always be `0`. For Any node at a lower level than 1, it will
  // be the most recent node occurring at a higher level. Id est, in a Document
  // wherein the first heading occurs at a level of `2`, the most recent
  // higher-level node would still be the zero-th node.
  Nodes          []*MetaNode
  // BufferSettings define "client" behaviors when parsing and interpretting the
  // data structures and heritability of properties, tags, etc. within an org
  // document.
  BufferSettings *BufferSettings
}

// Instantiate a new blank document with no contents
func New() *Document {
  d := &Document{}
  d.Nodes = []*MetaNode{
    &MetaNode{
      ParentIdx: 0,
      Node: &Node{
        Index: 0,
        Document: d,
  }}}

  bufSettings := &BufferSettings{}
  todoSettings := &TodoSettings{}
  todoSettings.Add(&TodoSequence{
    ProcessKeywords: []string{"TODO"},
    DoneKeywords: []string{"DONE"},
    Kind: TODO_SEQUENCE_STATE,
  })

  prioritySettings := &HeadingPrioritySetting{
    Kind: HEADING_PRIORITY_ALPHA,
    Highest: AlphaHeadingPriority("A"),
    Lowest: AlphaHeadingPriority("C"),
    Default: AlphaHeadingPriority("B"),
  }

  bufSettings.Priorities = prioritySettings
  bufSettings.TodoSettings = todoSettings
  d.BufferSettings = bufSettings

  return d
}

// MetaNode describes the given node's position within the document as well as
// its position within the overall node tree. Individual nodes do not hold
// direct references to their child nodes in favor of this structure, in order
// to allow for top-down processing of operations such as refiling and
// inheritance, as in org's own api implementation.
type MetaNode struct {
  // Refers to the index within Document.Nodes of the node which is above the
// given node in the parent tree. If this is the zero-th node, it refers to
  // itself.
  ParentIdx int
  // Points to the actual node itself in order to allow for operating on the
  // elements contained within said node
  Node      *Node
}

type HeadingOpt func(*Heading)

func WithPriority(p HeadingPriority) HeadingOpt {
  return func(h *Heading) {
    h.Priority = p
  }
}

func WithTags(tags []string) HeadingOpt {
  return func(h *Heading) {
    h.Tags = tags
  }
}

func WithHeadingIsComment() HeadingOpt {
  return func(h *Heading) {
    h.IsComment = true
  }
}

func (d *Document) AddHeading(lvl int, text string, opts... HeadingOpt) (*Document, error) {
  if lvl < 0 {
    return nil, NewInvalidHeadingLevelError(lvl)
  }
  h := &Heading{
    Level: lvl,
    Text: text,
  }

  for _, opt := range opts {
    opt(h)
  }

  return d, nil
}
