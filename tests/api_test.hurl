# Create a user and capture the user ID
POST http://localhost:8080/create-user
Content-Type: application/json

{
  "username": "user1"
}

HTTP/1.1 200
[Captures]
user1_id: jsonpath "$.user.ID"

###

# Create a user and capture the user ID
POST http://localhost:8080/create-user
Content-Type: application/json

{
  "username": "user2"
}

HTTP/1.1 200
[Captures]
user2_id: jsonpath "$.user.ID"

###

# Test the CreateAccount endpoint for user1
POST http://localhost:8080/create-account
Content-Type: application/json

{
  "user_id": "{{user1_id}}"
}

HTTP/1.1 200
[Captures]
account1_id: jsonpath "$.account.ID"


###
#
# Test the CreateAccount endpoint for user1
POST http://localhost:8080/create-account
Content-Type: application/json

{
  "user_id": "{{user2_id}}"
}

HTTP/1.1 200
[Captures]
account2_id: jsonpath "$.account.ID"


###

# Deposit funds into the first account
POST http://localhost:8080/deposit-funds
Content-Type: application/json

{
  "account_id": "{{account1_id}}",
  "amount": 100.0
}

HTTP/1.1 200

###

# Transfer funds from the first account to the second account
POST http://localhost:8080/transfer-funds
Content-Type: application/json

{
  "from_account": "{{account1_id}}",
  "to_account": "{{account2_id}}",
  "amount": 50.0
}

HTTP/1.1 200

###

# Get account1 details
GET http://localhost:8080/get-account?account_id={{account1_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.account.Balance" == 50.0

###

# Get account2 details
GET http://localhost:8080/get-account?account_id={{account2_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.account.Balance" == 50.0
