CREATE TABLE tasks (
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(200) NOT NULL,
                       active_at DATE NOT NULL,
                       status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'done')),
                       UNIQUE (title, active_at)
);