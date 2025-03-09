#!/bin/bash

# 获取当前时间作为提交信息的一部分
current_time=$(date "+%Y-%m-%d %H:%M:%S")

# 获取当前分支名
current_branch=$(git branch --show-current)

# 添加所有更改的文件
echo "正在添加更改的文件..."
git add .

# 获取更改的文件列表
changed_files=$(git diff --cached --name-only)

if [ -z "$changed_files" ]; then
    echo "没有需要提交的更改。"
    exit 0
fi

# 生成提交信息
commit_message="自动提交: $current_time"
echo "提交信息: $commit_message"
echo "更改的文件:"
echo "$changed_files"

# 提交代码
git commit -m "$commit_message"

# 推送到远程仓库
echo "正在推送到远程仓库的 $current_branch 分支..."
git push origin $current_branch

echo "完成！" 