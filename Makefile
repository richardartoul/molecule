gen-proto:
	rm -rf ./src/proto/gen
	docker run -v `pwd`/src/proto:/defs namely/protoc-all -l go -d /defs
