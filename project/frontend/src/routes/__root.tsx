import Navbar from "@/components/navbar";
import { createRootRoute, Link, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
    component: () => (
        <div className="w-screen max-w-screen">
            <Navbar/>
            <Outlet />
        </div>
    )
});
