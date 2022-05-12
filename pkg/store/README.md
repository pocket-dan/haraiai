# store

データベースとして Cloud Firestore を使用している。
1グループのデータが1KiBで収まると仮定すると、無料枠でおよそ

- 合計10万グループ
- 毎日2万件の割り勘リクエスト (1 read + 1 write)

が可能。



## Reference

- [Cloud Firestore  |  Firebase Documentation](https://firebase.google.com/docs/firestore)
- [Firebase Pricing](https://firebase.google.com/pricing)
