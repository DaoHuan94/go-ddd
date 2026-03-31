CREATE TABLE comments (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT,
  user_id BIGINT,
  content TEXT,
  created_at TIMESTAMP DEFAULT NOW(),

  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id)
);