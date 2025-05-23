.PHONY: build run clean demo help

build:
	go build -o simple-git .

run:
	go run . $(ARGS)

clean:
	rm -f simple-git
	rm -rf .git
	rm -f test.txt
	
demo: build
	@echo "=== Initializing repository ==="
	./simple-git init
	@echo "=== Creating a test file ==="
	echo "Hello, Simple Git!" > test.txt
	@echo "=== Adding file to staging ==="
	./simple-git add test.txt
	@echo "=== Creating first commit ==="
	./simple-git commit -m "Initial commit"
	@echo "=== Showing status ==="
	./simple-git status
	@echo "=== Showing log ==="
	./simple-git log

