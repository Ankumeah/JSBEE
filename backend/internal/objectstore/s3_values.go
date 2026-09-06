package objectstore

const backupBucket = "backup"
const publicBucket = "public"
const privateBucket = "private"

const backupBaseName = "jsbee.sql.bak."

var buckets []string = []string{
	backupBucket,
	publicBucket,
	privateBucket,
}
