# TODOs 
- set up golang nvim stuff
- find a library in go for human-readable string generation

- check transaction setup for Deposit and WithDrawl


- check transaction setup for Transfer 
  - should txn row be created outside of db transaction


- clean up commits and timestamps
- push to repo


# Part 1 - Unique Ids

- TODO id components "diagram"


- TODO explain choices for timestamp portion

- TODO explain remainder of string, and why you didn't choose human-readable



# Part 2 - Transaction Service


- TODO how to start the Service
  - TODO requirements to run the service

- TODO how to run Tests
  - go test for unit
  - hurl for API



## DB design
- TODO just list out rows and models

- TODO explain DB migrations left out

- TODO explain indexes left out

## Features

### Idempotency
- TODO explain implementation of Idempotency

### Concurrency
- TODO transaction isolation level



