package api

type Planning interface {
  Scheduled() Timing
  Deadline()  Timing
  Event()     Timing
}
