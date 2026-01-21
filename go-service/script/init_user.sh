#!/bin/bash

curl -X POST http://localhost:8080/api/auth/register \
-H "Content-Type: application/json" \
-d '{
    "email": "admin@admin.com",
    "password": "qweasd",
    "phone": "13800138000"
}'