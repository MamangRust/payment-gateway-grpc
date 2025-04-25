## Auth

## User Sender

### Login
```sh
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword"
  }'
```


### Register

```sh
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "John",
    "lastname": "Doe",
    "email": "john.doe@example.com",
    "password": "securepassword",
    "confirm_password": "securepassword"
  }'
```

### Refresh Token

```sh
curl -X POST \
  http://localhost:5000/api/auth/refresh-token \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM1Nzg2MTYxfQ.yEx98MCuT0fg8b63VuLl9XcPxszYG2BTlQtRVvEsMbI' \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMSIsImF1ZCI6WyJyZWZyZXNoIl0sImV4cCI6MTczNTc4NjE2MX0.Ti5BTb8xMbMUYDNE-vFU8MVbr6o7zQLWJ-CIetByFd4"
}'
```

### GetMe

```sh
curl -X GET http://localhost:5000/api/auth/me \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM1ODI5MDk0fQ.0MAChuYO1G458hK_HqVmFYAOdOnmeYkqFbTjbY0QDi8'

```


## User Receiver

### Login
```sh
curl -X POST http://0.0.0.0:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
     "email": "jane.doe@example.com",
  "password": "password123"
  }'
```


### Register

```sh
curl -X POST http://0.0.0.0:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
      "firstname": "Jane",
  "lastname": "Doe",
  "email": "jane.doe@example.com",
  "password": "password123",
  "confirm_password": "password123"
  }'
```

--------------------------------------------------

## User

## User Sender

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM1OTk0NDIxfQ.IFWGbahWa3VqrYc-M77KYJI9Q13rjRL9IPWeqw-P7Rs"
```

### FindById
```sh
curl -X GET http://0.0.0.0:5000/api/user/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NzE2ODczfQ.OsiUy5EuQVZLD8QnDXtRTlypHGqewz2x2J6utbH4bUg"
```


### FindByActive
```sh
curl -X GET http://0.0.0.0:5000/api/user/active \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### FindByTrashed

```sh
curl -X GET http://0.0.0.0:5000/api/user/trashed \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### Create

```sh
curl -X POST http://0.0.0.0:5000/api/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NjYxODEzfQ.LD51QxG5lfPUmkOw2L0yfQ8VjSzLHh2U0R1Dq7ca944" \
  -d '{
    "firstname": "Jane",
  "lastname": "Doe",
  "email": "jane.doe@example.com",
  "password": "password123",
  "confirm_password": "password123"
}'
```

### Update
```sh
curl -X POST http://0.0.0.0:5000/api/user/update/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NjYxODEzfQ.LD51QxG5lfPUmkOw2L0yfQ8VjSzLHh2U0R1Dq7ca944" \
  -d '{
    "user_id": 42,
    "firstname": "John",
    "lastname": "Doe",
    "email": "john.updated@example.com",
    "password": "newsecurepassword",
    "confirm_password": "newsecurepassword"
}'
```


### Trashed

```sh
curl -X POST http://0.0.0.0:5000/api/user/trashed/42 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NjYxODEzfQ.LD51QxG5lfPUmkOw2L0yfQ8VjSzLHh2U0R1Dq7ca944"
```

### Restore

```sh
curl -X POST http://0.0.0.0:5000/api/user/restore/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NjYxODEzfQ.LD51QxG5lfPUmkOw2L0yfQ8VjSzLHh2U0R1Dq7ca944"
```

### Permanent

```sh
curl -X DELETE hthttp://0.0.0.0:5000/api/user/permanent/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

## User Receiver

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### FindById
```sh
curl -X GET http://0.0.0.0:5000/api/user/22 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NzE2ODczfQ.OsiUy5EuQVZLD8QnDXtRTlypHGqewz2x2J6utbH4bUg"
```


### FindByActive
```sh
curl -X GET http://0.0.0.0:5000/api/user/active \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### FindByTrashed

```sh
curl -X GET http://0.0.0.0:5000/api/user/trashed \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### Create

```sh
curl -X POST http://0.0.0.0:5000/api/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s" \
  -d '{
    "firstname": "Jane",
  "lastname": "Doe",
  "email": "jane.doe@example.com",
  "password": "password123",
  "confirm_password": "password123"
}'
```

### Update
```sh
curl -X POST http://0.0.0.0:5000/api/user/update/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s" \
  -d '{
    "user_id": 2,
    "firstname": "John",
    "lastname": "Doe",
    "email": "john.updated@example.com",
    "password": "newsecurepassword",
    "confirm_password": "newsecurepassword"
}'
```


### Trashed

```sh
curl -X POST http://0.0.0.0:5000/api/user/trashed/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### Restore

```sh
curl -X POST http://0.0.0.0:5000/api/user/restore/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

### Permanent

```sh
curl -X DELETE hthttp://0.0.0.0:5000/api/user/permanent/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiMSIsImV4cCI6MTczNDY1MTIzN30.APhuIsM2DIyUaLZQLoapyJqsbghAdW155bFwhaxM1_s"
```

------------------------------------------------------------------------------


## Card

## User Sender

### FindAll

```sh
curl -X GET "http://0.0.0.0:5000/api/card" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM1OTk2Mjk5fQ.m_2bE3hpoTdwQ2_B9AjxchbUFAqcxls5vIrh4le2Yyo" \
-H "Content-Type: application/json"
```

### FindById

```sh
curl -X GET "http://0.0.0.0:5000/api/card/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindByUserId
```sh
curl -X GET "http://0.0.0.0:5000/api/card/user" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/card/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/card/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/card/card_number/1234567890" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/card/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "user_id": 1,
  "card_type": "Credit",
  "expire_date": "2025-12-31T00:00:00Z",
  "cvv": "123",
  "card_provider": "Visa"
}'

```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/card/update/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_id": 1,
  "user_id": 1,
  "card_type": "Debit",
  "expire_date": "2026-06-30T00:00:00Z",
  "cvv": "456",
  "card_provider": "MasterCard"
}'

```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/card/trashed/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/card/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/card/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```
## User Receiver

### FindAll

```sh
curl -X GET "http://0.0.0.0:5000/api/card" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById

```sh
curl -X GET "http://0.0.0.0:5000/api/card/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindByUserId
```sh
curl -X GET "http://0.0.0.0:5000/api/card/user" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/card/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/card/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/card/card_number/1234567890" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/card/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "user_id": 2,
  "card_type": "Credit",
  "expire_date": "2025-12-31T00:00:00Z",
  "cvv": "123",
  "card_provider": "Visa"
}'

```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/card/update/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_id": 2,
  "user_id": 2,
  "card_type": "Debit",
  "expire_date": "2026-06-30T00:00:00Z",
  "cvv": "456",
  "card_provider": "MasterCard"
}'

```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/card/trashed/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/card/restore/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/card/permanent/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

--------------------------------------------------------------------

## Saldo

## User Sender
### FindAll

```sh
curl -X GET "http://0.0.0.0:5000/api/saldo" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById

```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/card_number/1234567890" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "1234567890",
  "total_balance": 5000
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/update/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "saldo_id": 1,
  "card_number": "1234567890",
  "total_balance": 10000
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/trashed/123" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/restore/123" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/saldo/permanent/123" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

## User Receiver

### FindAll

```sh
curl -X GET "http://0.0.0.0:5000/api/saldo" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById

```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/2" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM1ODI5MDk0fQ.0MAChuYO1G458hK_HqVmFYAOdOnmeYkqFbTjbY0QDi8" \
-H "Content-Type: application/json"
```


### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/saldo/card_number/1234567890" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "1234567890",
  "total_balance": 5000
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/update/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "saldo_id": 2,
  "card_number": "1234567890",
  "total_balance": 10000
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/trashed/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/saldo/restore/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/saldo/permanent/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

-----------------------------------------------------

## Merchant

### FindAll

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindByName

```sh

```

### Api Key

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant/api-key" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindByMerchanUserId

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant/merchant-user/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Active

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed

```sh
curl -X GET "http://0.0.0.0:5000/api/merchant/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create

```sh
curl -X POST "http://0.0.0.0:5000/api/merchant/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "name": "New Merchant",
  "user_id": 1
}'
```

### Update

```sh
curl -X POST "http://0.0.0.0:5000/api/merchants/updates/1" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3NzE2ODczfQ.OsiUy5EuQVZLD8QnDXtRTlypHGqewz2x2J6utbH4bUg" \
-H "Content-Type: application/json" \
-d '{
  "merchant_id": 1,
  "name": "Updated Merchant",
  "user_id": 1,
  "status": "active"
}'
```

### Trashed

```sh
curl -X POST "http://0.0.0.0:5000/api/merchant/trashed/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore

```sh
curl -X POST "http://0.0.0.0:5000/api/merchant/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent

```sh
curl -X DELETE "http://0.0.0.0:5000/api/merchant/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

--------------------------------------------------------------------------------------

## Topup

## Topup Sender

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/topup" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/card_number/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```


### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "topup_no": "TOPUP001",
  "topup_amount": 50000,
  "topup_method": "Bank Transfer"
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/update/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "topup_id": 1,
  "topup_amount": 75000,
  "topup_method": "Credit Card"
}'

```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/trashed/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/topup/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

## Topup Receiver

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/topup" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/topup/card_number/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```


### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "topup_no": "TOPUP001",
  "topup_amount": 50000,
  "topup_method": "Bank Transfer"
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/update/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "topup_id": 2,
  "topup_amount": 75000,
  "topup_method": "Credit Card"
}'

```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/trashed/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/topup/restore/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/topup/permanent/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

--------------------------------------------------------------------------------



## Transaction

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/transaction" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/transactions/1" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MSIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzM3ODk0NTkwfQ.xiW8tYBVV2PK32vSM1Q5ntNy2vPfasdMLPYlLHWU62M" \
-H "Content-Type: application/json"
```


### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/transaction/card/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindByMerchantId
```sh
curl -X GET "http://0.0.0.0:5000/api/transaction/merchant/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/transaction/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/transaction/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/transaction/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "amount": 50000,
  "payment_method": "Credit Card",
  "merchant_id": 2,
  "transaction_time": "2024-06-17T15:04:05Z"
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/transaction/update" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "transaction_id": 1,
  "card_number": "123456789",
  "amount": 75000,
  "payment_method": "Bank Transfer",
  "merchant_id": 3,
  "transaction_time": "2024-06-18T12:00:00Z"
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/transaction/trashed/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/transaction/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/transaction/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

------------------------------------------------------------------------------


## Transfer

## User Sender
### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Transfer From
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/transfer_from/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Transfer To
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/transfer_to/987654321" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```


### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "transfer_from": "123456789",
  "transfer_to": "987654321",
  "transfer_amount": 50000
}'

```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/update/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "transfer_id": 1,
  "transfer_from": "123456789",
  "transfer_to": "987654321",
  "transfer_amount": 75000
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/trashed/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/transfer/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

## User Receiver

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Transfer From
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/transfer_from/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Transfer To
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/transfer_to/987654321" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```


### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/transfer/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```


### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "transfer_from": "123456789",
  "transfer_to": "987654321",
  "transfer_amount": 50000
}'

```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/update/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "transfer_id": 2,
  "transfer_from": "123456789",
  "transfer_to": "987654321",
  "transfer_amount": 75000
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/trashed/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/transfer/restore/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/transfer/permanent/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```



----------------------------------------------------------------------

## Withdraw

## User Sender

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/card_number/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "withdraw_amount": 50000,
  "withdraw_time": "2024-06-11T10:00:00Z"
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/update/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "withdraw_id": 1,
  "card_number": "123456789",
  "withdraw_amount": 75000,
  "withdraw_time": "2024-06-11T12:00:00Z"
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/trash/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/restore/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/withdraw/permanent/1" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```
## User Receiver

### FindAll
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### FindById
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Card Number
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/card_number/123456789" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Active
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/active" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Trashed
```sh
curl -X GET "http://0.0.0.0:5000/api/withdraw/trashed" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Create
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/create" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "card_number": "123456789",
  "withdraw_amount": 50000,
  "withdraw_time": "2024-06-11T10:00:00Z"
}'
```

### Update
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/update/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "withdraw_id": 2,
  "card_number": "123456789",
  "withdraw_amount": 75000,
  "withdraw_time": "2024-06-11T12:00:00Z"
}'
```

### Trashed
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/trash/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Restore
```sh
curl -X POST "http://0.0.0.0:5000/api/withdraw/restore/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"
```

### Permanent
```sh
curl -X DELETE "http://0.0.0.0:5000/api/withdraw/permanent/2" \
-H "Authorization: Bearer <YOUR_BEARER_TOKEN>" \
-H "Content-Type: application/json"

```


----------------------------------------------------------------------------------------
