from qdrant_client import QdrantClient
from qdrant_client.models import Distance, PointStruct, VectorParams


class QdrantClientApi:
    def __init__(
        self, qdrant_host: str, qdrant_port: int, collection_name: str = "audio_embeddings", embbedings_dim: int = 512
    ) -> None:
        self.qdrant_host = qdrant_host
        self.qdrant_port = qdrant_port
        self.collection_name = collection_name
        self.embbedings_dim = embbedings_dim

        self.qdrant_client = QdrantClient(host=qdrant_host, port=qdrant_port)

        # self.qdrant_client.recreate_collection(
        #     collection_name=collection_name, vectors_config=VectorParams(size=embbedings_dim, distance=Distance.COSINE)
        # )

        self.qdrant_client.recreate_collection(
            collection_name="val_embbedings", vectors_config=VectorParams(size=embbedings_dim, distance=Distance.COSINE)
        )

        self.id_counter = 0
        self.test_id_counter = 0

    def upload_vectors(self, embeddings_dict: dict[str, list[float]]) -> None:
        points = []
        for audio in embeddings_dict:
            embedding = embeddings_dict[audio]
            point = PointStruct(
                id=self.id_counter,
                vector=embedding,
                payload={
                    "audio": audio,
                },
            )
            points.append(point)
            self.id_counter += 1

        for point in points:
            self.qdrant_client.upsert(collection_name=self.collection_name, points=[point])

    def find_nearest_vectors(
        self, audios_paths: list[str], all_embbedings: dict[str, list[float]], score_treshold: float = 0.962
    ) -> list[str]:
        audio_hits = {}
        for audio in audios_paths:
            vector = all_embbedings[audio]

            hits = self.qdrant_client.search(collection_name=self.collection_name, query_vector=vector, limit=1000)

            hits_filtered = []
            for hit in hits:
                if hit.score >= score_treshold:
                    hits_filtered.append(hit.payload["audio"])

            audio_hits[audio] = hits_filtered

        return audio_hits
