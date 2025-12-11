#!/bin/bash

# MongoDB备份脚本

set -e

BACKUP_DIR="/backup/neuro-guide"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/$DATE"

echo "开始备份MongoDB数据..."
echo "备份路径: $BACKUP_PATH"

# 创建备份目录
mkdir -p $BACKUP_PATH

# 执行备份
docker exec neuro-guide-mongodb mongodump --out /data/backup/$DATE

# 压缩备份
echo "压缩备份文件..."
tar -czf $BACKUP_PATH.tar.gz -C $BACKUP_DIR $DATE

# 删除未压缩的备份目录
rm -rf $BACKUP_PATH

echo "备份完成: $BACKUP_PATH.tar.gz"

# 删除7天前的备份
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "清理完成"
