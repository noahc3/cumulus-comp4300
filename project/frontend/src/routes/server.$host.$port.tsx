import WebsocketTerminal from "@/components/terminal";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import WebsocketClient from "@/lib/websocketclient";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export const Route = createFileRoute("/server/$host/$port")({
    component: ServerViewer
});

function ServerViewer() {
    const { host, port } = Route.useParams();
    const [wsclient, _] = useState<WebsocketClient>(new WebsocketClient());

    const [cmd, setCmd] = useState("java -Xmx1024M -Xms1024M -jar server.jar nogui");
    const [wd, setWd] = useState("/opt/mc");

    const [cpu1s, setCpu1s] = useState("0");
    const [cpu10s, setCpu10s] = useState("0");
    const [cpu30s, setCpu30s] = useState("0");
    const [cpu60s, setCpu60s] = useState("0");
    const [memUsed, setMemUsed] = useState("0");
    const [memTotal, setMemTotal] = useState("0");
    const [diskUsed, setDiskUsed] = useState("0");
    const [diskTotal, setDiskTotal] = useState("0");

    const onDiagnosticEvent = (data: WsMessage) => {
        if (data.command == "cpu1" && data.value) {
            setCpu1s(data.value.toFixed());
        } else if (data.command == "cpu10" && data.value) {
            setCpu10s(data.value.toFixed());
        } else if (data.command == "cpu30" && data.value) {
            setCpu30s(data.value.toFixed());
        } else if (data.command == "cpu60" && data.value) {
            setCpu60s(data.value.toFixed());
        } else if (data.command == "mem" && data.used && data.total) {
            setMemUsed(data.used.toFixed(1));
            setMemTotal(data.total.toFixed());
        } else if (data.command == "disk" && data.used && data.total) {
            setDiskUsed((data.used / 1024.0).toFixed(2));
            setDiskTotal((data.total / 1024.0).toFixed());
        }
    }

    useEffect(() => {
        wsclient.setOnDiagnosticMessage(onDiagnosticEvent);
        wsclient.connect(host, parseInt(port));
    }, [wsclient]);

    const startServer = () => {
        const msg = {
            command: "start",
            target: cmd,
            directory: wd
        }

        wsclient.sendJson(msg);
    }

    const restartServer = () => {
        const msg = {
            command: "restart",
            target: cmd,
            directory: wd
        }

        wsclient.sendJson(msg);
    }

    const stopServer = () => {
        const msg = {
            command: "stop"
        }

        wsclient.sendJson(msg);
    }


    return (
        <div className="p-2">
            <h1 className="text-2xl font-bold mb-3 p-1">Server {host}:{port}</h1>
            <Tabs defaultValue="console" className="w-full mb-4">
                <TabsList className="grid w-full grid-cols-7">
                    <TabsTrigger value="console">Console</TabsTrigger>
                    <TabsTrigger value="files" disabled>Files</TabsTrigger>
                    <TabsTrigger value="databases" disabled>Databases</TabsTrigger>
                    <TabsTrigger value="backups" disabled>Backups</TabsTrigger>
                    <TabsTrigger value="network" disabled>Network</TabsTrigger>
                    <TabsTrigger value="startup" disabled>Startup</TabsTrigger>
                    <TabsTrigger value="settings" disabled>Settings</TabsTrigger>
                </TabsList>
            </Tabs>
            <div className="grid grid-cols-[250px_auto] gap-4">
                <div className="flex flex-col justify-between">
                    <Card className="p-4">
                        <div className="text-xs text-muted-foreground flex gap-2">
                            <Button className="p-4 h-8 m-0" size="sm" variant="outline" onClick={startServer}>
                                Start
                            </Button>
                            <Button className="p-4 h-8" size="sm" variant="outline" onClick={restartServer}>
                                Restart
                            </Button>
                            <Button className="p-4 h-8" size="sm" variant="destructive" onClick={stopServer}>
                                Stop
                            </Button>
                        </div>
                    </Card>
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-0">
                            <CardTitle className="text-sm font-medium">CPU Usage</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{cpu1s}%</div>
                            <div className="text-xs text-muted-foreground flex justify-between">
                                <div className="flex-col">10s {cpu10s}%</div>
                                <div className="flex-col">30s {cpu30s}%</div>
                                <div className="flex-col">1m {cpu60s}%</div>
                            </div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-0">
                            <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">
                                {memUsed} MB <span className="text-xs">/ {memTotal} MB</span>
                            </div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-0">
                            <CardTitle className="text-sm font-medium">Disk Usage</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">
                                {diskUsed} GB <span className="text-xs">/ {diskTotal} GB</span>
                            </div>
                        </CardContent>
                    </Card>
                </div>
                <div className="flex-col p-4 bg-black border rounded-lg w-full">
                    <WebsocketTerminal client={wsclient} />
                </div>
                <div></div>
                <div className="flex flex-col gap-4">
                    <div className="grid w-full items-center gap-1.5">
                        <Label htmlFor="command">Launch Command</Label>
                        <Input id="command" placeholder="java -jar ..." value={cmd} onChange={(e) => {setCmd(e.target.value);}}/>
                    </div>
                    <div className="grid w-full items-center gap-1.5">
                        <Label htmlFor="wd">Working Directory</Label>
                        <Input id="wd" placeholder="/opt/mc" value={wd} onChange={(e) => {setWd(e.target.value);}}/>
                    </div>
                </div>
            </div>
        </div>
    );
}
