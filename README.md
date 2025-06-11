# HeavyMetal: Music Streaming Mono Repo
HeavyMetal is a Music Streaming Mono Repo designed to Provide a Seamless Experience for Self Hosted Music Streaming.

**Built with Go**

Each Folder represents a separate service that runs indedepently to provide a complete music streaming experience.

1. [hosting_server/README.md](hosting_server/README.md) - The main server that hosts the music streaming service binary.
2. [indexer/README.md](indexer/README.md) - The service that indexes the music files and provides metadata for the Database.
3. [frontend/README.md](frontend/README.md) - The web frontend that provides a user interface for the music streaming service. (Not in Go)

## Important Note
The frontend is not built with Go obviously, and it is just a simple web interface that provides an UI for the music streaming service. It is not a full-fledged frontend, but it is enough to provide a basic user interface for the music streaming service.

It is there, just so that you can check if the backend is working or not.
By no means it is Optimized, however, I am working on a different client that can be used with this service which will be a full-fledged client with all the features you would expect from a music streaming service for your self hosted music.

## How to run?

1. Indexer
```bash
cd indexer
make build
# The built binary will be in the /tmp/ folder.
```

2. Hosting Server
```bash
cd hosting_server
make build
# The built binary will be in the /tmp/ folder.
```
3. Frontend
```bash
cd frontend
pnpm install
pnpm run build
# The built files will be in the /dist/ folder.
```

