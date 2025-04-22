#!/bin/bash
#The author of this file was entirely ChatGPT
# Exit if any command fails
set -e

# DB credentials – change as needed
DB_HOST="yourHost"
DB_PORT="yourPort"
DB_NAME="yourDB"
DB_USER="youruser"
DB_PASS="yourpassword"

# Connect and run SQL commands
echo "🚀 Seeding MySQL database '$DB_NAME' on $DB_HOST..."

mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" <<EOF

-- Optional: Drop existing data
DELETE FROM assignments;
DELETE FROM students;
DELETE FROM rooms;
DELETE FROM buildings;

-- Insert buildings
INSERT INTO buildings (building_id, building_name, has_ac, has_dining) VALUES
(1, 'Maple Hill South', 1, 1),
(2, 'Maple Hill East', 0, 1),
(3, 'Maple Hill West', 1, 0),
(4, 'Founders Hall', 1, 1),
(5, 'Hotz Honors Hall', 1, 1);

-- Insert rooms
$(for b in {1..5}; do
  for r in {101..110}; do
    bedrooms=$(( (RANDOM % 2) + 1 ))
    priv_bath=$(( RANDOM % 2 ))
    has_kitchen=$(( RANDOM % 2 ))
    echo "INSERT INTO rooms (building_id, room_number, num_bedroom, private_bathrooms, has_kitchen) VALUES ($b, $r, $bedrooms, $priv_bath, $has_kitchen);"
  done
done)

-- Insert students
$(for i in {1..100}; do
  ac=$(( RANDOM % 2 ))
  dining=$(( RANDOM % 2 ))
  kitchen=$(( RANDOM % 2 ))
  priv_bath=$(( RANDOM % 2 ))
  echo "INSERT INTO students (wants_ac, wants_dining, wants_kitchen, wants_private_bath) VALUES ($ac, $dining, $kitchen, $priv_bath);"
done)

-- Insert 50 assignments
$(for i in {1..50}; do
  b=$(( (RANDOM % 5) + 1 ))
  r=$(( (RANDOM % 10) + 101 ))
  echo "INSERT INTO assignments (student_id, building_id, room_number) VALUES ($i, $b, $r);"
done)

EOF

echo "✅ Database seeded successfully!"
