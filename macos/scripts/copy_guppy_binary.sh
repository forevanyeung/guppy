# Define the source and destination paths
SOURCE_PATH="${PROJECT_DIR}/../../cli/build/guppy"  # Adjust this path to where your Guppy binary is located
DEST_PATH="${BUILT_PRODUCTS_DIR}/${PRODUCT_NAME}.app/Contents/Resources"

# Ensure the destination directory exists
mkdir -p "$DEST_PATH"

# Copy the binary
cp "$SOURCE_PATH" "$DEST_PATH/guppy"

# Make the binary executable
chmod +x "$DEST_PATH/guppy"

echo "Copied Guppy binary to $DEST_PATH/guppy"
