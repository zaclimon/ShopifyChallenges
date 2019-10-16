def is_valid_extension(filepath):
    """Checks if the file is a valid image."""
    return filepath.suffix in (".jpg", ".jpeg", ".gif", ".png")
