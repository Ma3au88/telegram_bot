CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    chat_id NUMERIC NOT NULL,
    date TIMESTAMP(0) NOT NULL
);

CREATE TABLE crypto_prices (
   id SERIAL PRIMARY KEY,
   coin TEXT NOT NULL,
   value NUMERIC NOT NULL,
   date TIMESTAMP(0) NOT NULL
);

CREATE TABLE subscribers (
    id SERIAL PRIMARY KEY,
    chat_id NUMERIC UNIQUE NOT NULL,
    interval INTERVAL(0) NOT NULL,
    last_message_time TIMESTAMP(0),
    status BOOLEAN NOT NULL DEFAULT FALSE
);