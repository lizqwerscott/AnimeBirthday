#!/bin/bash

# Usage: ./ru.sh dev|prod

if [ $# -ne 1 ]; then
    echo "Usage: $0 dev|prod"
    exit 1
fi

# 获取传入的参数
env_type=$1

# 映射 dev 到 development，prod 到 production
case "$env_type" in
    "dev")
        env_type="development"
        ;;
    "prod")
        env_type="production"
        ;;
    *)
        echo "Invalid environment type: $env_type. Use 'dev' or 'prod'."
        exit 1
        ;;
esac

# 设置环境变量
export APP_ENV=$env_type

echo "APP_ENV set to: $APP_ENV"

go run main.go

export APP_ENV=""
