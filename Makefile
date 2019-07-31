test:
	go test $$GOPATH/src/github.com/ledongthuc/licensechecker
update:
	git clone https://github.com/spdx/license-list-data;
	go get -u github.com/jteeuwen/go-bindata/...;
	go-bindata -o internal/data/license.go -pkg "data" -prefix "license-list-data/text/" license-list-data/text/*;
	go-bindata -o internal/toc/toc.go -pkg "toc" -prefix "license-list-data/json/" license-list-data/json/licenses.json license-list-data/json/exceptions.json;
	rm -rf license-list-data;
