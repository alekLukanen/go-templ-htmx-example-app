GREEN='\033[0;32m'
NC='\033[0m'

codegen:
	@echo "${GREEN}Generating templ code...${NC}"
	templ generate ./core/ui/
	@echo "${GREEN}Generating SQLC code...${NC}"
	sqlc generate -f ./core/sqlc.yaml
	@echo "${GREEN}Generating Tailwind CSS code...${NC}"
	npx tailwindcss -c ./core/ui/tailwind.config.js -i ./core/ui/main.css -o ./core/ui/static/main.css --minify

migrate_up:
	@echo "${GREEN}Migrating up...${NC}"
	go run cmd/migrate/main.go upOne

build_ui: codegen
	@echo "${GREEN}Building UI...${NC}"
	go build -o bin/ui/ui cmd/ui/main.go
	cp -r ./core/ui/static ./bin/ui/

run_ui: codegen
	@echo "${GREEN}Running UI...${NC}"
	go run cmd/ui/main.go

