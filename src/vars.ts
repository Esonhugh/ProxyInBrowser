import { debug } from "./debug";

export const project: string = "__ENDPOINT__";
export const reverse_host: string = "localhost:9999";

export const trace_label: string = "x_client_trace_id";

export function uuidv4(): string {
  return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, (c) =>
    (
      +c ^
      (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (+c / 4)))
    ).toString(16)
  );
}

export let client_trace_id = initTraceId();

function initTraceId(): string {
  let client_trace_id = uuidv4();
  try {
    if (localStorage.getItem(trace_label) !== null) {
      client_trace_id = localStorage.getItem(trace_label)!;
      debug("Client trace id loaded");
    } else {
      localStorage.setItem(trace_label, client_trace_id);
      debug("Client trace id saved");
    }
  } catch (e) {
    debug("Cannot save client trace id");
    // ignore
  } finally {
    debug(`Client trace id: ${client_trace_id}`);
    return client_trace_id
  }
}
