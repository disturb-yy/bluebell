conf目录
使用viper读取config.yaml配置参数，并将其反序列化到conf变量

pkg目录
存放引用的第三方库

router
请求到达以后，由router根据访问的URL做路由转发，交给controller层

controller 
- 请求参数获取和校验 —— 需要反序列数据保存在struct中，
由于该struct在controller和logic都要用到，所以可以将其抽象
出来放在models
- 业务处理 —— 交由logic层处理
- response 和 code —— 封装状态码和状态信息

logic
具体的业务操作，如用户注册逻辑: 
- 判断用户存不存在 —— 数据库的查询
- 雪花算法id生成器 —— 用户Uid的生成
- 保存注册的用户信息 —— 数据库的插入

dao
处理数据库相关的操作
- mysql
- redis

models
- params: 保存参数相关的数据