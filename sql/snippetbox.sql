-- Create database
CREATE DATABASE IF NOT EXISTS snippetbox;

-- Switch to 'snippetbox' database
USE snippetbox;

-- Create snippets table
CREATE TABLE IF NOT EXISTS snippet (
  id INT NOT NULL AUTO_INCREMENT,
  title VARCHAR(50) NOT NULL,
  content TEXT NOT NULL,
  expires DATETIME NOT NULL,
  user INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user) REFERENCES users (users.id)
);

CREATE INDEX idx_snippets_created ON snippet(created_at);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(1024) NOT NULL,
  active BOOLEAN NOT NULL DEFAULT false,
  activeToken VARCHAR(1024),
  activeTokenEXPIRE DATETIME,
  passwordResetToken VARCHAR(1024),
  passwordResetTokenExpire DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);