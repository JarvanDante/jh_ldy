# jh_ldy

JH 推广落地页。结构和 `my_ldy` 一样：Go 认渠道、选包、换肤、拼 iOS 描述文件；浏览器写 Cookie 和剪贴板。

投放链接：

```
https://jh-ldy.panalow.cc/?source=JHAAA
https://jh-ldy.panalow.cc/s/JHAAA
```

渠道码需要和后台渠道管理里的码一致。点下载时剪贴板写入 `channel://码`，iOS 描述文件打开 `https://jh-h5.panalow.cc/?source=码`。

## 路由

| 路径 | 说明 |
|---|---|
| `GET /` | 落地页；`?source=` 会种 Cookie `ldy_source`（5 天） |
| `GET /s/{码}` | 路径带渠道 |
| `GET /api/boot` | 当前渠道的文案、皮肤、安卓/iOS 地址 |
| `GET /download/android.apk` | 读取 `downloads/app-release.apk` |
| `GET /download/ios.mobileconfig` | 现场拼 iOS 描述文件 |
| `GET /health` | 探活 |

## 配置

`config.json`。环境变量 `JH_LDY_ROOT`、`JH_LDY_CONFIG`。

默认 H5 是 `https://jh-h5.panalow.cc/`。`channels` 的 key 是渠道码，有值才覆盖默认包和皮肤。示例渠道 `JHAAA` 使用 `light` 皮肤。

安卓包不进 git。打好后放到 `downloads/app-release.apk`：

```bash
make apk
```

## 本地

```bash
make test
make run
```

页面：`http://127.0.0.1:5175/?source=JHAAA`

## 集群

镜像给 Docker Desktop 的 Kubernetes 用。`jh-test` 1 个副本，`jh-prod` 2 个副本。`jh-ldy.panalow.cc` 由 ingress-nginx 转到生产 Service `jh-ldy:80`。

```bash
make docker-build
kubectl apply -f k8s/deploy.yaml
```
