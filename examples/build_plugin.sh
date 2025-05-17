#!/bin/bash

# 编译示例插件的脚本

echo "正在编译示例插件..."

# 确保我们在正确的目录中
cd "$(dirname "$0")"

# 编译示例插件
go build -o sample_plugin.so -buildmode=plugin sample_plugin.go

if [ $? -eq 0 ]; then
    echo "编译成功！插件已保存为 sample_plugin.so"
    echo "使用方法："
    echo "1. 启动 Luna"
    echo "2. 在 Luna 中执行: load $(pwd)/sample_plugin.so"
    echo "3. 执行: list 查看已加载的插件"
    echo "4. 执行: exec sample_plugin your_target"
    echo "   或者:"
    echo "   执行: use sample_plugin"
    echo "   执行: set target your_target"
    echo "   执行: run"
else
    echo "编译失败，请检查错误信息"
fi