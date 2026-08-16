# 青岚校报文章工作台

一个自包含的校报编务工作台。编辑可按新闻、访谈、社论和活动栏目浏览稿件，在同一界面完成正文编辑、预览、保存和状态流转，并查看完成清单。

项目采用单一 Go module，使用固定内存 fixture，不依赖数据库、网络服务、系统时间或随机数据。HTTP 路由使用固定版本的 `github.com/go-chi/chi/v5`，前端为 `web/` 下的原生 HTML、CSS 和 JavaScript 三件套。

## 环境

- Go 1.23.12
- Node.js 20
- `CGO_ENABLED=0`

## 运行

```bash
CGO_ENABLED=0 go run .
```

打开 `http://localhost:8080`。可通过 `PORT` 环境变量指定其他端口。

## 构建

```bash
CGO_ENABLED=0 go build ./...
cd web && npm ci && npm run build
```

前端构建会生成 `web/dist/`，该目录仅用于本地构建验证，不提交到源树。

## 测试

```bash
CGO_ENABLED=0 go test -count=1 ./...
```
