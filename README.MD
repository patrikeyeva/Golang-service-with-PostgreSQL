# CURL запросы

```bash
curl -X POST localhost:9000/article -d '{"name":"Сats","rating":5}' -i

curl -X PUT localhost:9000/article -d '{"id":1,"name":"New Cats","rating":4}' -i

curl -X POST localhost:9000/comment -d '{"article_id":1,"text":"Cute"}' -i
curl -X POST localhost:9000/comment -d '{"article_id":1,"text":"Funny"}' -i

curl -X GET localhost:9000/article?id=1 -i

curl -X DELETE localhost:9000/article?id=1 -i
```