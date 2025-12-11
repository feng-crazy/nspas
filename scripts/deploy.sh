#!/bin/bash

# 神经科学修行助手部署脚本

set -e

echo "=========================================="
echo "神经科学修行助手部署脚本"
echo "=========================================="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker未安装，请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 检查.env文件
if [ ! -f .env ]; then
    echo "警告: .env文件不存在，正在从.env.example创建..."
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "请编辑.env文件，设置必要的环境变量"
        exit 1
    else
        echo "错误: .env.example文件不存在"
        exit 1
    fi
fi

# 构建前端
echo "构建Web前端..."
cd web-frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi
npm run build
cd ..

# 停止现有容器
echo "停止现有容器..."
docker-compose down

# 构建并启动服务
echo "构建并启动服务..."
docker-compose build
docker-compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务健康状态
echo "检查服务健康状态..."
if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✓ Go服务运行正常"
else
    echo "✗ Go服务启动失败"
    docker-compose logs go-service
    exit 1
fi

if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo "✓ Python AI服务运行正常"
else
    echo "✗ Python AI服务启动失败"
    docker-compose logs python-ai-service
    exit 1
fi

echo "=========================================="
echo "部署完成！"
echo "=========================================="
echo "Go服务: http://localhost:8080"
echo "Python AI服务: http://localhost:8000"
echo "Web前端: http://localhost"
echo ""
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo "=========================================="
