package org

import "github.com/lcyvin/gorgeous/internal/util"

type Document struct {
  NodeTree *MetaNodeTree

  // BufferSettings define "client" behaviors when parsing and interpretting the
  // data structures and heritability of properties, tags, etc. within an org
  // document.
  BufferSettings *BufferSettings
}

// Instantiate a new blank document with base defaults as needed to handle
// compatibility with orgmode behaviors. 
func New() *Document {
  d := &Document{
    NodeTree: NewMetaNodeTree(),
  }

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

type MetaNodeTree struct {
  Parent *MetaNodeTree
  Node *Node
  Subtree []*MetaNodeTree
}

func NewMetaNodeTree() *MetaNodeTree {
  nmt := &MetaNodeTree{
    Parent: nil,
    Node: &Node{},
    Subtree: make([]*MetaNodeTree, 0),
  }
  nmt.Node.Tree = nmt

  return nmt
}


func (mnt *MetaNodeTree) GetEndNodes() []*MetaNodeTree {
  out := make([]*MetaNodeTree, 0)
  for _, v := range mnt.Subtree {
    if len(v.Subtree) == 0 {
      out = append(out, v)
    }
  }

  return out
}

func (mnt *MetaNodeTree) AddNode(n *Node) *MetaNodeTree {
  newMetaNode := &MetaNodeTree{
    Node: n,
    Parent: mnt,
  }

  mnt.Subtree = append(mnt.Subtree, newMetaNode)

  return mnt
}

func (mnt *MetaNodeTree) AddSubtree(st *MetaNodeTree) *MetaNodeTree {
  mnt.Subtree = append(mnt.Subtree, st)
  return mnt
}

func (mnt *MetaNodeTree) Level() int {
  return mnt.Node.Level()
}

func (mnt *MetaNodeTree) WalkBackToLevel(targetLvl int) *MetaNodeTree {
  if mnt.Level() <= targetLvl {
    return mnt
  }

  tree := mnt.Parent

  for lvl := tree.Level(); lvl <= targetLvl; tree = tree.Parent {
    if lvl <= targetLvl {
      return tree
    }
  }

  return nil
}

func (mnt *MetaNodeTree) InheritTags(include, exclude []string, all bool) []string {
  upstreamTags := make([]string, 0)

  if len(include) == 0 && !all {
    return upstreamTags
  }

  for _, v := range mnt.Parent.GetNodeTags() {
    pass := false
    if !util.In(v, include) && !all {
      pass = true
    }

    if util.In(v, exclude) && all {
      pass = true
    }

    if pass {
      continue
    }

    upstreamTags = append(upstreamTags, v)
  }

  if mnt.Parent.Level() != 0 {
    upstreamTags = mnt.Parent.InheritTags(include, exclude, all)
  }

  return upstreamTags
}

func (mnt *MetaNodeTree) GetNodeTags() []string {
  if mnt.Node.Heading == nil {
    return make([]string, 0)
  }

  return mnt.Node.Heading.Tags
}

func (mnt *MetaNodeTree) GetNodesByProperties(propMap map[string][]string) []*Node {
  nodes := make([]*Node, 0)
  this := mnt.Node
  
  for _, prop := range this.Properties {
    if v, ok := propMap[prop.Key]; ok {
      if util.In(prop.Value, v) {
        nodes = append(nodes, this)
      }
    }
  }

  for _, v := range mnt.Subtree {
    subtreeNodes := v.GetNodesByProperties(propMap)
    if len(subtreeNodes) > 0 {
      nodes = append(nodes, subtreeNodes...)
    }
  }

  return nodes
}

type TagInheritOpts struct {
  InheritAll bool
  Inherit []string
  NoInherit []string
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

  if lvl == 1 {
    d.NodeTree.AddNode(n)
    return d, nil
  }

  endNodes := d.NodeTree.GetEndNodes()
  lastNode := endNodes[len(endNodes)-1]

  if lastNode.Level() < lvl {
    lastNode.AddNode(n)
  }

  if lastNode.Level() >= lvl {
    parent := lastNode.WalkBackToLevel(lvl-1)
    parent.AddNode(n)
  }

  return d, nil
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
