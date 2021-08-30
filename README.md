# fssync

此项目应用于大型应用集群配置文件的管理，以此尝试摒弃 ansible 或 puppet 之类复杂机械的运维体验

![fssync _1_.jpg](https://i.loli.net/2021/08/15/yae9p8OYoirzXCs.jpg)

API

```bash
# 同步文件
curl "http://127.0.0.1:8080/sync?mod=env&app=app1&zone=uat"
# 修改文件
curl -X POST -d '{"appName":"app2","envName":"uat","path":"\/wls","filename":"1","operate":"add","datas":{"key1":"value1","key2":"value2"}}' "http://127.0.0.1:8080/modify"
curl -X POST -d '{"appName":"app1","envName":"uat","path":"\/wls","filename":"1","operate":"edit","datas":{"key1":"value100","key2":""}}' "http://127.0.0.1:8080/modify"
curl -X POST -d '{"appName":"app1","envName":"uat","path":"\/wls","filename":"1","operate":"del","datas":{"key1":"value100","key2":""}}' "http://127.0.0.1:8080/modify"
# 获取文件列表
curl "http://127.0.0.1:8080/getfile?app=app1&zone=uat"
# 下载文件
wget http://127.0.0.1:8080/uat/app1/wls/1
```
