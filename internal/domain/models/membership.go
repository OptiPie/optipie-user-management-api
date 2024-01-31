package models

import "time"

type Membership struct {
	Type                string
	LiveMode            bool
	Attempt             int32
	Created             time.Time
	Updated             time.Time
	EventId             int64
	Id                  int64
	Amount              float64
	Object              string
	Paused              string
	Status              string
	Canceled            string
	Currency            string
	PspId               string
	MembershipLevelId   int64
	MembershipLevelName string
	StartedAt           time.Time
	CanceledAt          time.Time
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	SupporterId         int64
	SupporterEmail      string
	CurrentPeriodEnd    time.Time
	SupporterFeedback   string
	CancelAtPeriodEnd   string
	CurrentPeriodStart  time.Time
}
