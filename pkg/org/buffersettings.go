package org

import (
	"fmt"
	"strings"
)

// BufferSettings define various metadata and client behaviors, largely to
// handle how certain special keywords are handled or to override default
// values for parts of an element, primarily headings.
type BufferSettings struct {
  // Tags set by the FILETAGS property, to be inherited by all
  // headilnes within the document.
  FileTags      []string
  // Properties set within a property drawer at the top of a file
  // or with the #+PROPERTY keyword can be inherited by all nodes
  // within the document tree if org-use-property-inheritance is
  // set to "true", if the property key is within the list held by
  // the variable, or if the key matches the regex set by the variable.
  // Properties defined here whose keys contain the suffix "_All"
  // are always inherited, and apply restrictions to the allowed values
  // for the corresponding property sans-suffix.
  Properties    []*Property
  // The archive location (file name and optional mark) for any tree or
  // subtree within this document. This property is inherited by all nodes
  // within the document tree.
  Archive       string
  // Sets the category of the file, and is applied to all nodes within the
  // document. Used for agenda mode sorting and filtering.
  Category      string
  // Defines the formatting for column view
  Columns       string
  // Defines constants that table forumlas can make use of
  Constants     map[string]string
  // Defines abbreviations for links, allowing for shorthand referencces
  // to otherwise unweildy URLs. 
  // 
  // The map's key is the abbreviation, in the format of "org-string",
  // containing non-whitespace characters unless braced by quotation
  // marks. The value held by the key is the link to be used. If the
  // link contains the formatting mark '%s', the resulting link at the
  // point of usage should interpolate the string with the value of a 
  // tag accompanying the link abbreviation's invocation. E.G.:
  //    #+LINK: myabbreviation https://myurl.tld/%s
  //    ...
  //    [[myabbreviation:foo][see foo]]
  //    # results in linking to https://myurl.tld/foo
  Links         map[string]string
  // Sets the values to be used for heading priority levels. Can be either
  // alpha characters (a, b ,c), or numbers less than 65. The format for
  // this setting is:
  //
  //     #+PRIORITIES: A C B
  //
  // where A is the highest priority, C is the lowest, and B is the default
  // set when no priority is explicitly set.
  // 
  // You may set a broader range of priorities, for instance:
  //
  //     #+PRIORITIES: A E C
  // or
  //
  //     #+PRIORITIES: 1 10 5
  Priorities    *HeadingPrioritySetting
  // SetupFile contains additional buffer settings to be used in this file.
  // See BufferSettings.AddSetupFile for adding a setupfile to an existing
  // document. When parsing, this should be called if a setupfile setting is
  // encountered.
  SetupFile     *SetupFile
  // Todo keywords can be defined as a sequence of either states, represented
  // by all-caps strings containing only alphabet characters, or for backwards
  // compatibility as types, represented by strings of only alphabet characters
  // and starting with an uppercase Character. 
  // In a definition, the sequence "|" occurring in the list of states or types
  // indicates that all following items represent "done" states. If none is
  // present, the final element in the list is considered the "done" state.
  //
  // Multiple todo sequences may be defined in a file, for instance, for
  // different workflows. The same key should not be defined in multiple places
  // in order to allow for proper sequence cycling. E.G.:
  //
  //     #+TODO: BACKLOG TODO BLOCKED STARTED REVIEWING STAGED | DONE CANCELLED
  //     #+TODO: IDEA SPIKE RFC | PLANNED DENIED
  //     #+TODO: WAITING MEETING CALL EMAIL | HANDLED NOOP
  //
  // It is recommended to use tags in favor of types where relevant.
  TodoSettings  *TodoSettings
}

type HeadingPriority interface {
  String()  string
  Kind()    HeadingPriorityKind
  Higher()  bool
  Equal()   bool
}

type HeadingPrioritySetting struct {
  Kind      HeadingPriorityKind
  Highest   HeadingPriority
  Lowest    HeadingPriority
  Default   HeadingPriority
}

type HeadingPriorityKind int

const (
  HEADING_PRIORITY_INT HeadingPriorityKind = iota
  HEADING_PRIORITY_ALPHA
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
