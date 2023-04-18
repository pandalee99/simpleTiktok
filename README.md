# SimpleTiktok
第五届字节跳动青训营项目

项目方案：https://bytedance.feishu.cn/docs/doccnKrCsU5Iac6eftnFBdsXTof#

实现了大部分的接口，只做了单机版本的

能跑通所有功能，接口测试完备


实现注册、登录、获取视频流、关注列表和留言列表等功能，对每个视频的流的处理
完善，用户也能够上传短视频，并截取视频的画面作为图片同步上传，接口测试完备。
<br/>
<br/>
<br/>
1、使用 MySQL 作为关系型数据库，获取维护数据表的关系，使用 Gorm 进行连接。
<br/>
2、基于 Redis 高性能这一特点，在项目开发过程中使用 Redis 去提高关注和留言的模块的性能，对于一些关注消
息和留言内容，会先放入 Redis 的 Set 列表，提高访问速度。
<br/>
3、对于上传的视频流会使用 FFmpeg 去获取截图，并放入数据库。
<br/>
4、使用 Gin 作为开发框架，使用 Gorm 去访问数据库。
<br/>
5、采用 JWT 作为认证方式
<br/>
<br/>
测试图片：

![test](https://github.com/pandalee99/image_store/blob/master/TikTok/%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%20%E6%8A%96%E9%9F%B3%E5%AE%A2%E6%88%B7%E7%AB%AF.png?raw=true)
