# msql-monitor
使用:

1.下载main.go

2.使用go build main.go

使用示例:
Usage of ./st: 
  -h string 
        mysql host (default "127.0.0.1") 
  -p string 
        mysql user password 
  -port string 
        mysql port (default "3306") 
  -u string 
        mysql user 
展示样例:     

[root@txctest scripts]# ./st -u txc -p 123456 -port 3309

----------|------------|----- Innodb row operation ----|---Buffer Pool---|--lock--|----------|----Thread---
---Time---|--QPS--TPS--|  read inserted updated deleted|   Hit    dirty  |  lock  | slow logs|  conn   run

 11:27:51 |   0     0  |      0       0       0       0| 96.76      0.00 |     0  |      0   |     1     2 

 11:28:01 |   2     0  |      0       0       0       0| 96.76      0.00 |     0  |      0   |     2     2 

 11:28:11 |8036     397|167536     398     798      399| 96.76     27.69 |     0  |      0   |    13    10 

 11:28:21 |68113   3404|1419977    3403    6805    3402| 96.76     56.09 |     0  |      0   |    14     8 

 11:28:31 |70901   3546|1477381    3546    7091    3546| 96.76     56.37 |     0  |      0   |    15     9 
 
 QPS :=Questions 两次时间的差值
 
 TPS :=Com_select 两次时间的差值-Com_commit 两次时间的差值
 
 lock =Innodb_row_lock_current_waits 状态的当前值
