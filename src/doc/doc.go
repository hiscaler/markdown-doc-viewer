package doc

import "time"

type Doc struct {
	Title          string
	Author         string
	Sections       []Section
	ModifyDatetime time.Time
}
