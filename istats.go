package localcache

type IStats interface {
	miss()
	delHit()
	delMiss()
	collision()
	hit(key string)
	getMisses() int64
	getDelHits() int64
	getDelMisses() int64
	getCollisions() int64
	getHits() int64
	getKeyHits(key string) int64
}
