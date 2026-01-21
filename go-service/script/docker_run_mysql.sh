docker run -d \
  --name mysql-local \
  -e MYSQL_ROOT_PASSWORD=qweasd \
  -e MYSQL_DATABASE=nspas_db \
  -p 3306:3306 \
  -v /Users/hedengfeng/workspace/nspas/go-service/storage/mysql:/var/lib/mysql \
  mysql:8.0 \
  --default-authentication-plugin=mysql_native_password