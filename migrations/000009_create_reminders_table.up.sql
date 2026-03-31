CREATE TABLE reminders (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT,
  remind_at TIMESTAMP,
  type VARCHAR(50),
  created_at TIMESTAMP DEFAULT NOW(),

  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);