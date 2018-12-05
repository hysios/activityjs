export class HttpError extends Error {
    constructor(status, msg, json) {
        super(msg)
        this.status = status
        this.msg = msg
        this.json = json
    }
}
