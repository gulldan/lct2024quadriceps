from qdrant_client import QdrantClient

client = QdrantClient(host="localhost", port=18014)
print(client.collection_exists(collection_name="val_embbedings"))

