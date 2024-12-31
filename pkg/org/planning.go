package org

import "time"

// Planning elements directly follow headlines with no preceeding newlines
// and are marked by a keyword of "SCHEDULED", "DEADLINE", or "CLOSED", a
// colon, and a timestamp object.
type Planning struct {
  Kind PlanningKind
  Timestamp *TimestampRangeOrSexp
}

type PlanningKind string

const (
  // if we do not see a keyword before a timestamp, it is a simple event
  PLANNING_EVENT PlanningKind = ""
  // We expect to see these keywords depending on which keyword preceedes
  // the timestamp
  PLANNING_SCHEDULED = "SCHEDULED"
  PLANNING_DEADLINE = "DEADLINE"
)

type TimestampRangeOrSexp interface {
  // Should return out of either TIMESTAMP_KIND_TIMESTAMP or TIMESTAMP_KIND_SEXP
  Kind() TimestampKind

  // Returns true if the planning event held by the TimestampRange or sexp 
  // definition should be visible in an agenda view for the given window. 
  // Futher filtering, implementors should operate on the timestamps which
  // are valid for the window. For instance, an event with a delay that would
  // negate its display in the window is still considered valid at this level.
  //
  // Intended use pattern is to query for window tenancy, then handle further
  // filtering and data handling after type assertion based on the TimestampKind.
  InWindow(start, end time.Time) bool

  // Returns the time of day that the timestamp defines as the start, if a start
  // time is present. Else, returns 0, 0, 0
  Time() (int, int, int)
  
  // Returns the time of day that the timestamp defines as the end of the time
  // range, if a time range is set. Else, returns 0, 0, 0.
  EndTime() (int, int, int)
}

type TimestampKind string

const (
  // sentinel
  TIMESTAMP_KIND_UNKNOWN TimestampKind = ""
  TIMESTAMP_KIND_TIMESTAMP = "timestamp"
  TIMESTAMP_KIND_SEXP = "sexp"
)

type RepeatStamp interface {
  // Return the starting time of the timestamp object having a repeater, if one
  // is present.
  Start() time.Time
  // Return the ending time of the timestamp object having a repeater, if one is
  // present.
  End() time.Time
  // For implementors, return a custom Kind for the RepeatStamp for type assertion
  Kind() interface{}
  // string representation of the repeat directive (aka cookie) as would appear
  // within an org document (E.G., +1m)
  Cookie() string
  // Return true if this refers to an active timestamp, in orgmode terms
  Active() bool
  // Return true if the RepeatStamp would occur within the passed window, based
  // on the repeat cookie and implemented repetition logic.
  InWindow(start, end time.Time) bool
}

type Repeater interface {
  // Perform a single shift on a RepeatStamp based on the cookie it holds.
  Shift() RepeatStamp
  // Perform <n> shifts on a RepeatStamp based on the cookie it holds.
  Shiftn(i int) RepeatStamp
  // Perform as many shifts on a RepeatStamp as needed until it is less than or
  // equal to the time passed to in t
  ShiftUntil(t time.Time) RepeatStamp
  // Perform as many shifts on a RepeatStamp as needed to be in the future of
  // the time passed to t
  ShiftUntilAfter(t time.Time) RepeatStamp
}
