CREATE TABLE IF NOT EXISTS velocity_limits ( velocity_limit_id VARCHAR(255) NOT NULL, site_id INT NOT NULL, tid_id INT NULL, limit_level INT NOT NULL, scheme INT NULL, transaction_limit_daily INT NOT NULL, transaction_limit_batch INT NOT NULL, single_transaction_limit INT NOT NULL, PRIMARY KEY (velocity_limit_id), FOREIGN KEY (limit_level) REFERENCES limit_levels(limit_id), FOREIGN KEY (scheme) REFERENCES schemes(scheme_id));