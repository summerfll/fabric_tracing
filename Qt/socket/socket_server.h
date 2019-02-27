#ifndef SOCKET_SERVER_H
#define SOCKET_SERVER_H

#include <QDialog>
#include<QTcpServer>
#include<QTcpSocket>

namespace Ui {
class socket_server;
}

class socket_server : public QDialog
{
    Q_OBJECT

public:
    explicit socket_server(QWidget *parent = 0);
    ~socket_server();
    QTcpServer *tcpserver;
    QTcpSocket *tcpsocket;

private slots:
    void on_pushButton_2_clicked();
    void sendMessage();//发送信息
    void readMessage();//读取信息
    void connected();//有新连接

    void on_pushButton_clicked();

private:
    Ui::socket_server *ui;
};

#endif // SOCKET_SERVER_H
