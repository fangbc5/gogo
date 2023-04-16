app目录下每个子目录都是一个微服务

每个微服务应该按如下目录结构组织(不需要的可删除)
api                 controller控制器层
async               异步消息
conf                配置文件(只包含配置中心连接,其他配置从配置中心读取)
const               常量
dao                 数据访问层
model               实体类及模型
rpc                 rpc调用
  provider          对外提供接口
  consumer          调用外部接口
service             业务层
router.go           路由文件