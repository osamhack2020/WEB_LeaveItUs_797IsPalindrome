GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_HASH=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%F-%H:%M:%S')

info:
	@echo "[leave it us make]\nbuild information : ${GIT_TAG} - ${GIT_HASH} (${BUILD_DATE})"
	@echo "명령어 목록 : build, info"

build:
	rm -rf ./output && mkdir output
	cd frontend && npm run build
	cp -r frontend/dist output/static
	cd backend && go build -v -x -o ../output/

serve: 
	cd output && ./backend

clean:
	rm -rf ./output