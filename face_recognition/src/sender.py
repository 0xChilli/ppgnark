from flask import Flask, request, jsonify
from deepface import DeepFace
import numpy as np

app = Flask(__name__)

@app.route("/extract_features", methods=["POST"])
def extract_features():
    try:
        # Receive the image path from the request
        data = request.json
        image_path = data["image_path"]
        
        # Extract features using DeepFace
        embeddings = DeepFace.represent(img_path=image_path, model_name="Facenet")
        if embeddings and len(embeddings) > 0:
            feature_vector = embeddings[0]["embedding"]
            return jsonify({"status": "success", "embeddings": feature_vector})
        else:
            return jsonify({"status": "failure", "message": "No embeddings found"})
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)