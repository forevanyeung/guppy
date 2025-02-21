# Define the source and destination paths
SOURCE_PATH="${PROJECT_DIR}/../../../cli/guppy"  # Adjust this path to where your Guppy binary is located
DEST_PATH="${BUILT_PRODUCTS_DIR}/${FRAMEWORKS_FOLDER_PATH}"

# Ensure the destination directory exists
mkdir -p "$DEST_PATH"

# Copy the binary
cp "$SOURCE_PATH" "$DEST_PATH/guppy"

# Make the binary executable
chmod +x "$DEST_PATH/guppy"

echo "Copied Guppy binary to $DEST_PATH/guppy"
