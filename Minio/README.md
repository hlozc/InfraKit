# Minio 的连接和使用

oss(ojbect storage service) 对象存储服务，会将对象作为数据存储在分布式存储集群中，每个对象都有唯一的标识符（URL），可以通过这个标识符来访问和检索数据

对象是对象存储的基本单元，可以理解成任意类型的数据。存储桶是对象的载体，每个桶可以存储任意数量的对象
每个对象都由对象键(object key)、对象值(object value)、对象元数据(Metadata) 
* object key: 对象于对象桶中的唯一标识，可以理解为路径
* object value: 对象本身，即文件内容
* metadata: 一组键值对，存储了对象的元数据，可以理解为对象的属性

## Docker 安装 Minio 和 启动 Minio 容器
```shell
# install
docker pull minio/minio

# run
docker run --name minio \
-p 9000:9000 \
-p 9090:9090 \
-d --restart=always \
-e "MINIO_ROOT_USER=minio" \
-e "MINIO_ROOT_PASSWORD=minio123456Public" \
-v /home/minio/data:/data \
-v /home/minio/config:/root/.minio \
minio/minio server /data --console-address ":9090"
```

-p 9000:9000 表示 API 端口
-p 9090:9090 表示控制台端口
-v /home/minio/data:/data 表示数据目录
-v /home/minio/config:/root/.minio 表示配置文件目录

创建完桶之后，就可以通过
IP:9000/桶名称/文件名进行访问

## Golang 的连接
```
go get github.com/minio/minio-go
```
