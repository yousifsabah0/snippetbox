-- Create database
CREATE DATABASE IF NOT EXISTS snippetbox SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to 'snippetbox' database
USE snippetbox;

-- Create snippets table
CREATE TABLE IF NOT EXISTS snippet (
  id INT NOT NULL AUTO_INCREMENT,
  title VARCHAR(50) NOT NULL,
  content TEXT NOT NULL,
  expires DATETIME NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE INDEX idx_snippets_created ON snippets(created_at);