import {
    NavigationMenu,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
    navigationMenuTriggerStyle
} from "@/components/ui/navigation-menu";

import { Link } from "@tanstack/react-router";
import { ModeToggle } from "./theme/mode-switch";
import { Button } from "./ui/button";

export default function Navbar() {
    return (
        <div className="flex-row w-screen">
            <NavigationMenu className="w-screen max-w-screen justify-between p-2">
                <NavigationMenuList>
                    <NavigationMenuItem>
                        <NavigationMenuLink className={navigationMenuTriggerStyle()} asChild>
                            <Button variant="outline">
                                <Link to="/">Servers</Link>
                            </Button>
                        </NavigationMenuLink>
                    </NavigationMenuItem>
                </NavigationMenuList>
                <NavigationMenuList className="w-full justify-end">
                    <NavigationMenuItem>
                        <ModeToggle/>
                    </NavigationMenuItem>
                </NavigationMenuList>
            </NavigationMenu>
        </div>
    );
}
