from flask import Blueprint, request, jsonify, current_app as app
from app.routes.utils import is_valid_extension
from app.models.image import Image
from app.schemas.imageschema import ImageSchema
from werkzeug.utils import secure_filename
from pathlib import Path
from PIL import Image as PImage
import imagehash

search_bp = Blueprint("search", __name__)


@search_bp.route("/search", methods=("GET", "POST"))
def search():
    if request.method == "POST" and "image" in request.files:
        return search_by_image()
    elif request.method == "POST":
        return jsonify(error="You need to set one image!"), 400


def search_by_image():
    if len(request.files.getlist("image")) > 1:
        return jsonify(error="Please upload only one picture!"), 400

    image_schema = ImageSchema(many=True)
    image = request.files["image"]
    filename = secure_filename(image.filename)
    filepath = Path("{}/{}".format(app.config["UPLOAD_FOLDER"], filename))
    if not is_valid_extension(filepath):
        return jsonify(error="File in invalid!"), 400

    image.save(str(filepath))
    uploaded_image_hash = imagehash.phash(PImage.open(str(filepath)))
    similar_images = get_similar_images(uploaded_image_hash)
    filepath.unlink()
    return jsonify(images=image_schema.dump(similar_images))


def get_similar_images(compared_hash):
    # This is definitely not an optimal approach as we look over every image that have been upload.
    # A better solution for large amount of images would be to use MySQL which is able to do a BIT_COUNT
    # operation natively for easier comparison between images.
    # See: https://stackoverflow.com/questions/14925151/hamming-distance-optimization-for-mysql-or-postgresql

    images = Image.query.all()
    candidates = []
    # Let's set 8 bits as the distance threshold
    cutoff = 8

    for candidate in images:
        candidate_hash = imagehash.hex_to_hash(candidate.image_data.image_hash)
        if compared_hash - candidate_hash <= cutoff:
            candidates.append(candidate)

    return candidates
