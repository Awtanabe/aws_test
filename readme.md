

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

sudo yum install git



#### golang

curl -LO https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go verion

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

- 環境変数を利用s