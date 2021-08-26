# fssync

此项目应用于大型应用集群配置文件的管理，以此尝试摒弃 ansible 或 puppet 之类复杂机械的运维体验

![fssync _1_.jpg](https://i.loli.net/2021/08/15/yae9p8OYoirzXCs.jpg)

API

```bash
# 同步文件
curl "http://127.0.0.1:8080/sync?app=app1&env=uat"
# 修改文件
curl -X POST -d '{"appName":"app2","envName":"uat","path":"\/wls","filename":"1","operate":"add","datas":{"key1":"value1","key2":"value2"}}' "http://127.0.0.1:8080/edit"
# 获取文件列表
curl "http://127.0.0.1:8080/getfile?app=app1&env=uat1"
# 下载文件
wget http://127.0.0.1:8080/uat/app1/wls/1
```
