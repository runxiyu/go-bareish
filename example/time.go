package main

import "time"

type Time time.Time

func (t *Time) Decode(data []byte) error {
	tm, err := time.Parse("2006-01-02T15:04:05Z0700", string(data))
	if err != nil {
		return err
	}
	*t = Time(tm)
	return nil
}
