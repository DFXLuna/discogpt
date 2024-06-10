## Chroma
DiscoGPT supports using [ChromaDB](https://www.trychroma.com/) for retrieval augmented generation. This allows DiscoGPT bots to have moderately intelligent, long-term memory of user interactions.

### Operation
At runtime, when DiscoGPT is triggered, the contents of that message get turned into an embedding using a HuggingFace TEI server. That embedding is stored in ChromaDB. This means that you are storing any user message starting with the DiscoGPT trigger in Chroma INDEFINITELY until you take action to remove it.

Additionally, ChromaDB is queried with the content of the user message to find similar messages based on semantic similarity.

Even though DiscoGPT has no telemetry, ChromaDB and HuggingFace TEI may or may not have telemetry and you should inspect their documentation thoroughly before using them or this feature.

### Setup
- You'll need three pieces
    - DiscoGPT
    - [Chroma-core server](https://github.com/chroma-core/chroma/pkgs/container/chroma)
    - [HuggingFace TEI server](https://github.com/huggingface/text-embeddings-inference/pkgs/container/text-embeddings-inference)
- Add the full URLs (protocol included) into your config.yaml as described in the [top level readme](../README.md), DiscoGPT will automatically enable the chroma mod.

### Networking
- Depending on your setup, networking can be a little tricky.
  - All in one compose file
      - Use network links and service names with the internal container ports
      - DiscoGPT's config.yaml will also use the service name as the host name in the url 
        - e.g. With Chroma service `chroma`: `ChromaURL: http://chroma:8000`
  - Multiple compose groups or DiscoGPT otherwise running outside of the docker network that Chroma & HF TEI are running on
      - Use the standard docker compose network stuff: expose ports

### Embedding model selection
The example below uses `sentence-transformers/all-MiniLM-L6-v2` because it's small enough to run on the CPU without trouble. You may also use any of TEI's [supported models](https://github.com/huggingface/text-embeddings-inference?tab=readme-ov-file#supported-models). Keep in mind that larger models may require GPU acceleration, which requires extra setup in your container runtime.

### Example docker compose
```yaml
version: '3'

services:
    image: dfxluna/discogpt:latest
    volumes:
      - type: bind # Mount your config file with ChromaURL, ChromaTEIURL and CollectionName specified
        source: ./config.yaml
        target: /discogpt/config.yaml 
    restart: "unless-stopped"
  chroma:
    image: ghcr.io/chroma-core/chroma:0.4.24
    # ports:  # Add this if discogpt is accessing chroma from outside the compose group's network
    #   - 8000:8000
    volumes:
      - ./chroma-data:/chroma/chroma # Persistant storage is required if you want to keep data between restarts
    environment:
      - IS_PERSISTANT=TRUE #This enables writing to disk
    restart: "unless-stopped"
    healthcheck: # Health check isn't necessary
      test:
        [
          "CMD",
          "curl",
          "-f",
          "http://localhost:8000/api/v1/heartbeat"
        ]
      interval: 30s
      timeout: 10s
      retries: 3
  embedding-server:
    image: ghcr.io/huggingface/text-embeddings-inference:cpu-1.2 # change to gpu if planning on using accelerated embedding
    # ports:     # Add this if discogpt is accessing chroma from outside the compose group's network
    #   - 8080:80
    volumes:
      - ./tei-data:/data # Persistant storage is use to avoid pulling the model on every startup
    environment:
      - MODEL_ID=sentence-transformers/all-MiniLM-L6-v2 # The model's ID on HuggingFace
      - REVISION=8b3219a92973c328a8e22fadcfa821b5dc75636a # Choose the specific revision, this one uses a commit ID
    restart: "unless-stopped"
```
