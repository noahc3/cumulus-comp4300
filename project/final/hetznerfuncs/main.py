import os
import json
import socket

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from dotenv import load_dotenv

from hcloud import Client
from hcloud.images import Image
from hcloud.server_types import ServerType
from hcloud.locations import Location


app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

load_dotenv()

HETZNER_TOKEN = os.getenv("HETZNER_TOKEN")
hetzner = Client(token=HETZNER_TOKEN)

def get_hetzner_server_type(name: str):
    type = hetzner.server_types.get_by_name(name)
    return type

def get_hezner_location(location: str):
    loc = hetzner.locations.get_by_name(location)
    return loc

@app.get("/api/createserver")
async def create_server(location: str, instance_type: str, name: str):
    type = get_hetzner_server_type(instance_type.lower())
    loc = get_hezner_location(location.lower())
    image = hetzner.images.get_by_name("ubuntu-22.04")
    user_data = ""

    with open("cloud-init.yaml", 'r') as file:
        user_data = file.read()

    hetzner_server = hetzner.servers.create(
        name=name.replace(" ", "-").lower(),
        server_type=type,
        image=image,
        location=loc,
        user_data=user_data,
        ssh_keys=[hetzner.ssh_keys.get_by_name("noahcuroe@gmail.com")]
    )

    server = {
        "id": hetzner_server.server.id,
        "name": hetzner_server.server.name,
        "status": hetzner_server.server.status,
        "public_net": hetzner_server.server.public_net.ipv4.ip,
        "server_type": hetzner_server.server.server_type.name,
        "location": hetzner_server.server.datacenter.name
    }
    
    return JSONResponse(content=server)

@app.get("/api/getservers")
async def create_server():
    servers = []
    for server in hetzner.servers.get_all():
        if server.name == "nbg1-test-server" or server.name == "hil-test-mitm-server":
            continue

        ip = server.public_net.ipv4.ip
        
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(1)
        result = sock.connect_ex((ip, 1111))

        servers += [{
            "id": server.id,
            "name": server.name,
            "status": server.status,
            "ready": result == 0,
            "public_net": server.public_net.ipv4.ip,
            "server_type": server.server_type.name,
            "location": server.datacenter.name
        }]

    return JSONResponse(content=servers)