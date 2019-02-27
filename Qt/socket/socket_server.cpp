#include "socket_server.h"
#include "ui_socket_server.h"
#include<QTcpServer>
#include<QMessageBox>
socket_server::socket_server(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::socket_server)
{
    ui->setupUi(this);

    tcpserver=new QTcpServer(this);

    connect(tcpserver,SIGNAL(newConnection()),this,SLOT(connected()));

}

socket_server::~socket_server()
{
    delete ui;
}

void socket_server::on_pushButton_2_clicked()
{
    tcpsocket->close();
    this->close();
}


void socket_server::connected()
{
    tcpsocket=tcpserver->nextPendingConnection();
    QMessageBox::about(this,"提示","有新连接！");
    connect(tcpsocket,SIGNAL(readyRead()),this,SLOT(readMessage()));
}

void socket_server::readMessage()
{
    QByteArray arr=tcpsocket->readAll();
    QDataStream dst(arr);
    QString str1;
    QString str2;
    dst>>str1>>str2;
    this->ui->textEdit->setText(str1+str2);

}
void socket_server::sendMessage()
{

}

void socket_server::on_pushButton_clicked()
{
    if(!tcpserver->listen(QHostAddress::LocalHost,7777)){
        qDebug()<<tcpserver->errorString();
        close();
    }


}
