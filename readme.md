### skywalking默认不支持dingding告警所以我就写了一个小程序来做代理，这样就可以转发skywalking的告警请求到钉钉

## 使用方法

#### 1.编译main.go 生成可执行文件
#### 2.使用编译好的文件如main 则执行 main -h 查看使用帮助

## 举个栗子

> main -p 8080 -u https://oapi.dingtalk.com/robot/send?access_token=70baf9b1a90da92534c2b6c1x52737da0fd07146135435e0b72d2d0de815f1d4

> p 服务器端口号 

> u 钉钉告警机器人url
