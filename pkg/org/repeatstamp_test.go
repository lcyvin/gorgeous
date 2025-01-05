package org

import (
  "time"
  "testing"
)

func testingRepeatStamp(kind RepeatKind, amt int, interval RepeatIntervalKind) *RepeatStamp {
  return NewRepeatStamp(time.Date(
      2020, time.Month(1), 1,
      8, 30, 0, 0,
      time.Now().UTC().Location(),
      ),
    DefaultRepeatConfig,
    WithRepeat(&Repeat{
      Kind: kind,
      IntervalAmount: amt,
      Interval: interval,
    }))
}

func TestShiftFutureFixed(t *testing.T) {
  tk := REPEAT_KIND_SHIFT_FUTURE_FIXED
  testTime := time.Date(
    2020, time.Month(1), 1,
    12, 0, 0, 0,
    time.Now().UTC().Location(),
    )

  var tests = []struct {
    input *RepeatStamp
    want []int
  }{{
      testingRepeatStamp(tk, 6, REPEAT_INTERVAL_HOUR),
      []int{2020, 1, 1, 14, 30, 0},
    },{
      testingRepeatStamp(tk, 6, REPEAT_INTERVAL_DAY),
      []int{2020, 1, 7, 8, 30, 0},
    },{
      testingRepeatStamp(tk, 6, REPEAT_INTERVAL_WEEK),
      []int{2020, 2, 12, 8, 30, 0},
    }}

  for _, test := range tests {
    got := test.input.Shift(testTime)
    hour, minute, second := got.Time()
    testArr := []int{got.Year(), got.Month(), got.Day(), hour, minute, second}
    for i, v := range testArr {
      if v != test.want[i] {
        t.Errorf("ShiftFutureFixed(%v) = %v", test.input, testArr)
      }
    }
  }
}
