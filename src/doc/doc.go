package doc

import "time"

type Doc struct {
	Name           string
	Title          string
	Author         string
	Sections       map[string]Section
	ModifyDatetime time.Time
}
