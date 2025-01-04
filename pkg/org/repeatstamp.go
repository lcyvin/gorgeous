package org

import (
  "time"
)

var DefaultRepeatConfig RepeatConfig = RepeatConfig{
  ClampToEndOfMonth: true,
  ShiftByDays: false,
  FixedDate: true,
  FixedDateRef: nil,
}

type RepeatConfig struct {
  // ClampToEndOfMonth modifies the behavior of the Repeater to always set
  // shifted dates to be no later than the last day of the month following the
  // month referenced in the timstamp's current state. Its behavior varies 
  // depending on the values of other RepeatConfig struct items. The possible
  // behaviors are:
  //
  // - ClampToEndOfMonth only is true: date is always shifted to the last day
  //   of the target month.
  //
  // - ShiftByDays is also true: date is shifted by 30 days, unless said shift
  //   results in the new timestamp "Jumping" a month. E.G., with this setting
  //   a one month repeat on a date of January 31st would normally shift to 
  //   March 2nd or 3rd, but instead will clamp to the last day of the month,
  //   being Feb. 28th or 29th. 
  // 
  // - anch with personal remote
  ClampToEndOfMonth bool
  ShiftByDays bool
  FixedDate bool
  FixedDateRef interface{}
}

// RepeatStamp is a meta struct that implements the api.Repeater and
// api.RepeatStamp interfaces, providing handling for repeat directives set in
// a timestamp based on the behavior defined in RepeatConfig. If no
// RepeatConfig is set when a shift, window, etc. operation is called,
// DefaultRepeatConfig is used.
type RepeatStamp struct {
  Timestamp
  RepeatConfig RepeatConfig
}

func NewRepeatStampFromTimestamp(ts *Timestamp, cfg RepeatConfig) *RepeatStamp {
  return &RepeatStamp{
    Timestamp: *ts,
    RepeatConfig: DefaultRepeatConfig,
  }
}

func (rs *RepeatStamp) InWindow(start, end time.Time) bool {
  
  return false
}

