# How to run
## Prerequisites
To run the service, you should just need Make and Docker installed

For running unit tests - its just "go test" - so you'd need a local golang toolchain

For running the api test - you'd need hurl installed (https://hurl.dev)

## Run the service
```
make bootstrap-all
```
You should be able to use the makefile and just run "make bootstrap-all" 
to run both the Postgres container and service container

The Makefile and its comments should also provide a reasonable guide for how everything is set up to run, 
and the various operations you have available for the codebase. 



## Tests 
See the Makefile for the test functions available

# Part 1 - Unique Ids

## What I left out

### Human-readable component to IDs 
I originally implemented this, where the least-significant characters of the string were actually
picked from a dictionary of words (and the remaining characters that weren't the timestamp were random)
so it looked like "23JcebKwFMuUKNA-dogs" 
components of:
"23JcebKwFMu" - timestamp
"UKNA" - random filler
"-dogs" - human readable word with separator

I removed this feature, in part because I needed a separator to distinguish the word from the rest of the id, and I was using a special character to do that (which I wanted to avoid). Also - the number of random characters was pretty small, depending on the length of the readable word (sometimes only 2 or 1 character) - which made me concerned about guess-ability of the id

Also - in a prod system - making a choice to defer a feature like that actually wouldn't necessarily "lock you in" to an identifier schema that wasn't human readable. As long as the "most significant" characters remained the same style of timestamp - you could continue to tweak how the lower-order characters are generated, without affecting the ordering between the new and old ids

## components of id string
- "most significant" 11 characters: 
  - Unix Nanosecond timestamp (base 62 encoded) 

- remaining characters are random

## Design Choices
I chose the Nanosecond timestamp, because in my testing there were too many conflicts with millisecond
that does have the disadvantage of consuming more characters, but minimizing ids with identical timestamps seems worth it. 

I chose base 62 encoding, since that utilizes all alphanumeric characters 
(including upper and lower case alphabetic characters) - which allows us to minimize the length of the timestamp as much as possible. 

Although base64 is quite common, it adds the "+" and "/" characters, and there isn't really a consistent standard for how to lex-sort special characters and/or base64

by making the remaining characters random (or random with a dictionary word on the end) - that can function to tie-break the occasional case where two ids have identical timestamps. That way, two IDs created at or close to the exact same nanosecond will still have a strict ordering and one will be first in lex order

# Part 2 - Transaction Service

## What I left out
### List Transactions call
To implement the core functionality around transactions - I mainly just needed a 
GetTransaction function, so I implemented that to start. 

Implementing an additional API call to do pagination and range queries over the transactions table wouldn't 
require any major reword of the core logic, though, so IMO in a "real world" it wouldn't add technical debt to defer this feature 
in the interest of getting the core logic out. 

### GetBalance at different times - and Event Sourcing and for reconciliation 
I don't have an Event-sourcing setup here, which I know is common of payment systems to 
support reconciliation. 

Implementing this would be a substantial architectural shift over what I have here, but I mostly left this out since it wasn't explicitly called out as a requirement and I needed to save time! 

However implementing an event sourcing type pattern would also allow for retrieving account balances at any point in history

### BenchMark Tests & tests for race-conditions
In a real-world payment system, IMO a substantial investment into these sorts of stress tests would be justified. 
But, I just didn't have time so my testing is fairly minimal and mostly serves to check for correctness. 

I believe I took the right approaches for ensuring consistency under concurrent operations - but of course that should be validated in a real production system. 

### DB Migration tooling 
I just threw the db bootstrapping into an sql file, and just create a fresh DB every time I need to test the service.
this wouldn't fly in prod - and we'd need to use a DB migration tool of some sort

## DB - data model 
I'll just copy the sql code for bootstrapping the DB here (the sql bootstrapping files are localed in the DB folder)
primary keys are fixed-length to 20 chars - to be as efficient as possible while using the id generated code discussed above

```sql 
CREATE DATABASE transactions_db;

\connect transactions_db transaction_svc

CREATE TABLE IF NOT EXISTS users (
  id CHAR(20) PRIMARY KEY,   
  username CHAR(20) ,   
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS accounts (
  id CHAR(20) PRIMARY KEY,   
  user_id CHAR(20) ,   
	balance DECIMAL(15, 2) DEFAULT 0.00,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
  id CHAR(20) PRIMARY KEY,   
	to_account CHAR(20),
	from_account CHAR(20),
	amount DECIMAL(15, 2) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	idempotency_key TEXT UNIQUE,
	state TEXT
);

```
## Features

### Idempotency
To support full idempotency - the client needs to generate and provide an idempotency key to the API. 
This key is written to the transaction table when the transaction is run. 

The transaction table has the column for this key set as unique - so any subsequent transactions that are started with the same idempotency key will fail due to the uniqueness constraint. 

Also - any concurrent transactions (with the same idempotency key) that are triggered while the original transaction is in flight - will be blocked until the first transaction either commits or rolls back. This ensures even in the concurrent case, at most one of these transactions will succeed. 



