#!/bin/sh
echo "Waiting for Postgres..."
until pg_isready -h db -p 5432 -U postgres; do
  sleep 1
done

echo "Postgres is ready. Starting backend..."
/app/ashes-vote-backend
