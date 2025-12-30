#!/bin/bash
set -e

VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."
ASSETS_DIR="$PROJECT_ROOT/assets"
DB_GEN_DIR="$PROJECT_ROOT/cmd/generate-nutrition-db"

echo "🚀 Bundling assets for version $VERSION"
echo "================================================"

# Create directories
mkdir -p "$ASSETS_DIR/models"
mkdir -p "$ASSETS_DIR/databases"

# 1. Generate nutrition database
echo ""
echo "📊 Generating nutrition database..."
cd "$DB_GEN_DIR"
go run main.go

# Move to assets directory
mv nutrition.db "$ASSETS_DIR/databases/nutrition_v${VERSION}.db"
echo "✅ Database generated: nutrition_v${VERSION}.db"

# 2. Calculate checksums
echo ""
echo "🔐 Calculating checksums..."
cd "$ASSETS_DIR"

# Database checksum
if [[ "$OSTYPE" == "darwin"* ]]; then
    DB_CHECKSUM=$(shasum -a 256 "databases/nutrition_v${VERSION}.db" | awk '{print $1}')
    DB_SIZE=$(stat -f%z "databases/nutrition_v${VERSION}.db")
else
    DB_CHECKSUM=$(sha256sum "databases/nutrition_v${VERSION}.db" | awk '{print $1}')
    DB_SIZE=$(stat -c%s "databases/nutrition_v${VERSION}.db")
fi

echo "Database checksum: sha256:$DB_CHECKSUM"
echo "Database size: $DB_SIZE bytes"

# 3. Create checksums.json
cat > checksums.json <<EOF
{
  "version": "$VERSION",
  "generated_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "assets": {
    "nutrition_db": {
      "version": "$VERSION",
      "file": "databases/nutrition_v${VERSION}.db",
      "checksum": "sha256:$DB_CHECKSUM",
      "size_bytes": $DB_SIZE
    }
  }
}
EOF

echo ""
echo "✅ Checksums generated"
echo ""
echo "📦 Asset Bundle Summary:"
echo "  - Version: $VERSION"
echo "  - Database: nutrition_v${VERSION}.db ($DB_SIZE bytes)"
echo "  - Checksum: $DB_CHECKSUM"
echo ""
echo "🎉 Asset bundling complete!"
echo "================================================"
