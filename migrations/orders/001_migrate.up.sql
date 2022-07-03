CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customerid STRING,
    status STRING,
    createdon DATE,
    amount FLOAT
);