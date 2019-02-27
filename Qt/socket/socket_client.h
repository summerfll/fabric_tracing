#ifndef socket_client_H
#define socket_client_H

#include <QDialog>
#include <QAbstractSocket>
class QTcpSocket;
class QFile;

namespace Ui {
class socket_client;
}

class socket_client : public QDialog
{
    Q_OBJECT

public:
    explicit socket_client(QWidget *parent = 0);
    ~socket_client();

private:
    Ui::socket_client *ui;

    QTcpSocket *tcpsocket_client;
    QFile *localFile;     // 要发送的文件
    qint64 totalBytes;    // 发送数据的总大小
    qint64 bytesWritten;  // 已经发送数据大小
    qint64 bytesToWrite;  // 剩余数据大小
    qint64 payloadSize;   // 每次发送数据的大小
    QString fileName;     // 保存文件路径
    QByteArray outBlock;  // 数据缓冲区，即存放每次要发送的数据块

private slots:
    void openFile();
    void send();
    void startTransfer();
    void updatesocket_clientProgress(qint64);
    void displayError(QAbstractSocket::SocketError);


    void on_openButton_clicked();
    void on_sendButton_clicked();
};

#endif // socket_client_H
