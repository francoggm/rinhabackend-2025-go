-- Ensure the schema is created in the right database
\c payments;

CREATE TYPE processor_type AS ENUM ('default', 'fallback');

CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  correlation_id UUID NOT NULL,
  amount NUMERIC(10, 2) NOT NULL,
  type processor_type NOT NULL,
  requested_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_payments_type_timestamp ON payments(type, requested_at);
