package org

import "time"

// Planning elements directly follow headlines with no preceeding newlines
// and are marked by a keyword of "SCHEDULED", "DEADLINE", or "CLOSED", a
// colon, and a timestamp object.
type Planning struct {
  Kind PlanningKind
  TimestampRangeOrSexp TimestampRangeOrSexp
}

type PlanningKind string

const (
  // if we do not see a keyword before a timestamp, it is a simple event
  PLANNING_EVENT PlanningKind = ""
  // We expect to see these keywords depending on which keyword preceedes
  // the timestamp
  PLANNING_SCHEDULED PlanningKind = "SCHEDULED"
  PLANNING_DEADLINE  PlanningKind = "DEADLINE"
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

  // Returns an array of timestamp elements 
}

type TimestampKind string

const (
  // sentinel
  TIMESTAMP_KIND_UNKNOWN TimestampKind = ""
  TIMESTAMP_KIND_TIMESTAMP = "timestamp"
  TIMESTAMP_KIND_TIMESTAMP_RANGE = "timestamp-range"
  TIMESTAMP_KIND_SEXP = "sexp"
)
