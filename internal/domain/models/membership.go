package models

type Membership struct {
	Type                string
	LiveMode            bool
	Attempt             int32
	Created             int64
	EventId             int64
	Id                  int64
	Amount              float64
	Object              string
	Paused              string
	Status              string
	Canceled            string
	Currency            string
	PspId               string
	MembershipLevelId   string
	MembershipLevelName string
	StartedAt           int64
	CanceledAt          int64
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	SupporterId         int64
	SupporterEmail      string
	CurrentPeriodEnd    int64
	CurrentPeriodStart  int64
}
