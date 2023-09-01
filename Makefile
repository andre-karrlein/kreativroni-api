build-customers:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o customers cmd/customers/main.go
	mv customers $(ARTIFACTS_DIR)

build-customersDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o customersDelete cmd/customersDelete/main.go
	mv customersDelete $(ARTIFACTS_DIR)

build-customersPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o customersPut cmd/customersPut/main.go
	mv customersPut $(ARTIFACTS_DIR)

build-news:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o news cmd/news/main.go
	mv news $(ARTIFACTS_DIR)

build-newsDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o newsDelete cmd/newsDelete/main.go
	mv newsDelete $(ARTIFACTS_DIR)

build-newsPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o newsPut cmd/newsPut/main.go
	mv newsPut $(ARTIFACTS_DIR)

build-orders:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o orders cmd/orders/main.go
	mv orders $(ARTIFACTS_DIR)

build-ordersPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ordersPut cmd/ordersPut/main.go
	mv ordersPut $(ARTIFACTS_DIR)

build-product-variations:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o product-variations cmd/product-variations/main.go
	mv product-variations $(ARTIFACTS_DIR)

build-products:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o products cmd/products/main.go
	mv products $(ARTIFACTS_DIR)

build-products-import:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o products-import cmd/products-import/main.go
	mv products-import $(ARTIFACTS_DIR)

build-sections:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sections cmd/sections/main.go
	mv sections $(ARTIFACTS_DIR)