ALTER TABLE services
  ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN delete_time TIMESTAMPTZ DEFAULT NULL;