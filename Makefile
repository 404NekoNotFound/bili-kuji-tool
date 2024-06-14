gen:
	go build -o bili-kuji-tui *.go

amd64linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bili-kuji-tui *.go
	tar zcvf bili-kuji-tui-linux-amd64.tar.gz bili-kuji-tui

arm64linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o bili-kuji-tui *.go
	tar zcvf bili-kuji-tui-linux-arm64.tar.gz bili-kuji-tui

amd64windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bili-kuji-tui.exe *.go
	tar zcvf bili-kuji-tui-windows-amd64.tar.gz bili-kuji-tui.exe

amd64mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bili-kuji-tui *.go
	tar zcvf bili-kuji-tui-macOS-amd64.tar.gz bili-kuji-tui

clean:
	rm bili-kuji-tui*