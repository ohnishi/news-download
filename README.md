# nahaha
トレンド分析

cd /Users/ohnishi/home/go/src/github.com/ohnishi/nahaha


### fetch 5ch thread
go run github.com/ohnishi/nahaha/backend/cmd/nahahafetch 5ch --dest /Users/ohnishi/home/go/data/nahaha/fetch/5ch

### fetch yahoo thread
go run github.com/ohnishi/nahaha/backend/cmd/nahahafetch yahoo --dest /Users/ohnishi/home/go/data/nahaha/fetch/rss

### fetch rss thread
go run github.com/ohnishi/nahaha/backend/cmd/nahahafetch rss --src /Users/ohnishi/home/go/data/nahaha/fetch/rss --dest /Users/ohnishi/home/go/data/nahaha/fetch/rss


### transform 5ch thread
go run github.com/ohnishi/nahaha/backend/cmd/nahahatransform 5ch --src /Users/ohnishi/home/go/data/nahaha/fetch/5ch --dest /Users/ohnishi/home/go/data/nahaha/transform --date 20201030

### transform rss thread
go run github.com/ohnishi/nahaha/backend/cmd/nahahatransform rss --src /Users/ohnishi/home/go/data/nahaha/fetch/rss --dest /Users/ohnishi/home/go/data/nahaha/transform --date 20201030

### transform analysis trends
go run github.com/ohnishi/nahaha/backend/cmd/nahahaanalysis trends --src /Users/ohnishi/home/go/data/nahaha/transform --dest /Users/ohnishi/home/go/data/nahaha/trends --date 20201031

### transform analysis trends
go run github.com/ohnishi/nahaha/backend/cmd/nahahapublish trends --src /Users/ohnishi/home/go/data/nahaha/trends --dest /Users/ohnishi/home/go/src/github.com/ohnishi/nahaha/hugo/content/posts --date 20201031
