import algoliasearch from 'algoliasearch'
var client = algoliasearch('T1U4R3NIE6', '059b01b66c0f5e044216a1d1b08528f0');
var index = client.initIndex('product');

// const _search = Promise.promisify(index.search)

export function search(text) {
    return index.search(text, {
        "hitsPerPage": 10,
        "page": 0,
        "analytics": false,
        "attributesToRetrieve": "*",
        "getRankingInfo": true,
        "responseFields": "*",
        "attributesToHighlight": ["*"],
        "facets": []
    });
}