docker run -d \
  --name mongodb7 \
  --restart always \
  -p 27017:27017 \
  -v /Users/hedengfeng/workspace/nspas/go-service/storage/mongodb:/data/db \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=123456 \
  mongo:7.0