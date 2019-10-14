# Check if the file is a valid image
def is_valid_extension(filepath):
    return filepath.suffix in (".jpg", ".jpeg", ".gif", ".png")