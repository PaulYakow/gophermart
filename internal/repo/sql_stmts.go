package repo

// todo: Named stmt place here

const (
	schemaUsers = `
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    login         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    UNIQUE (login)
);
`
	schemaOrders = `
CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    number      BIGINT NOT NULL UNIQUE,
    status      VARCHAR(255),
    accrual     REAL,
    uploaded_at TIMESTAMP
);
`
	schemaLinkUsersOrders = `
CREATE TABLE IF NOT EXISTS users_orders
(
    id       SERIAL PRIMARY KEY,
    user_id  SERIAL,
    order_id SERIAL
);
`
	insertUser = `
INSERT INTO users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	selectUser = `
SELECT id
FROM users
WHERE login=$1 AND password_hash=$2;
`
)
