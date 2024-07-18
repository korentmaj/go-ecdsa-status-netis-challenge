#!/bin/bash

psql -U postgres -d your_database -f /path/to/project-root/internal/database/migrations/001_create_status_table.sql
