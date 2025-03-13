## WGET Project Structure

```
wget/
│── main.go
│── fileDownload/
│   └── download.go
│── bckgrdDownload/
│   └── background.go
│── downloader/
│   └── resources.go
│── inputDownload/
│   └── batch.go
│── rateDownload/
│   └── rate_limit.go
│── mirrorDownload/
│   └── mirror.go
│   └── pathfix.go
```

each directory has its own specific go file that carrys out a specific function

### Here’s how each file will contribute:

1. fileDownload/download.go → Handles single file downloads.
2. bckgrdDownload/background.go → Implements background downloading.
3. downloader/resources.go → Supports the implementation of the background function.
4. inputDownload/batch.go → Supports batch downloads from a file.
5. rateDownload/rate_limit.go → Implements rate-limited downloads.
6. mirrorDownload/mirror.go → Implements website mirroring.
