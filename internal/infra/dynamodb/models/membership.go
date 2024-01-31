package dbmodels

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type Membership struct {
	Email               string    `dynamodbav:"email"`
	Type                string    `dynamodbav:"type"`
	LiveMode            bool      `dynamodbav:"live_mode"`
	Attempt             int32     `dynamodbav:"attempt"`
	Created             time.Time `dynamodbav:"created"`
	Updated             time.Time `dynamodbav:"updated"`
	EventId             int64     `dynamodbav:"event_id,omitempty"`
	Id                  int64     `dynamodbav:"id"`
	Amount              float64   `dynamodbav:"amount"`
	Object              string    `dynamodbav:"object"`
	Paused              string    `dynamodbav:"paused"`
	Status              string    `dynamodbav:"status"`
	Canceled            string    `dynamodbav:"canceled"`
	Currency            string    `dynamodbav:"currency,omitempty"`
	PspId               string    `dynamodbav:"psp_id,omitempty"`
	MembershipLevelId   int64     `dynamodbav:"membership_level_id"`
	MembershipLevelName string    `dynamodbav:"membership_level_name"`
	StartedAt           time.Time `dynamodbav:"started_at"`
	CanceledAt          time.Time `dynamodbav:"canceled_at,omitempty"`
	NoteHidden          bool      `dynamodbav:"note_hidden,omitempty"`
	SupportNote         string    `dynamodbav:"support_note,omitempty"`
	SupporterName       string    `dynamodbav:"supporter_name,omitempty"`
	SupporterId         int64     `dynamodbav:"supporter_id"`
	CurrentPeriodEnd    time.Time `dynamodbav:"current_period_end"`
	SupporterFeedback   string    `dynamodbav:"supporter_feedback,omitempty"`
	CancelAtPeriodEnd   string    `dynamodbav:"cancel_at_period_end,omitempty"`
	CurrentPeriodStart  time.Time `dynamodbav:"current_period_start"`
}

func (m Membership) GetPrimaryKey() (map[string]types.AttributeValue, error) {
	email, err := attributevalue.Marshal(m.Email)
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{"email": email}, nil
}
