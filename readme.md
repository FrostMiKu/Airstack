# Airstack

> 文件快递柜

简单四步分享文件
1. 上传文件
2. 获得取件码
3. 输入取件码
4. 下载文件

ps：取件码将在取件后作废

## API

```
/api/ping
网页打开，显示pong (

/api/upload
POST 上传文件
返回 {code:200, pwd:取件码}

/api/download/下载码
GET 下载文件，取件码将在取件后作废

```

## TODO

- frontend pages
- config
- limit api access
- logger
- Resource Type
- cache resource