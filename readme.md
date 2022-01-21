# 一个极简聊天室

目前是HIML+Javascript 实现前端的Websocket连接，相应的Android实现在开发中。


## 运行
1. go run main.go
2. 访问localhost:8081/login，直接点击login
3. 访问localhost:8080/ 发送消息
## 已实现的功能
1. 发送消息并推送给全体在线成员
2. Login之后才能说话
3. 登录后才能看到消息。
1. 发送的每一条消息根据jwt设置发送方
2. 给用户添加昵称 *考虑在jwt里添加昵称项*
1. 用户离开*关闭websocket连接*后，在client中删除用户信息，
3. 在数据库而不是程序中保存用户信息
1. 真正的检查用户是否为合法用户，
1. 注册 *向 /sign post 一个json数据,请求示例`{"username":"teste3","password":"test2"}`* 
## 将要开发的功能
1. 聊天记录存储（用于获取历史消息）
5. 私聊
2. 密码在后端也加密处理
3. 设置jwt 的Redis保存，实现刷新jwt和鉴权jwt
4. 添加修改昵称的功能
5. 实现方便的聊天命令开发。