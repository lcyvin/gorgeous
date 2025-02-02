package org

import (
	"strconv"
	"strings"
)

// TODO: 
// - add list nesting comprehension
// - string output handling
// - cookie comprehensions
// - checkbox state tracking comprehension

var alphas string = "abcdefghijklmnopqrstuvwxyz"

type List struct {
  Ordered bool
  Suffix string
  Items []ListItem
  CounterKind CounterKind
}

func (l List) Kind() ElementKind {
  return ELEMENT_LIST
}

func (l List) IsGreaterElement() bool {
  return true
}

func (l *List) OrderedMap() map[string]ListItem {
  out := make(map[string]ListItem, 0)

  counter := 0
  for _, v := range l.Items {
    if v.Cookie != "" {
      c := v.CookieIdx(l.CounterKind)
      if c > -1 {
        counter = c
      }
    }

    if counter > 0 && v.Cookie == "" {
      counter += 1
    }

    out[l.CounterKind.StringAt(counter)] = v
  }

  return out
}

type CheckBox struct {
  State CheckBoxState
}

type CheckBoxState string

const (
  CHECKBOX_UNCHECKED CheckBoxState = " "
  CHECKBOX_PARTIAL CheckBoxState = "-"
  CHECKBOX_CHECKED CheckBoxState = "X"
)

func (cbs CheckBoxState) String() string {
  return string(cbs)
}

type ListItem struct {
  Cookie string
  Numerator int
  Elements []Element
  CheckBox *CheckBox
}

func (li ListItem) Kind() ElementKind {
  return ELEMENT_ITEM
}

func (li ListItem) IsGreaterElement() bool {
  return true
}

func (li *ListItem) String(idx int, suffix string) string {
  out := ""
  if suffix == "" && (li.Numerator > 0 || idx > 0) {
    suffix = "."
  }

  //TODO
  
  return out
}

func (li *ListItem) CookieIdx(k CounterKind) int {
  if k == COUNTER_KIND_NUM {
    i, err := strconv.ParseInt(li.Cookie, 10, 64)
    if err != nil {
      return -1
    }

    return int(i)
  }

  return strings.Index(li.Cookie, alphas)
}



type CounterKind string

const (
  COUNTER_KIND_ALPHA CounterKind = "alpha"
  COUNTER_KIND_NUM CounterKind = "number"
)

func (ck CounterKind) StringAt(i int) string {
  switch ck {
  case COUNTER_KIND_NUM:
    return strconv.Itoa(i)
  case COUNTER_KIND_ALPHA:
    return string(alphas[i]) 
  default:
    return ""
  }
}
