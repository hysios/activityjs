import { combineReducers } from 'redux'
import code from './code'
import compile from './compile'
import effect from './effect'
import api from './api'

import {
    ADD_ITEM,
    ITEM_ADD_QUANTITY,
    ITEM_CHANGE_QUANTITY,
    PRODUCT_TO_ITEM,
    REMOVE_ITEM,
    CHANGE_CURRENT_ITEM,
} from '../actionTypes';

const initialState = {
    products: [],
    product: {},
}

const initialOrder = {
    items: []
}

const initialItem = {
    productId: "",
    quantity: 1,
    defaultImage: "",
    imageUrls: [],
    price: 0,
    title: "",
    product: null
}

const initialCart = {
    items: []
}

function products(state = initialState, action) {
    switch (action.type) {
    case PRODUCT_TO_ITEM: {
        let { defaultImage, imageUrls, productId, ...rest } = action.product;
        return {
            product: {
                id: productId,
                defaultImage: defaultImage ? defaultImage : imageUrls[0],
                imageUrls,
                ...rest,
            }
        }
    }
    default:
        return state
    }
}

function order(state = initialOrder, action) {
    return initialOrder
}

function item(state = initialItem, action) {
    switch (action.type) {
    case PRODUCT_TO_ITEM: {
        let { defaultImage, imageUrls, ...rest } = action.product;
        return {
            defaultImage: defaultImage ? defaultImage : imageUrls[0],
            imageUrls,
            ...rest,
        }
    }
    case ITEM_ADD_QUANTITY: {
        let {quantity, ...rest} = state
        quantity += action.value
        if (quantity <= 1) {
            quantity = 1
        }
        return {
            quantity,
            ...rest
        }
    }
    case ITEM_CHANGE_QUANTITY: {
        let {quantity, ...rest} = state
        let n = parseInt(action.value)
        if (!isNaN(n)) {
            quantity = n
        }

        return {
            quantity,
            ...rest
        }
    }
    case CHANGE_CURRENT_ITEM: {
        return action.change
    }
    default:
        return state
    }
}

function cart(state = initialCart, action) {
    switch(action.type) {
    case ADD_ITEM: {
        if (action.cart) {
            let {items} = state
            let item = action.item

            let idx = items.findIndex((it) => it.productId == item.productId)
            if (idx >= 0 ) {
                let oldItem = items[idx]
                let quantity = oldItem.quantity + item.quantity
                return {
                    items: [...items.slice(0, idx), {...item, quantity}, ...items.slice(idx +1)]
                }
            } else {
                return {
                    items: [...items, item]
                }
            }
        }
        return state
    }
    case REMOVE_ITEM: {
        if (action.cart) {
            let {items} = state
            let idx = items.findIndex((it) => it.productId == item.productId)
            return {
                items: [...items.slice(0, idx), ...items.slice(idx +1)]
            }
        }
        return state
    }
    default:
        return state
    }
}

const app = combineReducers({
    products,
    order,
    item,
    cart,
    code,
    api,
    compile,
    effect,
})

export default app