import WebsocketTerminal from "@/components/terminal";
import WebsocketClient from "@/lib/websocketclient";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createLazyFileRoute("/direct-via-ws")({
    component: About
});

function About() {

    const [wsclient, _] = useState<WebsocketClient>(new WebsocketClient());

    useEffect(() => {
        wsclient.connect("localhost", 1111);
    }, [wsclient]);

    return (
        <div className="p-2">
            <h1 className="text-2xl font-bold">Direct via WebSocket</h1>
            <WebsocketTerminal client={wsclient} />
        </div>
    );
}
