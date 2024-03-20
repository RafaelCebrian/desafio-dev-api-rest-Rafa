CREATE TABLE IF NOT EXISTS holders (
    holder_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS accounts (
	account_id SERIAL PRIMARY KEY,
    fk_holder VARCHAR(11) REFERENCES holders(cpf),
    number VARCHAR(8) NOT null UNIQUE,
    agency VARCHAR(3) NOT NULL,
    balance NUMERIC(10, 2) NOT NULL, 
    blocked BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS operations (
    operation_id SERIAL PRIMARY KEY,
    fk_account VARCHAR(8) REFERENCES accounts(number),
    type VARCHAR(50) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50)
);

DROP TABLE IF EXISTS operations
DROP TABLE IF EXISTS accounts
DROP TABLE IF EXISTS holders