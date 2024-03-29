# Excalidraw Complete: A Self-Hosted Solution

Excalidraw Complete simplifies the deployment of Excalidraw, bringing an
all-in-one solution to self-hosting this versatile virtual whiteboard. Designed
for ease of setup and use, Excalidraw Complete integrates essential features
into a single Go binary. This solution encompasses:

- The intuitive Excalidraw frontend UI for seamless user experience.
- An integrated data layer ensuring fast and efficient data handling based on different data providers.
- A socket.io implementation to enable real-time collaboration among users.

The project goal is to alleviate the setup complexities traditionally associated with self-hosting Excalidraw, especially in scenarios requiring data persistence and collaborative functionalities.

## Installation

To get started, download the latest release binary:

```bash
# Visit https://github.com/PatWie/excalidraw-complete/releases/ for the download URL
wget <binary-download-url>
chmod +x excalidraw-complete
./excalidraw-complete
```

Once launched, Excalidraw Complete is accessible at `localhost:3002`, ready for
drawing and collaboration.

### Configuration

Excalidraw Complete adapts to your preferences with customizable storage solutions, adjustable via the `STORAGE_TYPE` environment variable:

- **Filesystem:** Opt for `STORAGE_TYPE=filesystem` and define `LOCAL_STORAGE_PATH` to use a local directory.
- **SQLite:** Select `STORAGE_TYPE=sqlite` with `DATA_SOURCE_NAME` for local SQLite storage, including the option for `:memory:` for ephemeral data.
- **AWS S3:** Choose `STORAGE_TYPE=s3` and specify `S3_BUCKET_NAME` to leverage S3 bucket storage, ideal for cloud-based solutions.

These flexible configurations ensure Excalidraw Complete fits seamlessly into your existing setup, whether on-premise or in the cloud.

## Building from Source

Interested in contributing or customizing? Build Excalidraw Complete from source with these steps:

```bash
# Clone and prepare the Excalidraw frontend
git clone https://github.com/excalidraw/excalidraw.git
cd excalidraw
git checkout tags/v0.17.3
git apply ../frontend.patch
# Follow build instructions to compile assets into the frontend directory
```

Compile the Go application:

```bash
go build -o excalidraw-complete main.go
```

Start the server:

```bash
./excalidraw-complete
```

Excalidraw Complete is now running on your machine, ready to bring your collaborative whiteboard ideas to life.

---

Excalidraw is a fantastic tool, but self-hosting it can be tricky. I welcome
your contributions to improve Excalidraw Complete — be it through adding new
features, improving existing ones, or bug reports.
