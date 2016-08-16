# Keepalive Server

## Intro

I want to know how to use Keepalive in Golang.

So I Use Socket try to Do it.

## Keepalive

When i quit `WeChat` in Windows, My Phone still show me that `Log in Windows WeChat`.

It is a Problem in TCP Socket, We can solve it By:

1. SetReadDeadline() in net.Conn when your communication protocol support Heartbeat.
2. SetKeepAlive() in net.TCPConn when your communication protocol not support Heartbeat.