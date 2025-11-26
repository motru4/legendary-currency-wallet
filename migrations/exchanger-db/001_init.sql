-- +goose Up
-- +goose StatementBegin
CREATE TABLE exchange_rates (
    id SERIAL PRIMARY KEY,
    base_currency VARCHAR(3) NOT NULL,      
    target_currency VARCHAR(3) NOT NULL,      
    rate DECIMAL(15, 6) NOT NULL,           
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_different_currencies CHECK (base_currency != target_currency),
    CONSTRAINT check_positive_rate CHECK (rate > 0),
    
    CONSTRAINT unique_currency_pair UNIQUE (base_currency, target_currency)
);

CREATE INDEX idx_exchange_rates_base_currency ON exchange_rates(base_currency);
CREATE INDEX idx_exchange_rates_target_currency ON exchange_rates(target_currency);
CREATE INDEX idx_exchange_rates_updated_at ON exchange_rates(updated_at);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_exchange_rates_updated_at 
    BEFORE UPDATE ON exchange_rates 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

INSERT INTO exchange_rates (base_currency, target_currency, rate) VALUES
('USD', 'EUR', 0.85),
('USD', 'RUB', 90.0),
('EUR', 'USD', 1.18),
('EUR', 'RUB', 105.0),
('RUB', 'USD', 0.011),
('RUB', 'EUR', 0.0095);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchange_rates;
DROP TRIGGER IF EXISTS update_exchange_rates_updated_at ON exchange_rates;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd