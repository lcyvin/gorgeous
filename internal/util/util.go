package util

import "reflect"

func In(val any, search any) bool {
  v := reflect.ValueOf(search)
  svals := make([]any, v.Len())
  for i := 0; i < v.Len(); i++ {
    svals[i] = v.Index(i).Interface()
  }

  for _, sv := range svals {
    if reflect.DeepEqual(val, sv) {
      return true
    }
  }

  return false
}
