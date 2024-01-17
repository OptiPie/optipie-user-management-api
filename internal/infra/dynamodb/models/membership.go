package dbmodels

type Membership struct {
	Email               string  `dynamodbav:"email"`
	Type                string  `dynamodbav:"type"`
	LiveMode            bool    `dynamodbav:"live_mode"`
	Attempt             int32   `dynamodbav:"attempt"`
	Created             int64   `dynamodbav:"created"`
	EventId             int64   `dynamodbav:"event_id,omitempty"`
	Id                  int64   `dynamodbav:"id"`
	Amount              float64 `dynamodbav:"amount"`
	Object              string  `dynamodbav:"object"`
	Paused              string  `dynamodbav:"paused"`
	Status              string  `dynamodbav:"status"`
	Canceled            string  `dynamodbav:"canceled"`
	Currency            string  `dynamodbav:"currency,omitempty"`
	PspId               string  `dynamodbav:"psp_id,omitempty"`
	MembershipLevelId   string  `dynamodbav:"membership_level_id"`
	MembershipLevelName string  `dynamodbav:"membership_level_name"`
	StartedAt           int64   `dynamodbav:"started_at"`
	CanceledAt          int64   `dynamodbav:"canceled_at"`
	NoteHidden          bool    `dynamodbav:"note_hidden,omitempty"`
	SupportNote         string  `dynamodbav:"support_note,omitempty"`
	SupporterName       string  `dynamodbav:"supporter_name,omitempty"`
	SupporterId         int64   `dynamodbav:"supporter_id"`
	CurrentPeriodEnd    int64   `dynamodbav:"current_period_end"`
	CurrentPeriodStart  int64   `dynamodbav:"current_period_start"`
}
