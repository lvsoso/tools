```shell
npm init
npm install --save express

node app.js

curl localhost:3000
curl localhost:3000 \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"name": "hahah"}'
```