package org

import (
  "time"
)

type ClockEntry struct {
  TimeRange Timestamp
}

func (ce *ClockEntry) Duration() time.Duration {
  if ce.TimeRange.End.IsZero() {
    return time.Duration(0)
  }
  return ce.TimeRange.End.Sub(ce.TimeRange.Start)
}
