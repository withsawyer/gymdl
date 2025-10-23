## 开发计划

-   [x] 搭建项目模块化框架
-   [x] 确定配置文件字段,加载配置文件
-   [x] 搭建日志框架 Zap.支持分模块,日志级别控制,可以指定输出模式(控制台,文件,混合)
-   [x] Gin 日志中间件 针对异常 需要输出不同日志级别的 log
-   [x] 所属平台解析器
-   [x] 搭建定时任务框架[gocron](https://github.com/go-co-op/gocron)
-   [x] gocron 的中间件实现(日志,鉴权)
-   [x] 通过Github自动化构建
-   [ ] 通过Github自动检测依赖更新,提交pr,合并pr后自动构建docker镜像
-   [ ] 支持指定多个 cookiecloud-key 多端同步
-   [ ] 初始化 cookiecloud
-   [ ] 研究 cookiecloud 的同步机制以及其存储方式,更新方式
-   [ ] 初始化 webdav
-   [ ] 初始化 AI
-   [ ] 搭建 telegram 框架[telebot](https://github.com/tucnak/telebot)
-   [ ] telebot 的中间件实现(日志,鉴权)
-   [ ] telegram 消息处理器开发
-   [ ] telegram 指令处理器开发
