#include "socket_client.h"
#include "ui_socket_client.h"
#include <QtNetwork>
#include <QFileDialog>

socket_client::socket_client(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::socket_client)
{
    ui->setupUi(this);

    payloadSize = 64*1024; // 64KB
    totalBytes = 0;
    bytesWritten = 0;
    bytesToWrite = 0;
    tcpsocket_client = new QTcpSocket(this);

    // 当连接服务器成功时，发出connected()信号，开始传送文件
    connect(tcpsocket_client, SIGNAL(connected()), this, SLOT(startTransfer()));
    connect(tcpsocket_client, SIGNAL(bytesWritten(qint64)),
            this, SLOT(updatesocket_clientProgress(qint64)));
    connect(tcpsocket_client, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(displayError(QAbstractSocket::SocketError)));
    ui->sendButton->setEnabled(false);
}

socket_client::~socket_client()
{
    delete ui;
}

void socket_client::openFile()
{
    fileName = QFileDialog::getOpenFileName(this);
    if (!fileName.isEmpty()) {
        ui->sendButton->setEnabled(true);
        ui->socket_clientStatusLabel->setText(tr("打开文件 %1 成功！").arg(fileName));
    }
}

void socket_client::send()
{
    ui->sendButton->setEnabled(false);

    // 初始化已发送字节为0
    bytesWritten = 0;
    ui->socket_clientStatusLabel->setText(tr("连接中…"));
    tcpsocket_client->connectToHost(ui->hostLineEdit->text(),
                             ui->portLineEdit->text().toInt());
}


void socket_client::startTransfer()
{
    localFile = new QFile(fileName);
    if (!localFile->open(QFile::ReadOnly)) {
        qDebug() << "socket_client: open file error!";
        return;
    }
    // 获取文件大小
    totalBytes = localFile->size();

    QDataStream sendOut(&outBlock, QIODevice::WriteOnly);
    sendOut.setVersion(QDataStream::Qt_4_0);
    QString currentFileName = fileName.right(fileName.size()
                                             - fileName.lastIndexOf('/')-1);
    // 保留总大小信息空间、文件名大小信息空间，然后输入文件名
    sendOut << qint64(0) << qint64(0) << currentFileName;

    // 这里的总大小是总大小信息、文件名大小信息、文件名和实际文件大小的总和
    totalBytes += outBlock.size();
    sendOut.device()->seek(0);

    // 返回outBolock的开始，用实际的大小信息代替两个qint64(0)空间
    sendOut << totalBytes << qint64((outBlock.size() - sizeof(qint64)*2));

    // 发送完文件头结构后剩余数据的大小
    bytesToWrite = totalBytes - tcpsocket_client->write(outBlock);

    ui->socket_clientStatusLabel->setText(tr("已连接"));
    outBlock.resize(0);
}

void socket_client::updatesocket_clientProgress(qint64 numBytes)
{
    // 已经发送数据的大小
    bytesWritten += (int)numBytes;

    // 如果已经发送了数据
    if (bytesToWrite > 0) {
        // 每次发送payloadSize大小的数据，这里设置为64KB，如果剩余的数据不足64KB，
        // 就发送剩余数据的大小
        outBlock = localFile->read(qMin(bytesToWrite, payloadSize));

        // 发送完一次数据后还剩余数据的大小
        bytesToWrite -= (int)tcpsocket_client->write(outBlock);

        // 清空发送缓冲区
        outBlock.resize(0);
    } else { // 如果没有发送任何数据，则关闭文件
        localFile->close();
    }
    // 更新进度条
    ui->socket_clientProgressBar->setMaximum(totalBytes);
    ui->socket_clientProgressBar->setValue(bytesWritten);
    // 如果发送完毕
    if(bytesWritten == totalBytes) {
        ui->socket_clientStatusLabel->setText(tr("传送文件 %1 成功").arg(fileName));
        localFile->close();
        tcpsocket_client->close();
    }
}

void socket_client::displayError(QAbstractSocket::SocketError)
{
    qDebug() << tcpsocket_client->errorString();
    tcpsocket_client->close();
    ui->socket_clientProgressBar->reset();
    ui->socket_clientStatusLabel->setText(tr("客户端就绪"));
    ui->sendButton->setEnabled(true);
}


// 打开按钮
void socket_client::on_openButton_clicked()
{
    ui->socket_clientProgressBar->reset();
    ui->socket_clientStatusLabel->setText(tr("状态：等待打开文件！"));
    openFile();

}

// 发送按钮
void socket_client::on_sendButton_clicked()
{
    send();
}
