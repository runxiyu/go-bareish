package example

import (
	"fmt"
	"time"

	bare "git.sr.ht/~runxiyu/go-bareish"
)

type Time time.Time

func (t *Time) Unmarshal(r *bare.Reader) error {
	st, err := r.ReadString()
	if err != nil {
		return fmt.Errorf("Time.Unmarshal: read string: %e", err)
	}

	tm, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return fmt.Errorf("Time.Unmarshal: parse time: %e", err)
	}

	*t = Time(tm)
	return nil
}
