build-customers:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/customers/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-customersDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/customersDelete/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-customersPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/customersPut/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-news:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/news/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-newsDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/newsDelete/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-newsPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/newsPut/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-orders:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/orders/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-ordersPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/ordersPut/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-productVariations:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/productVariations/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-products:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/products/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-productsImport:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/productsImport/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-sections:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/sections/main.go
	mv bootstrap $(ARTIFACTS_DIR)