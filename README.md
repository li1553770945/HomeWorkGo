## 基于Gin的班级作业网站管理系统
之前使用Django开发了班级作业网站管理系统，现在仍在使用中，后端请见[HomeWork](https://github.com/li1553770945/HomeWork)
,感觉当初的设计有些不完善的地方，现在打算用Go进行重构。因为第一次接触Go，如果有不对的地方还希望批评指正。

有关项目的详细信息，请见我的[个人博客](http://blog.peacesheep.xyz/categories/?category=%E4%BD%9C%E4%B8%9A%E7%BD%91%E7%AB%99%E5%90%8E%E7%AB%AF%E8%AE%BE%E8%AE%A1)
<!-- more -->
## 主要功能

### 用户管理模块

主要功能如下：

+ 注册
+ 登录
+ 查看我的信息
+ 登出
+ 修改密码功能。

### 小组管理模块

+ 用户可以创建自己的小组
+ 根据小组id和组织密码，加入一个小组
+ 可以根据id查看一个小组的信息
+ 可以查看自己创建的小组


### 作业管理模块

主要功能：
+ 组织创建者可以给自己的小组成员发布作业
+ 组织创建者可以查看作业的完成情况
+ 组织创建者可以一键打包所有的作业
+ 小组成员可以提交作业

### 文件管理模块（可能有）

该模块为专门的文件管理模块，在用户使用中不会感知到。