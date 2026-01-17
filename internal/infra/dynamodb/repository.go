package dynamorepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/cerrors"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"github.com/OptiPie/optipie-user-management-api/internal/domain/models"
	dbmodels "github.com/OptiPie/optipie-user-management-api/internal/infra/dynamodb/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

const (
	membershipPrimaryKey = "email"
	analyticsPrimaryKey  = "timestamp"

	// dynamodb condition expressions
	attributeNotExists = "attribute_not_exists"
	attributeExists    = "attribute_exists"
)

type Client struct {
	client              *dynamodb.Client
	membershipTableName string
	analyticsTableName  string
}

func NewRepository(client *dynamodb.Client, membershipTableName, analyticsTableName string) domain.Repository {
	return &Client{
		client:              client,
		membershipTableName: membershipTableName,
		analyticsTableName:  analyticsTableName,
	}
}

func (c *Client) CreateMembership(ctx context.Context, args domain.CreateMembershipArgs) error {
	membership := dbmodels.Membership{
		Email:               args.SupporterEmail,
		Type:                args.Type,
		LiveMode:            args.LiveMode,
		Attempt:             args.Attempt,
		Created:             args.Created,
		Updated:             args.Created,
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
		return fmt.Errorf("marshall_map error: %v", err)
	}
	conditionExpression := fmt.Sprintf("%v(%v)", attributeNotExists, membershipPrimaryKey)

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.membershipTableName), Item: item,
		ConditionExpression: aws.String(conditionExpression),
	})
	if err != nil {
		return fmt.Errorf("put_item error: %v", err)
	}

	return nil
}

func (c *Client) GetMembershipByEmail(ctx context.Context, email string) (models.Membership, error) {
	membership := dbmodels.Membership{Email: email}
	membershipPk, err := membership.GetPrimaryKey()
	if err != nil {
		return models.Membership{}, fmt.Errorf("get_primary_key error: %v", err)
	}
	response, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: membershipPk, TableName: aws.String(c.membershipTableName),
	})

	if err != nil {
		return models.Membership{}, fmt.Errorf("get_item error: %v", err)
	}

	err = attributevalue.UnmarshalMap(response.Item, &membership)
	if err != nil {
		return models.Membership{}, fmt.Errorf("unmarshal_map error: %v", err)
	}

	return models.Membership{
		Type:                membership.Type,
		LiveMode:            membership.LiveMode,
		Attempt:             membership.Attempt,
		Created:             membership.Created,
		EventId:             membership.EventId,
		Id:                  membership.Id,
		Amount:              membership.Amount,
		Object:              membership.Object,
		Paused:              membership.Paused,
		Status:              membership.Status,
		Canceled:            membership.Canceled,
		Currency:            membership.Currency,
		PspId:               membership.PspId,
		MembershipLevelId:   membership.MembershipLevelId,
		MembershipLevelName: membership.MembershipLevelName,
		StartedAt:           membership.StartedAt,
		CanceledAt:          membership.CanceledAt,
		NoteHidden:          membership.NoteHidden,
		SupportNote:         membership.SupportNote,
		SupporterName:       membership.SupporterName,
		SupporterId:         membership.SupporterId,
		SupporterEmail:      membership.Email,
		CurrentPeriodEnd:    membership.CurrentPeriodEnd,
		CurrentPeriodStart:  membership.CurrentPeriodStart,
	}, nil
}

func (c *Client) UpdateMembershipByEmail(ctx context.Context, email string, args domain.UpdateMembershipArgs) error {
	update := expression.Set(expression.Name("updated"), expression.Value(time.Now().UTC()))
	update.Set(expression.Name("paused"), expression.Value(args.Paused))
	update.Set(expression.Name("status"), expression.Value(args.Status))
	update.Set(expression.Name("canceled"), expression.Value(args.Canceled))
	update.Set(expression.Name("currency"), expression.Value(args.Currency))
	update.Set(expression.Name("psp_id"), expression.Value(args.PspId))
	update.Set(expression.Name("membership_level_id"), expression.Value(args.MembershipLevelId))
	update.Set(expression.Name("membership_level_name"), expression.Value(args.MembershipLevelName))
	update.Set(expression.Name("started_at"), expression.Value(args.StartedAt))
	update.Set(expression.Name("canceled_at"), expression.Value(args.CanceledAt))
	update.Set(expression.Name("note_hidden"), expression.Value(args.NoteHidden))
	update.Set(expression.Name("support_note"), expression.Value(args.SupportNote))
	update.Set(expression.Name("supporter_name"), expression.Value(args.SupporterName))
	update.Set(expression.Name("current_period_end"), expression.Value(args.CurrentPeriodEnd))
	update.Set(expression.Name("supporter_feedback"), expression.Value(args.SupporterFeedback))
	update.Set(expression.Name("cancel_at_period_end"), expression.Value(args.CancelAtPeriodEnd))
	update.Set(expression.Name("current_period_start"), expression.Value(args.CurrentPeriodStart))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return fmt.Errorf("expression_builder error: %v", err)
	}

	membership := dbmodels.Membership{Email: email}
	membershipPk, err := membership.GetPrimaryKey()
	if err != nil {
		return fmt.Errorf("get_primary_key error: %v", err)
	}

	conditionExpression := fmt.Sprintf("%v(%v)", attributeExists, membershipPrimaryKey)

	_, err = c.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(c.membershipTableName),
		Key:                       membershipPk,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String(conditionExpression),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueNone,
	})

	if err != nil {
		return fmt.Errorf("update_item error: %v", err)
	}

	return nil
}

func (c *Client) DeleteMembershipByEmail(ctx context.Context, email string) error {
	membership := dbmodels.Membership{Email: email}
	membershipPk, err := membership.GetPrimaryKey()
	if err != nil {
		return fmt.Errorf("get_primary_key error: %v", err)
	}

	conditionExpression := fmt.Sprintf("%v(%v)", attributeExists, membershipPrimaryKey)

	_, err = c.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:           aws.String(c.membershipTableName),
		Key:                 membershipPk,
		ConditionExpression: aws.String(conditionExpression),
	})

	if err != nil {
		var ccf *types.ConditionalCheckFailedException
		if errors.As(err, &ccf) {
			return cerrors.NewCustomError(err.Error(), cerrors.ConditionalCheckFailedException)

		}
		return fmt.Errorf("delete_item error: %v", err)
	}

	return nil
}

func (c *Client) CreateAnalytics(ctx context.Context, args domain.CreateAnalyticsArgs) error {
	analytics := dbmodels.Analytics{
		Timestamp:         args.Timestamp,
		StrategyName:      args.StrategyName,
		StrategySymbol:    args.StrategySymbol,
		StrategyPeriod:    args.StrategyPeriod,
		StrategyDateRange: args.StrategyDateRange,
	}

	item, err := attributevalue.MarshalMap(analytics)
	if err != nil {
		return fmt.Errorf("marshall_map error: %v", err)
	}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.analyticsTableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put_item error: %v", err)
	}

	return nil
}
