#!/bin/bash
set -e

# MongoDB init scripts run before auth is enabled, so we can connect without credentials
# This script runs automatically when the container starts for the first time

# Create users with proper roles using mongosh
mongosh -- "$MONGO_INITDB_DATABASE" <<EOF
// Create admin user in admin database
db = db.getSiblingDB('admin');

db.createUser({
  user: "$MONGO_USER_ADMIN",
  pwd: "$MONGO_USER_ADMIN_PASSWORD",
  roles: [
    {
      role: "userAdminAnyDatabase",
      db: "admin"
    }
  ]
});

// Switch to application database
db = db.getSiblingDB("$MONGO_INITDB_DATABASE");

db.createUser({
  user: "$MONGO_DB_OWNER",
  pwd: "$MONGO_DB_OWNER_PASSWORD",
  roles: [
    {
      role: "dbOwner",
      db: "$MONGO_INITDB_DATABASE"
    }
  ]
});

db.createUser({
  user: "$MONGO_DB_USER",
  pwd: "$MONGO_DB_USER_PASSWORD",
  roles: [
    {
      role: "readWrite",
      db: "$MONGO_INITDB_DATABASE"
    }
  ]
});

print("MongoDB users created successfully");
EOF
