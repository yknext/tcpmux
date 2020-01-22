# tcpmux

tcp 多路复用转发

server -local 0.0.0.0:80 -forward 127.0.0.1:443

server -local 0.0.0.0:8080 -forward 127.0.0.1:80
