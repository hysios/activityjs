import { camelCase, max } from 'lodash'

let modelAttributes = {
    "Product": {
        "productId": { type: "ID", name: "ID" },
        "title": "Title",
        "price": "DefaultPrice"
    },
    "OrderItem": {
        "id": { type: "ID", name: "ID" },
        "price": "Price",
        "title": "Title",
        "quantity": "Quantity",
        "subtotal": "Subtotal",
        "produceDate": "ProduceDate",
        "productID": { type: "ID", name: "ProductId" },
    }
}

let modelNames = {
    "Product": "Product",
    "Item": "OrderItem"
}

export function generateCode(model, value) {
    let varName = camelCase(model)
    let modelName = getModelName(model)

    return `var ${varName} = model.${modelName} {
${generateAttributes(modelName, value)}
}

`
}

export function generateVarCode(model, value) {
    let varName = camelCase(model)
    let modelName = getModelName(model)

    return `var ${varName} = model.${modelName} {
${generateAttributes(modelName, value)}
}
`
}

export function generateInSliceCode(model, value) {
    let varName = camelCase(model)
    let modelName = getModelName(model)

    return addIndent(1, `model.${modelName} {
${generateAttributes(modelName, value)}
},`)
}

const times = x => f => {
    if (x > 0) {
        f()               // has to be side-effecting function
        times(x - 1)(f)
    }
}

export function addIndent(n, code) {
    let lines = code.split("\n")
    let indents = ""
    times(n)(() => indents += '\t')
    return lines.map((l) => indents + l).join("\n")
}

function getModelName(model) {
    let name = modelNames[model]
    if (name) {
        return name
    } else {
        return model
    }
}

function generateAttributes(model, value) {
    let config = modelAttributes[model] || {}
    let pad = computeAttrsMaxLength(model, value)

    let attrs = Object.keys(config).map((key) => {
        let conv = getAttr(config[key])
        if (value[key]) {
            return "\t" + padRight(pad + 2, conv.name + ": ") + attrVal(conv, value[key]) + ","
        } else {
            return null
        }
    });
    attrs = attrs.filter((m) => m)
    return attrs.join("\n")
}

function attrVal(attr, val) {
    if (attr.type === 'ID') {
        if (typeof val === 'string') {
            return `model.StringID("${val}")`
        } else if (typeof val === 'number') {
            return `model.IntID(${val})`
        } else {
            return JSON.stringify(val)
        }
    } else {
        return JSON.stringify(val)
    }
}

function getAttr(conv) {
    if (typeof conv === 'object') {
        return conv
    } else if (typeof conv === 'string') {
        return { type: 'value', name: conv }
    }
}

function computeAttrsMaxLength(model, value) {
    let config = modelAttributes[model] || {}
    let lens = Object.keys(config).map((m) => config[m].length)
    let maxlength = max(lens)
    return maxlength
}

function padRight(n, str) {
    let rest = n - str.length
    let right = ''
    for (let i = 0; i < rest; i++) {
        right += ' '
    }
    return str + right
}

function padleft(n, str) {
    let rest = n - str.length
    let left = ''
    for (let i = 0; i < rest; i++) {
        left += ' '
    }
    return left + str
}