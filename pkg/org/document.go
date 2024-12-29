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

// Instantiate a new blank document with base defaults as needed to handle
// compatibility with orgmode behaviors. 
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

// HeadingOpt funcs provide handlers to control the instantiation of a new
// heading element as created by the Document.AddHeading method.
type HeadingOpt func(*Heading)

// Sets the priority (E.G., [#A]) of the newly created heading. The valid
// values for a heading's priority are defined in PrioritySettings. By default,
// the valid values are of type AlphaHeadingPriority with a string value of
// "A", "B", or "C". If no priority is set, it is implicitly seen by org as
// having the default priority (B, if no custom priority values are set).
func WithPriority(p HeadingPriority) HeadingOpt {
  return func(h *Heading) {
    h.Priority = p
  }
}

// Sets the collection of tags held by the headline
func WithTags(tags []string) HeadingOpt {
  return func(h *Heading) {
    h.Tags = tags
  }
}

// If the heading contains the token "COMMENT" (case sensitive) immediately
// preceeding the title token (I.E., the heading text before the tags), the
// entire node and its children are considered to be "commented out" and will
// be omitted from exports, queries, etc..
func WithHeadingIsComment() HeadingOpt {
  return func(h *Heading) {
    h.IsComment = true
  }
}

// Adds a new node and heading to end of the document's node list, setting
// inheritance and relatives based on the heading's level and the previous
// node's own heading level. If none save for the zero-th node are of a higher
// significance (lower level/fewer asterisks), its parent will be the zero-th
// node. AddHeading only encapsulates the information present on a single
// headline definition, and thus only creates a bare node without properties
// or section contents.
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

  n := &Node{
    Heading: h,
    Document: d,
  }

  doc, _, err := d.InsertMetaNode(n)
  if err != nil {
    return nil, err
  }

  return doc, nil
}

// Inserts a MetaNode entry to the end of the Document.Nodes list, returning
// a pointer to the modified document, and the meta node itself, and may
// additionally return an error.
func (d *Document) InsertMetaNode(n *Node) (*Document, *MetaNode, error) {
  if n == nil {
    return nil, nil, &NilMetaNodeError{}
  }

  if n.Heading == nil {
    return nil, nil, &NilNodeHeadingError{}
  }

  // set n.Index to -1 in order to signal we have not processed it, and to
  // prevent default assignment related issues.
  n.Index = -1

  mn := &MetaNode{Node: n}

  // determine MetaNode's parent based on heading level if it is not "1"
  for i := len(d.Nodes)-1; i >= 0; i-- {
    // we handle the default case of a new node with a heading of lvl 1 being
    // owned by the zero-th node in the first loop instantiation to avoid heavy
    // nesting
    if n.Heading.Level == 1 {
      mn.ParentIdx = 0
      n.Index = len(d.Nodes)
      d.Nodes = append(d.Nodes, mn)
      break
    }

    if nodeHdg := d.Nodes[i].Node.Heading; nodeHdg != nil {
      if nodeHdg.Level < n.Heading.Level {
        mn.ParentIdx = i
        n.Index = len(d.Nodes)
        d.Nodes = append(d.Nodes, mn)
        break
      }
    }
  }

  if n.Index < 0 {
    return nil, nil, UnknownInsertError{}
  }

  return d, mn, nil
}

type NilMetaNodeError struct {}
func (NilMetaNodeError) Error() string {
  return "Unable to insert nil meta node into document node tree"
}

type NilNodeHeadingError struct{}
func (NilNodeHeadingError) Error() string {
  return "Node heading must not be nil"
}

type UnknownInsertError struct {}
func (UnknownInsertError) Error() string {
  return "Unable to insert node, unexpected behavior encountered"
}
