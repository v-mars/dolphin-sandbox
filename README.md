# Dolphin-Sandbox
一个支持 Golang、Python、JavaScript、Shell 等语言的沙箱环境

## 项目简介
Dolphin-Sandbox 提供了一种在安全环境中运行不可信代码的简单方式。它适用于多租户场景，允许多个用户提交代码进行执行。代码将在受限的沙箱环境中运行，限制其可访问的资源和系统调用。

## 功能特性
- 支持多种编程语言（如: Golang、Python、JavaScript、Shell）
- 支持MCP调用
- 支持HTTP API调用
- 多租户安全隔离机制
- 基于 Seccomp 的系统调用过滤
- 易于集成到容器化平台

## 技术架构
Dolphin-Sandbox 主要由以下组件构成：
- 核心沙箱引擎：负责代码执行与资源隔离
- 构建脚本：用于生成不同架构下的二进制文件
- 安装脚本：自动安装必要的依赖库
- 运行时服务：提供 HTTP 接口接收任务并执行

## 使用指南
### 系统要求
目前仅支持 Linux 操作系统（专为 Docker 容器设计），需要以下依赖：
- libseccomp
- pkg-config
- gcc
- golang 1.23.0

### 安装与启动步骤
1. 克隆仓库：`git clone https://github.com/v-mars/dolphin-sandbox`
2. 进入项目目录并运行 `./install.sh` 安装依赖
3. 执行 `./build/build_[amd64|arm64].sh` 构建对应架构的沙箱二进制文件
4. 启动服务：`./main`

### 调试说明
如需调试服务，请先通过构建脚本生成沙箱库的二进制文件，然后使用 IDE 按照常规方式调试。

## 注意事项
- 请确保运行环境满足所有依赖项
- 若需扩展支持其他语言或架构，可根据构建脚本进行扩展
- 生产部署建议启用 Seccomp 配置增强安全性