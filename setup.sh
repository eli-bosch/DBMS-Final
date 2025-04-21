#!/bin/bash
set -e -v

echo "Creating and populating 'messages' table..."

mysql -u casonp -p <<EOFMYSQL

USE casonp;

-- Drop if exists 
DROP TABLE IF EXISTS messages;

-- Create a simple table
CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content VARCHAR(255)
) ENGINE=InnoDB;

-- Insert a test row
INSERT INTO messages (content) VALUES ('Hello from the database!');

EOFMYSQL

echo "✅ messages table created and populated."
