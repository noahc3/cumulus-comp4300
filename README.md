# Cumulus - Game Server Control Panel & Daemon (COMP4300 Final Project)

Implementation of core technologies required for deploying and communicating with game servers hosted in the
Cloud.

## How to run

There are three components to this project:

- Frontend
- Daemon
- Hetzner Auxiliary API

**Important: Only Ubuntu 22.04 and MacOS have been tested!** To run on Windows, please use WSL2 with Ubuntu.

### Running the Frontend

Bun is an alternative JavaScript runtime and package manager which was used to create this project.

- Install [bun](https://bun.sh/docs/installation)
- Open a terminal in the `/project/frontend` directory
- Run `bun install`
- Start the dev server with `bun run dev`.

Nodejs may work, but is untested.

### Running the daemon

The daemon is written in Go.

- Install [go](https://go.dev/doc/install)
- Open a terminal in the `/project/final/cloud` directory
- Start the daemon with `go run . --host <ip> --port <port>`
    - host and port are optional, will default to `0.0.0.0:1111`
- Build a production binary with `go build`

### Running the Hetzner Auxiliary API

The auxiliary is built with Python and FastAPI.

- Install [Python](https://www.python.org/downloads/)
    - Note: Tested with Python 3.10.12, other versions may work but are untested
- Open a terminal in the `/project/final/hetznerfuncs` directory
- Install uvicorn: `pip install uvicorn`
- Install other required packages: `pip install -r requirements`
- Start the server with `uvicorn main:app --reload`

**Important:** You must create a file named `.env` in this directory and supply a Hetzner API key:

```env
# hetznerfuncs/.env

HETZNER_TOKEN=<api-key>
```

You can create a Hetzner Cloud account here: <https://www.hetzner.com/cloud/>

Billing is hourly so it isn't expensive to launch some test servers, but we recommend just watching the demo
video if you don't want to throw away money.

You will need to adjust the `main.py` script to specify your own SSH key (or remove it, Hetzner will email you
a root password).

## Other Files

While the above is all you need to worry about for running the final implementation of our project, we also
include source code for the test programs we wrote during the exploration phase of our project, detailed below.

- **docs/report**: Source for our project report. This is a mix of Markdown and Latex rendered to PDF using Pandoc.
- **project/cloud**: Our initial cloud daemon with some bugs, we do not recommend running this.
- **project/final** Directory with our final cloud daemon and hetzner functions
    - **project/gs-go-daemon.service** - systemd service file to launch the daemon automatically. Pulled in the cloud-init.yml script.
- **project/frontend** - Full frontend including testing interfaces
- **project/httpclient** - Simple WS client we used to generate latency data results.
- **project/middleware** - Simple middleware server we used to generate latency data results.
- **project/mockprocess** - Two simple programs, one to echo any input that is typed, and another which will read a mock log file and spit it out to standard output at specified lines per second.
- **project/tcpclient** - Test TCP client
- **project/udpclient** - Test UDP client
- **project/wsclient** - Test WebSocket client

