test:
	go test -v -count 1 ./...

test100:
	go test -v -count 100 ./...

race:
	go test -v -race -count 1 ./...


# Название файла с точкой входа
MAIN=main.go

# Где лежат исходники
APP_DIR=cmd/prove

# Путь где создать бинарник
TARGET_DIR=bin

# Операционная система
OS=`go env GOOS`

# Архитектура системы
ARCH=amd64

# Наименование бинарника
BIN_NAME=prove

# Версия сборки
VERSION_BUILD=`git describe --tags`

# Полное имя для бинарника Linux
BIN_NAME_L=$(BIN_NAME)_$(VERSION_BUILD)_$(OS)_$(ARCH)

# Полное имя для бинарника Windows
BIN_NAME_W=$(BIN_NAME)_$(VERSION_BUILD)_win_$(ARCH).exe

# Полное имя для бинарника Mac
BIN_NAME_M=$(BIN_NAME)_$(VERSION_BUILD)_mac_$(ARCH)

# Покрытие тестами
COVER_OUT=profiles/coverage.out

# Запечь версию и время создания сборки
BAKING="-X 'github.com/SETTER2000/prove/internal/app.OSString=`go env GOOS`' -X 'github.com/SETTER2000/prove/internal/app.archString=`go env GOARCH`' -X 'github.com/SETTER2000/prove/internal/app.dateString=`date`' -X 'github.com/SETTER2000/prove/internal/app.versionString=`git describe --tags`' -X 'github.com/SETTER2000/prove/internal/app.commitString=`git rev-parse HEAD`'"

# Подключение к базе данных
#DB=postgres://bob:Bob-2023@194.67.104.167:5432/prove?sslmode=disable
#DB=postgres://bob:Bob2023@postgres:5433/prove?sslmode=disable
DB=postgres://prove:DBprove_2023@127.0.0.1:5432/prove?sslmode=disable


# HELP по флагамq
short_h:
	./$(APP_DIR)/$(BIN_NAME) -h

# Запустить сервис prove сконфигурировав его от файла json (config/config.json)
conf:
	./$(APP_DIR)/$(BIN_NAME) -c config/config_tt.json


# Запустить перед запуском сервера в отдельном терминале (для OS Linux)
accrual:
	./cmd/accrual/accrual_linux_amd64 -a localhost:8088

# Скомпилировать и запустить бинарник сервиса prove (prove) с запечёнными аргументами сборки
short:
	go build -ldflags $(BAKING) -o $(APP_DIR)/$(BIN_NAME) $(APP_DIR)/$(MAIN)
	./$(APP_DIR)/$(BIN_NAME)

docker:
	go build -ldflags "-X 'github.com/SETTER2000/prove/internal/app.OSString=`go env GOOS`' -X 'github.com/SETTER2000/prove/internal/app.archString=`go env GOARCH`' -X 'github.com/SETTER2000/prove/internal/app.dateString=`date`' -X 'github.com/SETTER2000/prove/internal/app.versionString=`git describe --tags`' -X 'github.com/SETTER2000/prove/internal/app.commitString=`git rev-parse HEAD`'" -o prove ./cmd/prove/main.go
	./prove


# Скомпилировать и запустить бинарник сервиса prove (prove) с подключением к DB
short_d:
	go build -ldflags $(BAKING) -o $(APP_DIR)/$(BIN_NAME) $(APP_DIR)/$(MAIN)
	./$(APP_DIR)/$(BIN_NAME) -d $(DB)

cover:
	go test -v -count 1 -race -coverpkg=./... -coverprofile=$(COVER_OUT) ./...
	go tool cover -func $(COVER_OUT)
	go tool cover -html=$(COVER_OUT)
	rm $(COVER_OUT)

cover1:
	go test -v -count 1  -coverpkg=./... -coverprofile=cover.out.tmp ./...
	cat cover.out.tmp | grep -v mocks/*  > cover.out2.tmp
	cat cover.out2.tmp | grep -v log/*  > $(COVER_OUT)
	go tool cover -func $(COVER_OUT)
	go tool cover -html=$(COVER_OUT)
	rm cover.out.tmp cover.out2.tmp
	rm $(COVER_OUT)

# Запустить сервис с документацией
# Доступен здесь: http://rooder.ru:6060/pkg/github.com/SETTER2000/prove/?m=all
godoc:
	godoc -http rooder.ru:6060

# Запустить сервис с документацией в фоновом режиме
doc:
	godoc -http=rooder.ru:6060  -play >/dev/null &

# Скомпилировать и запустить бинарник сервиса с подключением к DB
# и разрешением чтения заголовка X-Real-IP или X-Forwarded-For
short_x:
	go build -tags pro -ldflags $(BAKING) -o $(APP_DIR)/$(BIN_NAME) $(APP_DIR)/$(MAIN)
	./$(APP_DIR)/$(BIN_NAME) -d $(DB) -resolve_ip_using_header true


short_g:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/controller/grpc/proto/app.proto



# Скомпилировать для Linux
build:
	go build -ldflags $(BAKING) -o $(TARGET_DIR)/$(BIN_NAME_L) $(APP_DIR)/$(MAIN)

# Запустить сервис с подключением к DB
# FILE_STORAGE_PATH=;DATABASE_DSN=postgres://prove:DBprove_2023@127.0.0.1:5432/prove?sslmode=disable
run:
	./$(TARGET_DIR)/$(BIN_NAME_L) -d $(DB)

# Скомпилировать для Linux и запустить c БД
build_d:
	go build -ldflags $(BAKING) -o $(TARGET_DIR)/$(BIN_NAME_L) $(APP_DIR)/$(MAIN)
	./$(TARGET_DIR)/$(BIN_NAME_L) -d $(DB)

# Запустить сервис in File
short_f:
	go build -ldflags $(BAKING) -o $(TARGET_DIR)/$(BIN_NAME) $(APP_DIR)/$(MAIN)
	./$(TARGET_DIR)/$(BIN_NAME) -d= -f=storage.txt

# Запустить сервис prove с протоколом HTTPS
hs:
	sudo -S ./$(TARGET_DIR)/$(BIN_NAME_L)  -s -d $(DB)

# Запустить сервис prove и с протоколом HTTPS в фоновом режиме
hsf:
	sudo ./$(APP_DIR)/$(BIN_NAME) -s >/dev/null &

# .. и запустить
build_r:
	./$(TARGET_DIR)/$(BIN_NAME_L)

# Скомпилировать для Windows
build_w:
	GOOS=windows GOARCH=$(ARCH) go build -ldflags $(BAKING)  -o $(TARGET_DIR)/$(BIN_NAME_W) $(APP_DIR)/$(MAIN)

# Скомпилировать для Mac
build_m:
	GOOS=darwin GOARCH=$(ARCH) go build -ldflags $(BAKING) -o $(TARGET_DIR)/$(BIN_NAME_M) $(APP_DIR)/$(MAIN)

# Удалить папку bin со скомпилированными версиями приложения
drop_files:
	rm -rf $(TARGET_DIR)

# Скомпилировать для Linux, Windows, Darwin
build_a: drop_files build build_m build_w
	file $(TARGET_DIR)/$(BIN_NAME_W)
	file $(TARGET_DIR)/$(BIN_NAME_L)   # получить дополнительную информацию об этом файле, подтвердив его сборку для Linux
	file $(TARGET_DIR)/$(BIN_NAME_M)   # получить дополнительную информацию об этом файле, подтвердив его сборку для Linux
