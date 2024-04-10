import { Link } from "@tanstack/react-router";
import { Button } from "./ui/button";
import { BsFillLightningFill } from "react-icons/bs";

export default function Sidebar() {
    const pages = [
        { title: "Home", path: "/" },
        { title: "About", path: "/about" }
    ];
    
    return (
        <div className="flex-col">
            <Link to="/" className="flex">
                <Button variant={"outline"} size={"nav"}>Home</Button>
            </Link>
            <Link to="/direct-via-ws" className="flex">
                <Button variant={"outline"} size={"nav"}><BsFillLightningFill className="mr-2 h-4 w-4"/> Direct via WS</Button>
            </Link>

        </div>
    )
}