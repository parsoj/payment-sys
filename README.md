# TODOs 

- check transaction setup for Deposit and WithDrawl
  - check to/from fields
  - need repeatableread isolation? 


- check transaction setup for Transfer 
  - should txn row be created outside of db transaction

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

- Human-readable IDs 

## components of id string

- TODO id components "diagram"

## Design Choices

- TODO explain choices for timestamp portion

- TODO explain remainder of string, and why you didn't choose human-readable


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

## DB design
I'll just copy the sql code for bootstrapping the DB here (the sql bootstrapping files are localed in the DB folder)

## Features

### Idempotency
To support full idempotency - the client needs to provide an idempotency key when they submit a transaction. 
This key is checked against any existing transactions in the DB to ensure there are no accidental double transfers


### Concurrency
- TODO 



