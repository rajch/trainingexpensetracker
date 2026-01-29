package model

import "time"

type Expense struct {
	Id          int
	Date        time.Time
	Description string
	Amount      float64
}
