#ifndef ADD_STORE_H
#define ADD_STORE_H

#include <QWidget>

namespace Ui {
class add_store;
}

class add_store : public QWidget
{
    Q_OBJECT

public:
    explicit add_store(QWidget *parent = 0);
    ~add_store();
signals:

    void sendData(QString,QString,QString);//向socket客户端发送信息
private slots:
    void on_pushButton_3_clicked();

    void on_pushButton_2_clicked();

    void on_pushButton_clicked();
    void displayphone_num(QString staff_name);

private:
    Ui::add_store *ui;
};

#endif // ADD_STORE_H
