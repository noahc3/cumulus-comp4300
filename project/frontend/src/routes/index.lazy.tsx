import { Link, createLazyFileRoute } from "@tanstack/react-router";
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow
} from "@/components/ui/table";
import { BsFillLightningFill } from "react-icons/bs";
import { PiShareNetworkBold } from "react-icons/pi";
import { FaCloud } from "react-icons/fa";
import { GiCube } from "react-icons/gi";
import { IconType } from "react-icons";
import { useEffect, useRef, useState } from "react";
import { BeatLoader } from "react-spinners";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
  } from "@/components/ui/select"
import { error } from "console";
import { DialogClose } from "@radix-ui/react-dialog";

export const Route = createLazyFileRoute("/")({
    component: Index
});

function Index() {
    const [servers, setServers] = useState<Server[] | null>(null);
    const [serverRows, setServerRows] = useState<JSX.Element[]>([]);

    const [location, setLocation] = useState<string>("hil");
    const [config, setConfig] = useState<string>("CCX13");
    const [name, setName] = useState<string>("")

    const makeRow = (title: string, location: string, Icon: IconType, disabled: boolean = false) => {
        return (
            <TableRow key={title}>
                <Link to={location} disabled={disabled}>
                    <div className="p-2">
                        <Icon className="mr-2 h-4 w-4 inline" />
                        {title}
                    </div>
                </Link>
            </TableRow>
        );
    };

    function deployNewServer() {
        fetch(`http://localhost:8000/api/createserver?name=${name}&location=${location}&instance_type=${config}`, {
            method: "GET"
        })
        .then((res) => res.json())
        .then((data) => {
            console.log(data);
        })
        .catch((err) => {
            console.error(err);
        });
    }

    function refreshServers(repeat: boolean = false) {
        fetch("http://localhost:8000/api/getservers")
            .then((res) => res.json())
            .then((data) => {
                setServers(data);
            });
    }

    useEffect(() => {
        refreshServers(true);
    }, []);

    useEffect(() => {
        if (servers && servers.map) {
            setServerRows(
                servers.map((server) => {
                    return makeRow(
                        server.ready ? `${server.name} (${server.public_net})` : `${server.name} (Initializing...)`,
                        `/server/${server.public_net}/1111`,
                        FaCloud
                    );
                })
            );
        }
    }, [servers]);

    return (
        <div className="p-2">
            <h1 className="mb-4">Servers</h1>
            <div className="pb-4">
                <h2>Hetzner</h2>
                <hr className="mb-2" />
                <Dialog>
                    <DialogTrigger asChild>
                        <Button size='sm' className="mb-2" variant="outline">Create New Server</Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>Create New Server</DialogTitle>
                            <DialogDescription>
                                Deploy a new Minecraft server to Hetzner Public Cloud.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="grid gap-4 py-4">
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="name" className="text-right">
                                    Name
                                </Label>
                                <Input id="name" className="col-span-3" value={name} onChange={(e) => {setName(e.currentTarget.value)}} />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label className="text-right">
                                    Location
                                </Label>
                                <Select onValueChange={(e) => {setLocation(e)}}>
                                    <SelectTrigger>
                                        <SelectValue>{location == "" ? "Choose a location..." : location}</SelectValue>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectItem value="nbg1">nbg1 (Nerumburg, DE)</SelectItem>
                                            <SelectItem value="hel1">hel1 (Helsinki, FI)</SelectItem>
                                            <SelectItem value="hil">hil (Ohio, USA)</SelectItem>
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label className="text-right">
                                    Configuration
                                </Label>
                                <Select onValueChange={(e) => {setConfig(e)}}>
                                    <SelectTrigger>
                                        <SelectValue>{config}</SelectValue>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectItem value="CPX21">CPX21 (4GB RAM, 2 Shared vCPU)</SelectItem>
                                            <SelectItem value="CCX13">CCX13 (8GB RAM, 2 vCPU)</SelectItem>
                                            <SelectItem value="CCX23">CCX23 (16GB RAM, 4 vCPU)</SelectItem>
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                        <DialogFooter>
                            <DialogClose>
                                <Button type="submit" onClick={() => {
                                    deployNewServer();
                                }}>Deploy</Button>
                            </DialogClose>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
                <div className="p-4 border">
                    <Table>
                        <TableBody className="w-full">
                            {serverRows.length > 0 ? (
                                serverRows
                            ) : (
                                <TableRow className="w-full flex">
                                    <BeatLoader className="m-auto" color="#36d7b7" />
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </div>
            </div>
            <div className="pb-4">
                <h2>Test Servers</h2>
                <hr className="mb-4" />
                <div className="p-4 border">
                    <Table>
                        <TableBody>
                            {makeRow("Localhost daemon", "/server/localhost/1111", GiCube)}
                            {makeRow("Direct to Cloud via WebSocket", "/direct-via-ws", BsFillLightningFill)}
                            {makeRow("Proxy via WebSocket", "/mitm-via-ws", PiShareNetworkBold)}
                        </TableBody>
                    </Table>
                </div>
            </div>
        </div>
    );
}
