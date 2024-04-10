import React, { useEffect, useRef, useState } from "react";
import { Xterm } from "xterm-react";
import { Terminal as XtermTerminal } from "xterm";
import WebsocketClient from "@/lib/websocketclient";

interface WebsocketTerminalProps {
    client: WebsocketClient;
}

export default function WebsocketTerminal(props: WebsocketTerminalProps) {
    const [terminal, setTerminal] = useState<XtermTerminal | null>(null);
    const [input, setInput] = useState("");
    const [buffer, setBuffer] = useState("");
    const divRef = useRef<HTMLDivElement>(null);
    const client = props.client;

    const onTermInit = (term: XtermTerminal) => {
        setTerminal(term);
        term.reset();
    };

    const onTermDispose = (term: XtermTerminal) => {
        setTerminal(null);
    };

    const onData = (data: string) => {
        const code = data.charCodeAt(0);

        if (terminal) {
            // If the user hits empty and there is something typed echo it.
            if (code === 13) {
                terminal.write("\r\n");
                client.send(input + "\r\n");
                setInput("");
            } else if (code == 127) {
                // Handle backspace
                if (input.length > 0) {
                    terminal.write("\b \b");
                    setInput(input.slice(0, -1));
                }
            } else if (code < 32 || code === 127) {
                // Disable control Keys such as arrow keys
                console.log(code);
                return;
            } else {
                // Add general key press characters to the terminal
                terminal.write(data);
                setInput(input + data);
            }
        }
    };

    useEffect(() => {
        client.setOnMessage((message) => {
            if (terminal && message.command == "output" && message.output) {
                terminal.write(message.output.replace(/\r/gim, "").replace(/\n/gim, "\r\n"));
            }
        });
    }, [terminal, client]);

    useEffect(() => {
        if (!divRef.current) return;
        const resizeObserver = new ResizeObserver(() => {
            if (terminal && divRef?.current) {
                const cols = terminal.cols;
                const xtermw = document.getElementsByClassName("xterm-screen")[0].clientWidth;
                const divw = divRef.current.clientWidth;
    
                terminal.resize(Math.floor(cols * (divw / xtermw)), terminal.rows);
            }
        });
        resizeObserver.observe(divRef.current);
        return () => resizeObserver.disconnect(); // clean up
    }, [terminal]);

    return (
        <div ref={divRef}>
            <Xterm onInit={onTermInit} onDispose={onTermDispose} onData={onData} />
        </div>
    );
}
