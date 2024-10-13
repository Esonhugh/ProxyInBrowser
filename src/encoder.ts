import { debug } from './debug'

export type TrafficEncoder = (c :string) => string

export let EncodeTraffic :TrafficEncoder = (c :string) => {
    return c;
}

export let DecodeTraffic :TrafficEncoder = (c :string) => {
    debug("decode to", c)
    return c
}

