APP := jh_ldy
JH_APP := /Users/wangdante/D/kugou/jh/jh_app
APK_SRC := $(JH_APP)/build/app/outputs/flutter-apk/app-release.apk
APK_DST := downloads/app-release.apk

.PHONY: run test apk docker-build

run:
	go run .

test:
	go test ./...
	node web/js/source.test.js

apk:
	@test -f $(APK_SRC) || (echo "先在 jh_app 执行: flutter build apk --release"; exit 1)
	mkdir -p downloads
	cp "$(APK_SRC)" "$(APK_DST)"
	@ls -lh "$(APK_DST)"

docker-build:
	docker build -t jh_ldy:latest .
