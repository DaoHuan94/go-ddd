CREATE TABLE tasks (
  id BIGSERIAL PRIMARY KEY,
  project_id BIGINT,
  parent_task_id BIGINT NULL,

  title VARCHAR(255) NOT NULL,
  description TEXT,

  status VARCHAR(50),
  priority VARCHAR(50),

  assigned_to BIGINT NULL,
  created_by BIGINT,

  due_date TIMESTAMP NULL,
  start_date TIMESTAMP NULL,

  position INT,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (assigned_to) REFERENCES users(id),
  FOREIGN KEY (created_by) REFERENCES users(id),
  FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE
);