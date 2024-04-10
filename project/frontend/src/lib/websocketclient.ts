export default class WebsocketClient {
    host: string;
    port: number;
    ws: WebSocket | null;

    onMessage: (message: WsMessage) => void;
    onDiagnosticMessage: (message: WsMessage) => void;

    constructor() {
        this.host = "";
        this.port = 0;
        this.ws = null;

        this.onMessage = (message: WsMessage) => {
            console.log(message);
        }

        this.onDiagnosticMessage = (message: WsMessage) => {
            console.log(message);
        }
    }

    connect(host: string, port: number) {
        this.host = host;
        this.port = port;

        this.ws = new WebSocket(`ws://${this.host}:${this.port}`);

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.command == "output") {
                this.onMessage(data);
            } else {
                this.onDiagnosticMessage(data);
            }
        }
    }

    setOnMessage(onMessage: (message: WsMessage) => void) {
        this.onMessage = onMessage;
    }

    setOnDiagnosticMessage(onMessage: (message: WsMessage) => void) {
        this.onDiagnosticMessage = onMessage;
    }

    send(message: string) {
        if (this.ws) {
            const msg = {
                command: "input",
                value: message
            }
            this.ws.send(JSON.stringify(msg));
        }
    }

    sendJson(message: Record<string, unknown>) {
        if (this.ws) {
            this.ws.send(JSON.stringify(message));
        }
    }
}