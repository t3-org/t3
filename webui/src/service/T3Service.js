export const Source = { // ticket source values.
    GRAFANA: "grafana"
}

export const SeveritiesList = [
    {name: "Low", value: "low"},
    {name: "Medium", value: "medium"},
    {name: "High", value: "high"},
]

export default class T3Service {
    static baseAPI = "api/v1"
    static perPage = 15

    url(p) {
        const trimmed = p.startsWith("/") ? p.slice(1) : p
        return `${T3Service.baseAPI}/${trimmed}`
    }

    async fetchTickets(query, page, perPage) {
        if (!perPage) {
            perPage = T3Service.perPage
        }

        const params = new URLSearchParams({
            query,
            page,
            'per_page': perPage ? perPage : T3Service.perPage,
        })
        const resp = await fetch(this.url(`tickets?`) + params)
        return {status: resp.status, body: await resp.json()}
    }

    async patchTicket(id, ticket) {
        const rawResponse = await fetch(this.url(`tickets/${id}`), {
            method: 'PATCH',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(ticket)
        });
        const content = await rawResponse.json()
        console.log(content);
        return {status: rawResponse.status, body: rawResponse.json()}
    }

    getProductsSmall() {
        return fetch('demo/data/products-small.json')
            .then((res) => res.json())
            .then((d) => d.data);
    }

    getProducts() {
        return fetch('demo/data/products.json')
            .then((res) => res.json())
            .then((d) => d.data);
    }

    getProductsWithOrdersSmall() {
        return fetch('demo/data/products-orders-small.json')
            .then((res) => res.json())
            .then((d) => d.data);
    }
}
