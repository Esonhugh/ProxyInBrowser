export const DEBUG = true

export function debug(...c :any) {
    if (DEBUG) {
        console.log(c)
    }
}