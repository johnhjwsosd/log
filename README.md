Log
=======

日志系统，暂提供mongo驱动

Installation
----

<pre><code>
go get github.com/johnhjwsosd/log
</code></pre>



Example
----------
通过注册日志对象获取ID
---------
stType : 驱动类型
stHost : 目标地址
stPort : 目标端口
stName : 存储库名称
appName : 应用名称对应collection
msName : 服务名称
<pre><code>
    http://192.168.1.183:8081/reg?stType=mongo&stHost=127.0.0.1&stPort=27017&stName=test&appName=testapp&msName=testms
</code></pre>

写日志
---------
appID  : 日志对象（通过注册获取）
content : 日志内容（key:value）
level  : 日志级别(info,trace,warn,error,fatal)
xxx,temp : 日志内容 对应content
<pre><code>
  http://192.168.1.183:8081/wl/?appID=cb7996eabaad966395b2f4a9e4188a81&level=info&xxx=223&temp=test1
</code></pre>


读日志
----------
appID : 日志对象（通过注册获取）
where : 查询条件（key:value）
xxx,temp : 具体查询条件
<pre><code>
  http://192.168.1.183:8081/rl/?appID=cb7996eabaad966395b2f4a9e4188a81&xxx=223&temp=test1
</code></pre>

