CREATE TABLE IF NOT EXISTS minute_prices (
    exchange TEXT NOT NULL,
    pair_name TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    average_price DOUBLE PRECISION NOT NULL,
    min_price DOUBLE PRECISION NOT NULL,
    max_price DOUBLE PRECISION NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_minute_prices_sym_ex_ts ON minute_prices (pair_name, exchange, timestamp DESC);