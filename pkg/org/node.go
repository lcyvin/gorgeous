package org

type Node struct {
  Index       int
  Heading     *Heading
  Section     *Section
  Properties  []Property
  Document    *Document
}
