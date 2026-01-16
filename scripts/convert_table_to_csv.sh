#!/bin/bash

# Script to convert markdown table (table_books.md) to CSV format
# Appends to books.csv in seeds/files directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
INPUT_FILE="$PROJECT_ROOT/seeds/files/table_books.md"
OUTPUT_FILE="$PROJECT_ROOT/seeds/files/books.csv"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "❌ Error: Input file not found: $INPUT_FILE"
    echo "   Please create table_books.md with your markdown table data."
    exit 1
fi

# Create output file with header if it doesn't exist
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "📝 Creating new books.csv with header..."
    echo "bidang_studi_code,jenjang_code,curriculum_code,kelas,pages,lks_price,pg_price,periode,year,merk_code,stock,publisher_code" > "$OUTPUT_FILE"
fi

# Count rows before processing
ROWS_BEFORE=$(wc -l < "$OUTPUT_FILE" | tr -d ' ')

# Check if table has PG column by reading header line
HAS_PG_COLUMN=false
while IFS= read -r line; do
    if echo "$line" | grep -qi "Bidang Studi"; then
        if echo "$line" | grep -qi "| *PG *|"; then
            HAS_PG_COLUMN=true
        fi
        break
    fi
done < "$INPUT_FILE"

if [ "$HAS_PG_COLUMN" = true ]; then
    echo "📊 Detected 11-column format (with PG column)"
else
    echo "📊 Detected 10-column format (without PG column)"
fi

# Process the markdown table
# Skip header row (line 1) and separator row (line 2 with dashes)
# Convert remaining rows from markdown table to CSV
echo "📖 Reading from: $INPUT_FILE"
echo "📝 Appending to: $OUTPUT_FILE"

while IFS= read -r line; do
    # Skip empty lines
    [ -z "$line" ] && continue
    
    # Skip separator lines (lines containing consecutive dashes like ---)
    if echo "$line" | grep -q -- '---'; then
        continue
    fi
    
    # Skip header line (contains "Bidang Studi")
    if echo "$line" | grep -qi "Bidang Studi"; then
        continue
    fi
    
    # Check if line is a valid table row (starts and ends with |)
    if echo "$line" | grep -qE '^\|.*\|$'; then
        # Remove leading and trailing |, then split by |
        # Trim whitespace from each field
        if [ "$HAS_PG_COLUMN" = true ]; then
            # 11-column format: bidang_studi,jenjang,curriculum,kelas,pages,lks,pg,periode,year,merk,stock[,publisher]
            csv_line=$(echo "$line" | \
                sed 's/^|//; s/|$//' | \
                awk -F'|' '{
                    for(i=1; i<=NF; i++) {
                        gsub(/^[ \t]+|[ \t]+$/, "", $i)
                        printf "%s", $i
                        if(i<NF) printf ","
                    }
                    printf "\n"
                }')
        else
            # 10-column format: bidang_studi,jenjang,curriculum,kelas,pages,lks,periode,year,merk,stock[,publisher]
            # Insert empty pg_price after lks_price (position 6)
            csv_line=$(echo "$line" | \
                sed 's/^|//; s/|$//' | \
                awk -F'|' '{
                    for(i=1; i<=NF; i++) {
                        gsub(/^[ \t]+|[ \t]+$/, "", $i)
                        printf "%s", $i
                        if(i==6) printf "," # Add empty pg_price after lks_price
                        if(i<NF) printf ","
                    }
                    printf "\n"
                }')
        fi
        
        # Append to CSV file
        echo "$csv_line" >> "$OUTPUT_FILE"
    fi
done < "$INPUT_FILE"

# Count rows after processing
ROWS_AFTER=$(wc -l < "$OUTPUT_FILE" | tr -d ' ')
ROWS_ADDED=$((ROWS_AFTER - ROWS_BEFORE))

echo ""
echo "✅ Conversion complete!"
echo "   Rows added: $ROWS_ADDED"
echo "   Total rows in CSV: $((ROWS_AFTER - 1)) (excluding header)"
echo ""
echo "💡 Next steps:"
echo "   1. Clear table_books.md if you want to add more data later"
echo "   2. Run the seed: go run cmd/seed/main.go books"
