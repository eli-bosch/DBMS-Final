#!/bin/bash
set -e

DB_HOST="yourHost"
DB_PORT="yourPort"
DB_NAME="yourDB"
DB_USER="youruser"
DB_PASS="yourpassword"

# array of names to cycle through
names=(Alice Bob Carol Dave Eve Frank Grace Hugo Ivy Jack)

echo "🚀 Seeding MySQL database '$DB_NAME' on $DB_HOST..."

mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" <<EOF

DELETE FROM assignments;
DELETE FROM students;
DELETE FROM rooms;
DELETE FROM buildings;

INSERT INTO buildings (building_id, building_name, has_ac, has_dining) VALUES
  (1, 'Maple Hill South', 1, 1),
  (2, 'Maple Hill East', 0, 1),
  (3, 'Maple Hill West', 1, 0),
  (4, 'Founders Hall', 1, 1),
  (5, 'Hotz Honors Hall', 1, 1);

$(for b in {1..5}; do
  for r in {101..110}; do
    bedrooms=$(( (RANDOM % 2) + 1 ))
    priv_bath=$(( RANDOM % 2 ))
    has_kitchen=$(( RANDOM % 2 ))
    printf "INSERT INTO rooms (building_id, room_number, num_bedroom, private_bathrooms, has_kitchen) VALUES (%d, %d, %d, %d, %d);\n" \
      "$b" "$r" "$bedrooms" "$priv_bath" "$has_kitchen"
  done
done)

$(for i in {1..100}; do
  name="${names[$(( (i-1) % ${#names[@]} ))]}"
  ac=$(( RANDOM % 2 ))
  dining=$(( RANDOM % 2 ))
  kitchen=$(( RANDOM % 2 ))
  priv_bath=$(( RANDOM % 2 ))
  printf "INSERT INTO students (name, wants_ac, wants_dining, wants_kitchen, wants_private_bath) VALUES ('%s', %d, %d, %d, %d);\n" \
    "$name" "$ac" "$dining" "$kitchen" "$priv_bath"
done)

$(for i in {1..50}; do
  b=$(( (RANDOM % 5) + 1 ))
  r=$(( (RANDOM % 10) + 101 ))
  printf "INSERT INTO assignments (student_id, building_id, room_number) VALUES (%d, %d, %d);\n" \
    "$i" "$b" "$r"
done)

EOF

echo "✅ Database seeded successfully!"
