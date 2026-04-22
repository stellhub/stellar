## 使用方式

```bash
npm install
npm run docs:dev
```

默认情况下，VitePress 会以 `docs/` 作为站点根目录，本仓库现有文档会直接作为页面参与构建。

## linux覆盖nginx配置
```shell
sudo cp ./nginx.conf /etc/nginx/nginx.conf
sudo nginx -s reload
```

## linux 定时更新
```shell
chmod +x scripts/auto-pull.sh
crontab -e

* * * * * /bin/bash /path/to/stellar/scripts/auto-pull.sh

```