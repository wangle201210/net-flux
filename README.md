# net-flux

MediaBridge 项目的基础库，提供网络、日志、文件等基本操作，以及服务发现、数据上报等基础能力。

## 功能概览

| 模块 | 说明 |
|------|------|
| 网络 | 网络通信与协议封装 |
| 日志 | 统一日志接口与输出 |
| 文件 | 文件读写与路径等基础操作 |
| 服务发现 | 服务注册、注销与查询 |
| 数据上报 | 机器指标、流状态等数据上报 |
| 配置分发 | 配置下发与同步 |
| 事件 / 告警 | 事件与告警上报 |
| 控制 | 远程控制指令 |

协议层通过 Protobuf 定义一级/二级命令（`CMD` / `SCMD*`），涵盖系统心跳、服务发现、数据上报等场景。

## 安装

```bash
go get github.com/dellinger2023/net-flux
```

## 项目结构

```
net-flux/
├── proto/          # Protobuf 协议定义
├── gen/            # 由 proto 生成的 Go 代码
├── examples/       # 使用示例
└── generate.go     # 代码生成入口
```

## 生成代码

依赖 [protoc](https://github.com/protocolbuffers/protobuf) 与 Go 工具链。在项目根目录执行：

```bash
go generate ./...
```

生成的代码输出至 `gen/` 目录。

## 协议命令

一级命令（`CMD`）包括：

- `SYSTEM` — 系统相关（如心跳 Ping/Pong）
- `DISCOVERY` — 服务发现（注册、注销、查询）
- `DATA_REPORT` — 数据上报（机器指标、流增删改查等）
- `CONFIG` — 配置分发
- `EVENT` / `ALARM` — 事件与告警上报
- `CONTROL` — 控制指令

详见 [`proto/base.proto`](proto/base.proto) 与 [`proto/system.proto`](proto/system.proto)。

## 许可证

[MIT License](LICENSE)
