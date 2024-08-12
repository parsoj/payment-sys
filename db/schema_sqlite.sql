
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
	idempotency_key TEXT,
	state TEXT
);
