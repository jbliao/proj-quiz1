Quiz1
===
## 使用方式
### Build Image
```
docker build -t proj-quiz1:latest .
```

##### 環境變數說明
- `QUIZ1_APP_PORT`: Server listen port
- `QUIZ1_DB_HOST`: 連線至資料庫的Hostname
- `QUIZ1_DB_NAME`: 資料庫名稱
- `QUIZ1_DB_USER`: 資料庫認證使用者名稱
- `QUIZ1_DB_PASSWD`: 資料庫密碼

### Run with docker command
1. 設定環境變數
```
export DBPASSWD=<資料庫密碼>
export SERVERIP=`hostname -I | cut -d' ' -f 1`
```
2. (Optinal) 建立資料庫
```
docker run -d --name mariadb -p 3306:3306 \
    -v ${PWD}/sql:/docker-entrypoint-initdb.d \
    -e MARIADB_ROOT_PASSWORD=${DBPASSWD} \
    mariadb:latest --init-file /docker-entrypoint-initdb.d/init.sql
```
3. 啟動 server container
```
docker run --rm --name server -p 8080:8080 \
    -e QUIZ1_APP_PORT=8080 \
    -e QUIZ1_DB_NAME=quiz1 \
    -e QUIZ1_DB_USER=root \
    -e QUIZ1_DB_PASSWD=${DBPASSWD} \
    -e QUIZ1_DB_HOST=${SERVERIP} \
    proj-quiz1:latest
```

### Run with docker-compose
```
docker-compose up
```

### Try Server
#### Create comment
```
curl --location "${SERVERIP}:8080/quiz/v1/comment" \
    --header 'Content-Type: application/json' \
    --data '{
        "uuid": "",
        "parentid": "a1205dab-824a-4e3a-bcd2-ed6102e60ae9",
        "comment": "根據中央氣象局地震測報中心地震報告，這起規模...",
        "author": "氣象局網站",
        "update": null,
        "favorite": false
    }'
```
#### Get comment
**Notice:** _Please manually assign the **UUID** gether from response of the post request_
```
curl --location "${SERVERIP}:8080/quiz/v1/comment/<UUID>"
```
#### Update comment
**Notice:** _Please manually assign the **UUID** gether from response of the post request_
```
curl --location --request PUT "${SERVERIP}:8080/quiz/v1/comment/<UUID>" \
    --header 'Content-Type: application/json' \
    --data '{
        "uuid": "",
        "parentid": "a1205dab-824a-4e3a-bcd2-ed6102e60ae9",
        "comment": "根據中央氣象局地震測報中心地震報告，這起規模...",
        "author": "氣象局網站",
        "update": null,
        "favorite": true
    }'
```
#### Delete comment
**Notice:** _Please manually assign the **UUID** gether from response of the post request_
```
curl --location --request DELETE "${SERVERIP}:8080/quiz/v1/comment/<UUID>"
```