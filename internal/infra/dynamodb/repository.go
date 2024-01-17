package dynamorepo

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"github.com/OptiPie/optipie-user-management-api/internal/domain/models"
	dbmodels "github.com/OptiPie/optipie-user-management-api/internal/infra/dynamodb/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	primaryKey = "email"

	// dynamodb condition expressions
	attributeNotExists = "attribute_not_exists"
)

type Client struct {
	client *dynamodb.Client
}

func NewRepository(client *dynamodb.Client) domain.Repository {
	return &Client{
		client: client,
	}
}

func (c Client) CreateMembership(ctx context.Context, args domain.CreateMembershipArgs) error {
	membership := dbmodels.Membership{
		Email:               args.SupporterEmail,
		Type:                args.Type,
		LiveMode:            args.LiveMode,
		Attempt:             args.Attempt,
		Created:             args.Created,
		EventId:             args.EventId,
		Id:                  args.Id,
		Amount:              args.Amount,
		Object:              args.Object,
		Paused:              args.Paused,
		Status:              args.Status,
		Canceled:            args.Canceled,
		Currency:            args.Currency,
		PspId:               args.PspId,
		MembershipLevelId:   args.MembershipLevelId,
		MembershipLevelName: args.MembershipLevelName,
		StartedAt:           args.StartedAt,
		CanceledAt:          args.CanceledAt,
		NoteHidden:          args.NoteHidden,
		SupportNote:         args.SupportNote,
		SupporterName:       args.SupporterName,
		SupporterId:         args.SupporterId,
		CurrentPeriodEnd:    args.CurrentPeriodEnd,
		CurrentPeriodStart:  args.CurrentPeriodStart,
	}

	item, err := attributevalue.MarshalMap(membership)
	if err != nil {
		return fmt.Errorf("marshall error: %v", err)
	}
	conditionExpression := fmt.Sprintf("%v(%v)", attributeNotExists, primaryKey)

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("test"), Item: item,
		ConditionExpression: aws.String(conditionExpression),
	})
	if err != nil {
		return fmt.Errorf("put item error: %v", err)
	}

	return nil
}

func (c Client) GetMembershipByEmail(ctx context.Context, email string) (models.Membership, error) {
	//TODO implement me
	panic("implement me")
}