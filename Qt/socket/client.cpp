#include "client.h"
#include "ui_client.h"
#include <QtNetwork>
#include <windows.h>
#include <QTextCodec>


Client::Client(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::Client)
{
    ui->setupUi(this);
    tcpSocket = new QTcpSocket(this);
    QString host="192.168.242.159";
    //QString host="10.22.42.230";
    int port=8888;
    tcpSocket->connectToHost(host,port);

    connect(tcpSocket,SIGNAL(connected()),this,SLOT(connected()));
    connect(tcpSocket, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(displayError(QAbstractSocket::SocketError)));
}

Client::~Client()
{
    delete ui;
}

void Client::connected()
{
    ui->lineEdit->setText("连接成功");
    qDebug()<<"-------连接成功-------";
}

void Client::displayError(QAbstractSocket::SocketError)
{
    qDebug() <<tcpSocket->errorString();
}

void Client::on_pushButton_clicked()
{


     QString data=ui->textEdit->toPlainText();
     QByteArray senddata;
     //将字符串转为QByteArray且设置字体为utf8
     //发送json格式字符串
     /*
      *json格式
      * `{
      * "ss":"ss",
      * "aa":"bb"
      * }`
*/
     senddata.append("{");
     senddata.append(data.toUtf8());
     senddata.append(",\"function\":");
     senddata.append("\"initOwner\"");
     senddata.append("}");
     tcpSocket->write(senddata);
}

void Client::displayData(QString args1, QString args2, QString args3, QString args4, QString args5, QString args6)
{
    QString display;
    display.append("\"goods_id\":");//转换为json格式
    display.append("\"");
    display.append(args1);
    display.append("\"");
    display.append(",\"goods_name\":");
    display.append("\"");
    display.append(args2);
    display.append("\"");
    display.append(",\"borntime\":");
    display.append("\"");
    display.append(args3);
    display.append("\"");
    display.append(",\"bornspace\":");
    display.append("\"");
    display.append(args4);
    display.append("\"");
    display.append(",\"manufacturer\":");
    display.append("\"");
    display.append(args5);
    display.append("\"");
    display.append(",\"owner_id\":");
    display.append("\"");
    display.append(args6);
    display.append("\"");

    ui->textEdit->setText(display);
}
void Client::displayData2(QString args1, QString args2, QString args3)
{
    QString display;
    display.append("\"store_id\":");//转换为json格式
    display.append("\"");
    display.append(args1);
    display.append("\"");
    display.append(",\"store_name\":");
    display.append("\"");
    display.append(args2);
    display.append("\"");
    display.append(",\"merchant\":");
    display.append("\"");
    display.append(args3);
    display.append("\"");
    ui->textEdit->setText(display);
}
void Client::readMessage()
{
    QDataStream in(tcpSocket);
    // 设置数据流版本，这里要和服务器端相同
    in.setVersion(QDataStream::Qt_5_7);

    // 如果是刚开始接收数据
    if (blockSize == 0) {
        //判断接收的数据是否大于两字节，也就是文件的大小信息所占的空间
        //如果是则保存到blockSize变量中，否则直接返回，继续接收数据
        if(tcpSocket->bytesAvailable() < (int)sizeof(quint16)) return;
        in >> blockSize;
    }
    // 如果没有得到全部的数据，则返回，继续接收数据
    if(tcpSocket->bytesAvailable() < blockSize) return;
    // 将接收到的数据存放到变量中
    in >> message;
    // 显示接收到的数据
    ui->textEdit_2->setText(message);
}
