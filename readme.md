

linuxosのパッケージマネージャーは aptなど

apt


### dockerコマンド

runはビルド時

cmd はコンテナ起動時

### ビルド
-o はアウトプット

➜  aws_test go build -v -o test/hello main.go
➜  aws_test ./test/hello


### docker 起動

docker run -p 8080:8080 aws_tests

### aws

vpc

パブリックサブネット

ec2を置く

コードを clone


### memo

- vpc
- 


### aws

#### git

sudo yum install -y git



#### golang

curl -LO https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version

### 構成

vpc
サブネット
ルートテーブル
igw
local
SG group
  22
  8080
IGW


### cloud watch

- インストール
sudo yum install -y amazon-cloudwatch-agent

// 最後に
sudo mkdir -p /var/awslogs/etc
sudo touch /var/awslogs/etc/awslogs.conf

// これで設定ファイル作成
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-config-wizard
// 起動
sudo systemctl enable amazon-cloudwatch-agent
sudo systemctl start amazon-cloudwatch-agent
// 起動確認
sudo systemctl status amazon-cloudwatch-agent

// 設定保存先
/opt/aws/amazon-cloudwatch-agent/bin/config.json



- iamロール
ex2のアクセス権限

ポリシー作成
1.CloudWatchAgentServerPolicy
2.CloudWatchLogsFullAccess
3.AmazonEC2ReadOnlyAccess

アタッチ


```
aws ec2 associate-iam-instance-profile --instance-id i-xxxxxxxxxxxxxxxxx --iam-instance-profile Name=YourProfileName

```


- カスタムメトリクス

デフォルト
/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json

ない場合

sudo nano /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json


/opt/aws/amazon-cloudwatch-agent/etc/ 配下

```
{
  "metrics": {
    "append_dimensions": {
      "InstanceId": "${aws:InstanceId}"
    },
    "metrics_collected": {
      "mem": {
        "measurement": [
          "mem_used_percent"
        ],
        "metrics_collection_interval": 60
      }
    }
  }
}
```

- 適応


```
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
    -a fetch-config \
    -m ec2 \
    -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json \
    -s
```

起動

```
sudo systemctl enable amazon-cloudwatch-agent
sudo systemctl start amazon-cloudwatch-agent

// 起動確認
sudo systemctl status amazon-cloudwatch-agent
```

画面確認

AWS マネジメントコンソール → CloudWatch → メトリクス → CWAgent

### 2025/02/02

### go mysql docker-compose

https://zenn.dev/ajapa/articles/443c396a2c5dd1

### マルチステージビルド

docker-compose build は as buildが実行

docker-compose up は fromの as buildではない方が実行される

- devと本番を分けるのが少し難しい

下記の設定に

```
docker-compose build --build-arg TARGET=dev
docker-compose up -d

docker-compose build --build-arg TARGET=prd
docker-compose up -d
```

- 詰まった内容

docker fileで build して mainを実行して、composeでそのまま実行するとうまくいかない
=> compose のcommandで 上書きして main.goを実行させる

### mysqlをローカルで接続する

- 環境変数を利用

### mysql ポイント

- amazon linux 2では mariadb？を使ってて微妙に違うので、mysqlを入れる方が良い
  - clientからの接続時も認証方式が違うとかで接続できない
- api sg からのネットワークだけを受け付けた
- クライアントからのアクセス
 - ユーザーごとにどこからのホストからアクセスするのか？の設定が必要
   - ホストは解放して良いと思う。sgで制限するので



### aws mysql

- nat gw

elastic ip必要で、nat gwがインターネットにアクセスするために必要らしいs

- 環境

amazon linux 2


```
sudo yum update -y

sudo amazon-linux-extras enable mysql8.0

sudo yum install -y mysql mysql-server

// 起動
sudo systemctl start mysqld
sudo systemctl enable mysqld

// 初期パスワード確認
sudo cat /var/log/mysqld.log | grep 'temporary password'

// 初期設定
sudo mysql_secure_installation
・一時パスワード(上記のパスワードかな)
・新しいルートパスワード
・設定進める

// 起動
mysql -u root -p
```

### mysql sg group

- 0.0.0.0/0 からの 3306 の許可
セキュリティグループを指定して良いかも

```
// これだとローカルでもアクセスできてしまいそう
mysql -u testuser -p -h <EC2のパブリックIP>
```

### データベースとユーザー

```
CREATE DATABASE testdb;
CREATE USER 'testuser'@'%' IDENTIFIED BY 'password123';
GRANT ALL PRIVILEGES ON testdb.* TO 'testuser'@'%';
FLUSH PRIVILEGES;
```

### バッションサーバー mysql クライアント


```
sudo yum install -y mariadb

```

### mysql エラー

- sg group でポート開けたのにエラー
原因はmysqlの ホストの制限が localhostのみをdefaultで設定してるみたい

```
[ec2-user@ip-10-0-0-166 ~]$ mysql -u root -h 10.0.1.118 -p
Enter password:
ERROR 1130 (HY000): Host '10.0.0.166' is not allowed to connect to this MySQL server
```

- ホストの確認 && 権限を付与

```
SELECT Host, User FROM mysql.user;


// 新しく付与
CREATE USER 'user'@'%' IDENTIFIED BY 'Str0ngP@ssw0rd!';
GRANT ALL PRIVILEGES ON *.* TO 'user'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

- 初期パスワードは変更

```
ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;

// 再起動
sudo systemctl restart mysqld

```


###　環境変数


```
echo 'export DB_HOST="10.0.1.118"' >> ~/.bashrc
echo 'export DB_USER="user"' >> ~/.bashrc
echo 'export DB_PASSWORD="Str0ngP@ssw0rd!"' >> ~/.bashrc
echo 'export DB_NAME="db"' >> ~/.bashrc
echo 'export DB_PORT="3306"' >> ~/.bashrc

# 設定を反映
source ~/.bashrc
```

### ロードバランサー

- http のポートを80にして詰まった
アプリは8080で開けていたため

- sg group
alb用を作成 カスタムtcp 8080で作成
http, httpsも対応できるようにしたい

api-server-sg インバウンドルールを alg-sgに


## 次

- https
- ドメイン設定
- ecs
- cloud watch


### cloud watch 系

https://qiita.com/hiroaki-u/items/09c7492f68c1fc6c1437

https://www.site24x7.jp/blog/cloudwatch-ec2-monitoring/

https://qiita.com/kaburagi_/items/58364e957b63b981a0bc

https://qiita.com/m_t_u_r_/items/2cf73b6a32c11357bb37

- [] logを吐き出すコード
- aws ネットワーク構築
 - vpc
 - パブリックサブネット
 - igw
 - router パブリック
- cloud watchの設定

#### ログ転送の仕組み

Q. logファイルへの書き込みいらないのか？
A. いらないらしい 

理由は下記
仕組みの確認：CloudWatch Logs に Go のログを送るときに log ファイルは必要か？
結論として、標準出力 (stdout) に直接ログを出力すれば log ファイルは不要 です。CloudWatch Logs エージェント (awslogs) は 標準出力 から直接 CloudWatch Logs に送ることが可能だからです。



### cloud watch 設定

### ポイント
- ロール作成
- log出力ロジック && logファイル作成 && 権限
  - todo ログファイルのローテーション設定
- cloud watch logsインストール && 設定
  - ecsなら stdoutで検知するがec2ならログファイル必要
  - cloud watch のロググループに出力される
  - cpu メモリ
    - cpu は　cpu_usage_idleで検索
     - totalcpu trueにする必要性あるかも


- ロール作成

```
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"logs:CreateLogGroup",
				"logs:CreateLogStream",
				"logs:PutLogEvents"
			],
			"Resource": "arn:aws:logs:ap-northeast-1:682033467141:log-group:/ecs/api-server:*"
		},
		{
			"Effect": "Allow",
			"Action": "ec2:DescribeTags",
			"Resource": "*"
		},
		{
			"Effect": "Allow",
			"Action": [
				"cloudwatch:PutMetricData"
			],
			"Resource": "*"
		}
	]
}
```

- ログファイル作成

権限ないって言われるから

```

sudo touch /var/log/go_app.log
sudo chown ec2-user:ec2-user /var/log/go_app.log
sudo chmod 664 /var/log/go_app.log
```

- cloud watch logsインストール && 設定
1. agent install

```
sudo yum install -y amazon-cloudwatch-agent
```

2. 設定

```
[ec2-user@ip-172-31-7-118 ~]$ sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-config-wizard
================================================================
= Welcome to the Amazon CloudWatch Agent Configuration Manager =
=                                                              =
= CloudWatch Agent allows you to collect metrics and logs from =
= your host and send them to CloudWatch. Additional CloudWatch =
= charges may apply.                                           =
================================================================
On which OS are you planning to use the agent?
1. linux
2. windows
3. darwin
default choice: [1]:
1
Trying to fetch the default region based on ec2 metadata...
I! imds retry client will retry 1 timesAre you using EC2 or On-Premises hosts?
1. EC2
2. On-Premises
default choice: [1]:
1
Which user are you planning to run the agent?
1. cwagent
2. root
3. others
default choice: [1]:
1
Do you want to turn on StatsD daemon?
1. yes
2. no
default choice: [1]:
2
Do you want to monitor metrics from CollectD? WARNING: CollectD must be installed or the Agent will fail to start
1. yes
2. no
default choice: [1]:
2
Do you want to monitor any host metrics? e.g. CPU, memory, etc.
1. yes
2. no
default choice: [1]:
1
Do you want to monitor cpu metrics per core?
1. yes
2. no
default choice: [1]:
2
Do you want to add ec2 dimensions (ImageId, InstanceId, InstanceType, AutoScalingGroupName) into all of your metrics if the info is available?
1. yes
2. no
default choice: [1]:
2
Do you want to aggregate ec2 dimensions (InstanceId)?
1. yes
2. no
default choice: [1]:
2
Would you like to collect your metrics at high resolution (sub-minute resolution)? This enables sub-minute resolution for all metrics, but you can customize for specific metrics in the output json file.
1. 1s
2. 10s
3. 30s
4. 60s
default choice: [4]:
4
Which default metrics config do you want?
1. Basic
2. Standard
3. Advanced
4. None
default choice: [1]:
2
Current config as follows:
{
	"agent": {
		"metrics_collection_interval": 60,
		"run_as_user": "cwagent"
	},
	"metrics": {
		"metrics_collected": {
			"cpu": {
				"measurement": [
					"cpu_usage_idle",
					"cpu_usage_iowait",
					"cpu_usage_user",
					"cpu_usage_system"
				],
				"metrics_collection_interval": 60,
				"totalcpu": false // ⭐️ trueにする必要がるかも
			},
			"disk": {
				"measurement": [
					"used_percent",
					"inodes_free"
				],
				"metrics_collection_interval": 60,
				"resources": [
					"*"
				]
			},
			"diskio": {
				"measurement": [
					"io_time"
				],
				"metrics_collection_interval": 60,
				"resources": [
					"*"
				]
			},
			"mem": {
				"measurement": [
					"mem_used_percent"
				],
				"metrics_collection_interval": 60
			},
			"swap": {
				"measurement": [
					"swap_used_percent"
				],
				"metrics_collection_interval": 60
			}
		}
	}
}
Are you satisfied with the above config? Note: it can be manually customized after the wizard completes to add additional items.
1. yes
2. no
default choice: [1]:
1
Do you have any existing CloudWatch Log Agent (http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AgentReference.html) configuration file to import for migration?
1. yes
2. no
default choice: [2]:
2
Do you want to monitor any log files?
1. yes
2. no
default choice: [1]:
1
Log file path:
/var/log/go_app.log
Log group name:
default choice: [go_app.log]
api-log
Log group class:
1. STANDARD
2. INFREQUENT_ACCESS
default choice: [1]:
1
Log stream name:
default choice: [{instance_id}]

Log Group Retention in days
1. -1
2. 1
3. 3
4. 5
5. 7
6. 14
7. 30
8. 60
9. 90
10. 120
11. 150
12. 180
13. 365
14. 400
15. 545
16. 731
17. 1096
18. 1827
19. 2192
20. 2557
21. 2922
22. 3288
23. 3653
default choice: [1]:
7
Do you want to specify any additional log files to monitor?
1. yes
2. no
default choice: [1]:
2
Do you want the CloudWatch agent to also retrieve X-ray traces?
1. yes
2. no
default choice: [1]:
2
Existing config JSON identified and copied to:  /opt/aws/amazon-cloudwatch-agent/etc/backup-configs
Saved config file to /opt/aws/amazon-cloudwatch-agent/bin/config.json successfully.
Current config as follows:
{
	"agent": {
		"metrics_collection_interval": 60,
		"run_as_user": "cwagent"
	},
	"logs": {
		"logs_collected": {
			"files": {
				"collect_list": [
					{
						"file_path": "/var/log/go_app.log",
						"log_group_class": "STANDARD",
						"log_group_name": "api-log",
						"log_stream_name": "{instance_id}",
						"retention_in_days": 30
					}
				]
			}
		}
	},
	"metrics": {
		"metrics_collected": {
			"cpu": {
				"measurement": [
					"cpu_usage_idle",
					"cpu_usage_iowait",
					"cpu_usage_user",
					"cpu_usage_system"
				],
				"metrics_collection_interval": 60,
				"totalcpu": false
			},
			"disk": {
				"measurement": [
					"used_percent",
					"inodes_free"
				],
				"metrics_collection_interval": 60,
				"resources": [
					"*"
				]
			},
			"diskio": {
				"measurement": [
					"io_time"
				],
				"metrics_collection_interval": 60,
				"resources": [
					"*"
				]
			},
			"mem": {
				"measurement": [
					"mem_used_percent"
				],
				"metrics_collection_interval": 60
			},
			"swap": {
				"measurement": [
					"swap_used_percent"
				],
				"metrics_collection_interval": 60
			}
		}
	}
}
Please check the above content of the config.
The config file is also located at /opt/aws/amazon-cloudwatch-agent/bin/config.json.
Edit it manually if needed.
Do you want to store the config in the SSM parameter store?
1. yes
2. no
default choice: [1]:
2
Program exits now.
```

3. 反映

```
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
  -a fetch-config \
  -m ec2 \
  -c file:/opt/aws/amazon-cloudwatch-agent/bin/config.json \
  -s
```

- 起動

```
sudo systemctl start amazon-cloudwatch-agent

// 自動起動
sudo systemctl enable amazon-cloudwatch-agent

// 起動確認
sudo systemctl status amazon-cloudwatch-agent
```

- go バックグラウンド

```
nohup go run main.go > output.log 2>&1 &


// これだけでバックグラウンド
go run main.go &
```

### デバック

- 
- cloud watch log agentでエラーを吐いていないか？
 sudo tail -f /opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log