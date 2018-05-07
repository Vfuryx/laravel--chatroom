<?php

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
        $userId = mt_rand(1000, 10000);
        $server->wsTable->set("uid:{$userId}", ['value' => $request->fd]);// 绑定uid到fd的映射
        $server->wsTable->set("fd:{$request->fd}", ['value' => $userId]);// 绑定fd到uid的映射

//         throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }

    public function onMessage(\swoole_websocket_server $server, \swoole_websocket_frame $frame)
    {
        \Log::info('Received message', [$frame->fd, $frame->data, $frame->opcode, $frame->finish]);
        try {

            $data = json_decode($frame->data);
            switch ($data->type) {
                case 'ping'://心跳检测
                    $server->push($frame->fd, json_encode([
                        'stat' => 3,
                        'fd' => $frame->fd,
                        'msg' => $data->type,
                        'time' => time(),
                        'num' => count($server->ports[0]->connections)
                    ], true));
                    break;

                case 'content'://通信
                    foreach ($server->ports[0]->connections as $fd) {
                        $server->push($fd, json_encode([
                            'stat' => 1,
                            'nick' => $server->infoTable->get("fd:{$frame->fd}")['nick'],
                            'fd' => $frame->fd,
                            'msg' => $data->msg,
                            'time' => time(),
                            'num' => count($server->ports[0]->connections)
                        ], true));
                    }

                    break;

                case 'login'://登记nick
                    $server->infoTable->set("fd:{$frame->fd}", ['nick' => $data->nick]);//绑定用户信息
                    $members = [];
                    foreach ($server->ports[0]->connections as $fd) {
                        if ($fd == $frame->fd)
                            continue;
                        $members[$fd] = $server->infoTable->get("fd:{$fd}")['nick'];
                        $server->push($fd, json_encode([
                            'stat' => 4,//添加新用户
                            'nick' => $data->nick,
                            'fd' => $frame->fd,
                            'time' => time(),
                            'num' => count($server->ports[0]->connections)
                        ], true));
                    }
                    $server->push($frame->fd, json_encode(['stat' => 2, 'fd' => $frame->fd, 'members' => $members, 'num' => count($server->ports[0]->connections)], true));
                    break;
            }
        } catch (\Error $error) {
            print_r($error);
        }


        // throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }

    public function onClose(\swoole_websocket_server $server, $fd, $reactorId)
    {
        $uid = app('swoole')->wsTable->get('fd:' . $fd);
        if ($uid !== false) {
            $server->wsTable->del('uid:' . $uid['value']);// 解绑uid映射
        }
        $server->wsTable->del('fd:' . $fd);// 解绑fd映射
        $server->infoTable->del('fd:' . $fd);// 去除信息

        foreach ($server->ports[0]->connections as $cfd) {//发送登出信息
            if ($cfd == $fd)
                continue;
            $server->push($cfd, json_encode([
                'stat' => 5,//删除用户
                'fd' => $fd,
                'time' => time(),
                'num' => count($server->ports[0]->connections)
            ], true));
        }

        // throw new \Exception('an exception');// 此时抛出的异常上层会忽略，并记录到Swoole日志，需要开发者try/catch捕获处理
    }


    public function log($server)
    {
        if ($fp = fopen(public_path('log.log'), 'w')) {
            if (flock($fp, LOCK_EX)) {
                $str = '';
                foreach ($server->wsTable as $key => $value) {
                    $str .= "$key:" . json_encode($value) . ';';
                }
                fwrite($fp, $str);
                flock($fp, LOCK_UN);
            }
            fclose($fp);
        }
    }
}