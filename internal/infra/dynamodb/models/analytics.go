package dbmodels

import (
	"github.com/OptiPie/optipie-user-management-api/internal/domain/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Analytics struct {
	Timestamp      int64  `dynamodbav:"timestamp"`
	StrategyName   string `dynamodbav:"strategy_name"`
	StrategySymbol string `dynamodbav:"strategy_symbol"`
	StrategyPeriod string `dynamodbav:"strategy_period"`
}

func (a *Analytics) GetKey() map[string]types.AttributeValue {
	timestamp, err := attributevalue.Marshal(a.Timestamp)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{
		"timestamp": timestamp,
	}
}

func (a *Analytics) ToDomain() *models.Analytics {
	return &models.Analytics{
		Timestamp:      a.Timestamp,
		StrategyName:   a.StrategyName,
		StrategySymbol: a.StrategySymbol,
		StrategyPeriod: a.StrategyPeriod,
	}
}

func FromDomain(analytics *models.Analytics) *Analytics {
	return &Analytics{
		Timestamp:      analytics.Timestamp,
		StrategyName:   analytics.StrategyName,
		StrategySymbol: analytics.StrategySymbol,
		StrategyPeriod: analytics.StrategyPeriod,
	}
}
