package org

import "github.com/lcyvin/gorgeous/internal/util"

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
      continue
    }

    out = append(out, v.GetEndNodes()...)
  }

  if len(out) == 0 {
    out = append(out, mnt)
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

// Inserts the subtree immediately after the given node. This causes downstream
// tree elements to be potentially moved into the inserted tree, based on the 
// levels of the subsequent node(s). This returns a /new/ MetaNodeTree. 
func (mnt *MetaNodeTree) InsertSubtree(t *MetaNodeTree) *MetaNodeTree {
  retree := *mnt
  retree.Subtree = []*MetaNodeTree{}

  tree := t.Flatten()
  
  for _, st := range mnt.Subtree {
    tree = append(tree, st.Flatten()...)
  }

  subtrees := buildTreesFromList(tree)
  for _, st := range subtrees {
    if st.Level() <= retree.Level() {
      parent := mnt.WalkBackToLevel(st.Level()-1)
      st.Parent = parent
      parent.AddSubtree(st)
    }

    if st.Level() > retree.Level() {
      st.Parent = &retree
      retree.AddSubtree(st)
    }
  }

  mnt.Subtree = retree.Subtree

  return mnt
}

func buildTreesFromList(l []*MetaNodeTree) []*MetaNodeTree {
  trees := make([]*MetaNodeTree, 0)

  if len(l) == 0 {
    return trees
  }

  last := l[0]
  root := l[0]
  for _, t := range l {
    t.Node.Tree = t
    if t.Level() <= root.Level() {
      root = t
      trees = append(trees, t)
      continue
    }

    if t.Level() == last.Level() {
      t.Parent = last.Parent
      last.Parent.AddSubtree(t)
      continue
    }

    if t.Level() < last.Level() {
      p := last.WalkBackToLevel(t.Level()-1)
      if p != nil {
        t.Parent = p
        p.AddSubtree(t)
        continue
      }
      
      panic("Unable to place node")
    }
  }

  return trees
}

func (mnt *MetaNodeTree) Flatten() []*MetaNodeTree {
  mntcp := *mnt
  mntcp.Subtree = []*MetaNodeTree{}
  mntcp.Parent = nil

  out := []*MetaNodeTree{&mntcp}

  for _, st := range mnt.Subtree {
    out = append(out, st.Flatten()...)
  }

  return out
}

func (mnt *MetaNodeTree) WalkBackToLevel(targetLvl int) *MetaNodeTree {
  if mnt.Level() <= targetLvl {
    return mnt
  }

  if mnt.Parent == nil {
    return nil
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

  if mnt.Parent.Level() != 0 || mnt.Parent.Node == nil {
    upstreamTags = mnt.Parent.InheritTags(include, exclude, all)
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

  // this is just to de-dupe the list we have. Using map keys lets us emulate
  // 'set' behavior in this instance
  coalesce := make(map[string]struct{}, 0)
  for _, v := range upstreamTags {
    if _, exists := coalesce[v]; !exists {
      coalesce[v] = struct{}{}
    }
  }

  out := make([]string, 0)
  for k := range coalesce {
    out = append(out, k)
  }

   return out
}

func (mnt *MetaNodeTree) GetNodeTags() []string {
  if mnt.Node.Heading == nil {
    return make([]string, 0)
  }

  return mnt.Node.Heading.Tags
}

// returns copies of the nodes held by each parent up the meta node tree.
func (mnt *MetaNodeTree) GetParentNodes() []*Node {
  if mnt.Level() == 0 || mnt.Parent == nil || mnt.Parent == mnt {
    return []*Node{}
  }

  parent := *mnt.Parent.Node
  return append(mnt.Parent.GetParentNodes(), &parent)
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
