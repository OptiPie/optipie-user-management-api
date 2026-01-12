# Deployment: Google OAuth JWT Auth + Anonymous Analytics

## Infrastructure Changes Required

### 1. Environment Variables
```bash
GOOGLE_OAUTH_CLIENT_ID=<your-client-id-from-google-console>
DYNAMODB_TABLE_NAME_ANALYTICS=optipie-analytics-production
```

### 2. Create DynamoDB Table
```bash
aws dynamodb create-table \
  --table-name optipie-analytics-production \
  --attribute-definitions AttributeName=timestamp,AttributeType=N \
  --key-schema AttributeName=timestamp,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST
```

### 3. New Endpoint
- **Route**: `POST /api/v1/analytics/collect`
- **Auth**: JWT (Bearer token in Authorization header)
- **Rate Limit**: 100 req/min per IP
- **Data**: Anonymous (strategy_name, strategy_symbol, strategy_period)

### 4. Test After Deployment
```bash
# Should return 401 (no token)
curl -X POST http://localhost:3000/api/v1/analytics/collect

# Should return 200 (with valid Google JWT)
curl -X POST http://localhost:3000/api/v1/analytics/collect \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"strategy_name":"test","strategy_symbol":"AAPL","strategy_period":"1d"}'
```

---

**Note**: Analytics is fully anonymous (GDPR compliant). No email or user ID stored.