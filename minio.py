#!/usr/bin/env python
# coding: utf-8

# In[1]:


from glob import glob
import requests
import warnings
import json
from tqdm import tqdm

warnings.filterwarnings("ignore")


# In[2]:


bff_port = ""
bff_host = ""

index_path = "rutube_data/rutube/index/compressed_index/"

all_videos = glob(index_path + "*")


# In[3]:


def update_database(
    file_path: str,
    bff_host: str,
    bff_port: str,
    not_upload_embeddings: bool
    ) -> None:
    url = f"http://{bff_host}:{bff_port}/api/v1/original/upload"
    with open(file_path, "rb") as file:
        file_name = file_path.split("/")[-1]
        response = requests.post(
            url,
            files={"file": file},
            data={
                "not_upload_embeddings": not_upload_embeddings,
                "name": file_name
            }
        )

    if response.ok:
        print(f"Video {file_name} uploaded")
    else:
        print(
            f"Request failed with status code {response.status_code} for video {file_name}"
        )
        return None


# In[4]:


for video in tqdm(all_videos):
    update_database(
        video,
        bff_host,
        bff_port,
        True
    )

