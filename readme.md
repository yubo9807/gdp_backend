# 私有包管理服务

## 接口服务

### 1. 打包
```shell
./scripts/build.sh  # linux，其他系统请自行修改参数
```

### 2. 文件上传至服务器

```bash
mkdir config_service.yml
```

### 3. 配置参数
```yaml
prefix: "/api"
port: 9800  # 启动端口

logDir: "./logs"  # 日志目录
logReserveTime: 10  # 日志保留时间(d)

dbDir: "./db"  # 存数据目录

modulesDir: "./modules"  # 依赖包目录
allowedCatalogs: ["modules"]  # 允许接口访问的目录
```

### 4. 启动服务

```bash
./server &
```

## 前端

前往前端仓库： https://github.com/yubo9807/gdp_frontend

## 命令行工具

### 1. 打包

```shell
./scripts/build.sh
```

将文件放到 `~/gdp/` 下，并配置环境变量。日志、配置文件都会在该目录下。
