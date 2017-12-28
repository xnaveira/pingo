.PHONY: clean

CONTAINER := xnaveira/pingo


pingo: */*.go
	docker build -t ${CONTAINER} . && touch pingo

clean:
	@rm pingo || true
	docker rmi $(CONTAINER)
