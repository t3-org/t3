export const Source = { // ticket source values.
    GRAFANA: "grafana"
}

export const SeveritiesList = [
    {name: "Low", value: "low"},
    {name: "Medium", value: "medium"},
    {name: "High", value: "high"},
]

export const ticketZeroVal = {title: "", is_spam: false, is_firing: true}

export default class T3Service {
    static baseAPI = "/api/v1"
    static perPage = 15

    url(p) {
        const trimmed = p.startsWith("/") ? p.slice(1) : p
        return `${T3Service.baseAPI}/${trimmed}`
    }

    async toResp(rawResp) {
        const contentType = rawResp.headers.get("content-type");
        const isJson = contentType && contentType.indexOf("application/json") !== -1
        const body = isJson ? await rawResp.json() : await rawResp.text()

        if (rawResp.status !== 200) {
            console.log("failed request", {raw: rawResp, body})
        }

        return {status: rawResp.status, body}
    }

    async getTicket(id) {
        return this.toResp(await fetch(this.url(`tickets/${id}`)))
    }

    async queryTickets(query, page, perPage) {
        if (!perPage) {
            perPage = T3Service.perPage
        }

        const params = new URLSearchParams({
            query,
            page,
            'per_page': perPage ? perPage : T3Service.perPage,
        })
        return this.toResp(await fetch(this.url(`tickets?`) + params))
    }

    async patchTicket(id, ticket) {
        const resp = await fetch(this.url(`tickets/${id}`), {
            method: 'PATCH',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(ticket)
        });
        return this.toResp(resp)
    }

    async createTicket(ticket) {
        const resp = await fetch(this.url(`tickets`), {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(ticket)
        });

        return this.toResp(resp)
    }
}
