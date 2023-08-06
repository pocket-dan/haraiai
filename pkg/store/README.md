# store

データベースとして Cloud Firestore を使用している。
- [Cloud Firestore](https://firebase.google.com/docs/firestore?hl=ja)
- [Cloud Firestore の課金について](https://firebase.google.com/docs/firestore/pricing?hl=ja)
- [ストレージ サイズの計算](https://firebase.google.com/docs/firestore/storage-size?hl=ja)

## Storage

## group
1グループあたりのドキュメントサイズが200B程度と仮定すると、無料枠で下記が可能。
- 合計500万グループ
- 毎日2万件の折半リクエスト (1 read + 1 write)
が可能。

## payment
今のところindexはなくドキュメントあたり大体100B程度。
無料枠と照らし合わせると、データ保存は合計 1 GiBまで、すなわち1千万件まで無料で保存できる


## Reference

- [Cloud Firestore  |  Firebase Documentation](https://firebase.google.com/docs/firestore)
- [Firebase Pricing](https://firebase.google.com/pricing)
