from flask import Blueprint, request, jsonify, current_app as app
from werkzeug.utils import secure_filename
from itsdangerous import JSONWebSignatureSerializer, BadSignature
from google.cloud import storage
from google.cloud.exceptions import NotFound, BadRequest
from app.models import db
from app.models.user import User
from app.models.image import Image
from app.models.imagedata import ImageData
from app.routes.utils import is_valid_extension
from pathlib import Path
from PIL import Image as PImage
import os
import imagehash

upload_bp = Blueprint("upload", __name__)


@upload_bp.route("/upload", methods=("GET", "POST"))
def upload():

    if request.method == "POST" and "images" in request.files:

        if os.getenv("GOOGLE_APPLICATION_CREDENTIALS") is None:
            return jsonify(error="This application cannot upload to cloud storage"), 400

        access_token = request.form["token"]
        user_id = get_id_from_token(access_token)

        if user_id is None:
            return jsonify(error="Token not found or invalid"), 400

        user = User.query.filter_by(id=user_id).first()
        images_bucket = get_storage_bucket("utsuru-images")
        images = request.files.getlist("images")
        skipped_images = []

        for image in images:
            filename = secure_filename(image.filename)
            filepath = Path("{}/{}".format(app.config["UPLOAD_FOLDER"], filename))

            if Image.query.filter_by(filename=filename).first() is not None and not is_valid_extension(filepath):
                # Skip images whose filenames are already present
                skipped_images.append(image)
                continue

            # Save the file temporarily in the server's filesystem
            image.save(str(filepath.absolute()))
            image_hash = imagehash.phash(PImage.open(filepath))
            blob = save_to_bucket(images_bucket, user_id, filename, filepath)
            save_to_db(filename, blob.public_url, blob.size, user, image_hash)

        return jsonify(status="Success", not_uploaded_images=skipped_images)


def get_id_from_token(token):

    try:
        serializer = JSONWebSignatureSerializer(secret_key=app.config["SECRET_KEY"])
        user_id_json = serializer.loads(token)
        return user_id_json["id"]
    except BadSignature:
        return None


def get_storage_bucket(bucket_name):

    storage_client = storage.Client()
    try:
        bucket = storage_client.get_bucket(bucket_name)
    except (NotFound, BadRequest):
        # Create the images bucket
        # Note: By default, the bucket will be created on US multi-region with standard storage
        bucket = storage_client.create_bucket(bucket_name)
        print("Bucket {} is created!".format(bucket))
    return bucket


def save_to_bucket(bucket, user_id, filename, filepath):
    # Save the file to Google Cloud Storage
    blob = bucket.blob("{}/{}".format(user_id, filename))
    blob.upload_from_filename(str(filepath))
    blob.make_public()
    # Remove the file from the server after upload
    filepath.unlink()
    return blob


def save_to_db(filename, url, size, user, phash):
    # Save the image info in the database
    image_model = Image(filename=filename, url=url, size=size)
    image_data = ImageData(image_hash=str(phash))
    image_model.image_data = image_data
    user.images.append(image_model)
    db.session.add(user)
    db.session.commit()
