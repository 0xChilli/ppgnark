from deepface import DeepFace
import numpy as np

def extract_features(image_path, model_name='Facenet'):
    """
    Extracts feature embeddings from an image using DeepFace.
    
    Parameters:
    - image_path (str): Path to the image file.
    - model_name (str): DeepFace model to use for feature extraction. Default is 'Facenet'.
      Other options include 'VGG-Face', 'OpenFace', 'DeepFace', 'DeepID', and 'ArcFace'.

    Returns:
    - np.ndarray: Feature embeddings as a numpy array.
    - str: Message if analysis fails or the image is invalid.
    """
    try:
        # Analyze image using the specified model
        embeddings = DeepFace.represent(img_path=image_path, model_name=model_name)
        
        # Extract the feature vector (embedding)
        if embeddings and len(embeddings) > 0:
            feature_vector = np.array(embeddings[0]["embedding"])
            return feature_vector
        
        return "No embeddings found for the image."
    except Exception as e:
        return f"Error in extracting features: {str(e)}"

# Example usage
if __name__ == "__main__":
    # Path to your image
    image_path = "/Users/archer/vscode/bachelor_project/ppzkp/face_recognition/dataset/lfw-deepfunneled/Alan_Ball/Alan_Ball_0001.jpg"
    
    # Extract features
    features = extract_features(image_path, model_name='Facenet')   # Use Facenet for larger embeddings 
                                                                    # VGG-Face: 2622 dimensions (dense but smaller embeddings).
	                                                                # Facenet: 128 dimensions (more discriminative and compact).
	                                                                # ArcFace: 512 dimensions (large and highly discriminative).
    
    # Check and print the result (for my Checking)
    if isinstance(features, np.ndarray):
        print(f"Feature Vector (Size {features.size}):")
        print(features)  # Prints the entire embedding vector
    else:
        print(features)