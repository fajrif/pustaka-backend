#!/bin/bash

# Database Seeding Script
# Wrapper for Go seed tool with friendly interface

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Load .env file if exists
if [ -f .env ]; then
  echo -e "${BLUE}Loading configuration from .env file...${NC}"
  export $(grep -v '^#' .env | xargs)
fi

# Configuration
SEED_BINARY="./bin/seed"
GO_FILE="./cmd/seed/main.go"

# Check if seed binary exists, build if not
check_seed_binary() {
  if [ ! -f "$SEED_BINARY" ]; then
    echo -e "${YELLOW}Seed tool not found. Building...${NC}"
    mkdir -p bin
    go build -o "$SEED_BINARY" "$GO_FILE"
    echo -e "${GREEN}✓ Seed tool built${NC}"
  fi
}

# Show usage
show_usage() {
  echo -e "${CYAN}==========================================${NC}"
  echo -e "${CYAN}  Database Seeding Tool${NC}"
  echo -e "${CYAN}==========================================${NC}"
  echo ""
  echo "Usage: $0 <command> [options]"
  echo ""
  echo -e "${GREEN}Commands:${NC}"
  echo "  all                  Run all seeders"
  echo "  <seeder_name>        Run specific seeder"
  echo "  list                 List available seeders"
  echo "  rebuild              Rebuild seed tool"
  echo "  help, -h, --help     Show this help message"
  echo ""
  echo -e "${GREEN}Available Seeders:${NC}"
  echo "  jenis_buku           Seed jenis_buku table"
  echo ""
  echo -e "${BLUE}Examples:${NC}"
  echo "  $0 all                        Run all seeders"
  echo "  $0 jenis_buku                 Run jenis_buku seeder only"
  echo "  $0 list                       List available seeders"
  echo ""
  echo -e "${YELLOW}Environment Variables:${NC}"
  echo "  DB_HOST         Database host (default: localhost)"
  echo "  DB_PORT         Database port (default: 5432)"
  echo "  DB_USER         Database user (default: deployer)"
  echo "  DB_PASSWORD     Database password"
  echo "  DB_NAME         Database name (default: pustaka_db)"
  echo ""
}

# Run all seeders
run_all_seeders() {
  check_seed_binary
  echo -e "${BLUE}Running all seeders...${NC}"
  echo ""
  "$SEED_BINARY"
}

# Run specific seeder
run_specific_seeder() {
  check_seed_binary
  seeder_name="$1"
  echo -e "${BLUE}Running seeder: $seeder_name${NC}"
  echo ""
  "$SEED_BINARY" "$seeder_name"
}

# List available seeders
list_seeders() {
  echo -e "${CYAN}==========================================${NC}"
  echo -e "${CYAN}  Available Seeders${NC}"
  echo -e "${CYAN}==========================================${NC}"
  echo ""
  echo -e "${GREEN}▶ curriculum${NC}       - Seed curriculum table with curriculum types"
  echo -e "${GREEN}▶ jenis_buku${NC}       - Seed jenis_buku table with book types"
  echo -e "${GREEN}▶ bidang_studi${NC}     - Seed bidang_studi table with field types"
  echo -e "${GREEN}▶ jenjang_studi${NC}    - Seed jenjang_studi table with degree types"
  echo -e "${GREEN}▶ kelas${NC}            - Seed kelas table with class types"
  echo -e "${GREEN}▶ merk_buku${NC}        - Seed merk_buku table with brand types"
  echo -e "${GREEN}▶ sales_associates${NC} - Seed sales_associates table with sales associate data"
  echo -e "${GREEN}▶ expeditions${NC}      - Seed expeditions table with all expedition data"
  echo -e "${GREEN}▶ cities${NC}           - Seed cities table with all indonesian cities"
  echo -e "${GREEN}▶ books${NC}            - Seed books table with defined data in seeds/files"
  echo ""
  echo -e "${YELLOW}To add more seeders:${NC}"
  echo "  1. Create a new file in seeds/ directory"
  echo "  2. Implement the seeder function"
  echo "  3. Register it in cmd/seed/main.go"
  echo ""
}

# Rebuild seed tool
rebuild_tool() {
  echo -e "${BLUE}Rebuilding seed tool...${NC}"
  rm -f "$SEED_BINARY"
  mkdir -p bin
  go build -o "$SEED_BINARY" "$GO_FILE"
  echo -e "${GREEN}✓ Seed tool rebuilt${NC}"
}

# Main command handler
case "${1:-help}" in
all)
  run_all_seeders
  ;;

list | ls)
  list_seeders
  ;;

rebuild)
  rebuild_tool
  ;;

help | -h | --help)
  show_usage
  ;;

*)
  if [ -n "$1" ]; then
    # Try to run as specific seeder
    run_specific_seeder "$1"
  else
    show_usage
  fi
  ;;
esac
