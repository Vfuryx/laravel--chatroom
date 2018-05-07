地址 http://swoole.furyx.cn/

##说明

基于[laravel-s](https://github.com/hhxsv5/laravel-s)来进行开发的匿名聊天室

## 要求

| 依赖 | 说明 |
| -------- | -------- |
| [PHP](https://secure.php.net/manual/zh/install.php) | `>= 5.5.9` |
| [Swoole](https://www.swoole.com/) | `>= 1.7.19` `推荐最新的稳定版` `从2.0.12开始不再支持PHP5` |
| Gzip[可选的] | [zlib](https://zlib.net/)，用于压缩HTTP响应，检查本机`libz`是否可用 *ldconfig -p&#124;grep libz* |
| Inotify[可选的] | [inotify](http://pecl.php.net/package/inotify)，用于修改代码后自动Reload Worker进程，检查本机`inotify`是否可用 *php --ri inotify* |


## 安装
```
composer install 
```
## 配置

具体配置可以参考[laravel-s](https://github.com/hhxsv5/laravel-s)

## 运行
> `php artisan laravels {start|stop|restart|reload|publish}`

| 命令 | 说明 |
| --------- | --------- |
| `start` | 启动LaravelS，展示已启动的进程列表 *ps -ef&#124;grep laravels* |
| `stop` | 停止LaravelS |
| `restart` | 重启LaravelS |
| `reload` | 平滑重启所有worker进程，这些worker进程内包含你的业务代码和框架(Laravel/Lumen)代码，不会重启master/manger进程 |
| `publish` | 发布配置文件到你的项目中`config/laravels.php` |



## 启用WebSocket服务器
> WebSocket服务器监听的IP和端口与Http服务器相同。

1.创建WebSocket Handler类，并实现接口`WebsocketHandlerInterface`。
```PHP
namespace App\Services;
use Hhxsv5\LaravelS\Swoole\WebsocketHandlerInterface;
/**
 * @see https://wiki.swoole.com/wiki/page/400.html
 */
class WebsocketService implements WebsocketHandlerInterface
{
    // 声明没有参数的构造函数
    public function __construct()
    {
    }
    public function onOpen(\swoole_websocket_server $server, \swoole_http_request $request)
    {
        \Log::info('New Websocket connection', [$request->fd]);
        $server->push($request->fd, 'Welcome to LaravelS');
        // throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }
    public function onMessage(\swoole_websocket_server $server, \swoole_websocket_frame $frame)
    {
        \Log::info('Received message', [$frame->fd, $frame->data, $frame->opcode, $frame->finish]);
        $server->push($frame->fd, date('Y-m-d H:i:s'));
        // throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }
    public function onClose(\swoole_websocket_server $server, $fd, $reactorId)
    {
        // throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }
}
```

2.更改配置`config/laravels.php`。
```PHP
// ...
'websocket'      => [
    'enable'  => true,
    'handler' => \App\Services\WebsocketService::class,
],
'swoole'         => [
    //...
    // dispatch_mode只能设置为2、4、5，https://wiki.swoole.com/wiki/page/277.html
    'dispatch_mode' => 2,
    //...
],
// ...
```

3.使用`swoole_table`绑定FD与UserId，可选的，[Swoole Table示例](https://github.com/hhxsv5/laravel-s/blob/master/README-CN.md#%E4%BD%BF%E7%94%A8swoole_table)。也可以用其他全局存储服务，例如Redis/Memcached/MySQL，但需要注意多个`Swoole Server`实例时FD可能冲突。

与Nginx配合使用（推荐）
> 参考 [WebSocket代理](http://nginx.org/en/docs/http/websocket.html)

```Nginx
map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
}
upstream laravels {
    server 192.168.0.1:5200 weight=5 max_fails=3 fail_timeout=30s;
    #server 192.168.0.2:5200 weight=3 max_fails=3 fail_timeout=30s;
    #server 192.168.0.3:5200 backup;
}
server {
    listen 80;
    server_name laravels.com;
    root /xxxpath/laravel-s-test/public;
    access_log /yyypath/log/nginx/$server_name.access.log  main;
    autoindex off;
    index index.html index.htm;
    # Nginx处理静态资源(建议开启gzip)，LaravelS处理动态资源。
    location / {
        try_files $uri @laravels;
    }
    # 当请求PHP文件时直接响应404，防止暴露public/*.php
    #location ~* \.php$ {
    #    return 404;
    #}
    # Http和Websocket共存，Nginx通过location区分
    # Javascript: var ws = new WebSocket("ws://laravels.com/ws");
    location =/ws {
        proxy_http_version 1.1;
        # proxy_connect_timeout 60s;
        # proxy_send_timeout 60s;
        # proxy_read_timeout：如果60秒内客户端没有发数据到服务端，那么Nginx会关闭连接；同时，Swoole的心跳设置也会影响连接的关闭
        # proxy_read_timeout 60s;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-PORT $remote_port;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $host;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_pass http://laravels;
    }
    location @laravels {
        proxy_http_version 1.1;
        # proxy_connect_timeout 60s;
        # proxy_send_timeout 60s;
        # proxy_read_timeout 60s;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-PORT $remote_port;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $host;
        proxy_pass http://laravels;
    }
}
```
