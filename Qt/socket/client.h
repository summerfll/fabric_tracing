#ifndef CLIENT_H
#define CLIENT_H

#include <QDialog>
#include <QAbstractSocket>
class QTcpSocket;

namespace Ui {
class Client;
}

class Client : public QDialog
{
    Q_OBJECT

public:
    explicit Client(QWidget *parent = 0);
    ~Client();

private:
    Ui::Client *ui;
    QTcpSocket *tcpSocket;
    QString message;
    quint16 blockSize;
signals:
    void senddata(QByteArray&);
    //void sendToClient(QString,QString,QString,QString,QString,QString);
private slots:
    void connected();//连接成功
    void readMessage();//接收数据
    void displayError(QAbstractSocket::SocketError);//socket错误
    void on_pushButton_clicked();
    void displayData(QString,QString,QString,QString,QString,QString);//展示商品初始化信息
    void displayData2(QString,QString,QString);//展示仓储初始化信息

};

#endif // CLIENT_H
