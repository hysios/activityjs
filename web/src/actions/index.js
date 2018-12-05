import {
  ADD_ITEM,
  CHANGE_ITEM,
  FETCH_PRODUCT,
  FETCH_PRODUCTS,
  ITEM_ADD_QUANTITY,
  ITEM_CHANGE_QUANTITY,
  CHANGE_CURRENT_ITEM,
  PRODUCT_TO_ITEM,
  REMOVE_ITEM,
  SET_PRODUCTS,
} from '../actionTypes';
import {apiAction} from './api'

export function productToItem(product, quantity = 1) {
    return {
        type: PRODUCT_TO_ITEM,
        product: {
            productId: product.id,
            title: product.name,
            price: product.price,
            imageUrls: product.imageUrls,
            quantity: quantity
        }
    }
}

export function itemAddQuantity(value = 1) {
    return {
        type: ITEM_ADD_QUANTITY,
        value
    }
}


export function itemChangeQuantity(value) {
    return {
        type: ITEM_CHANGE_QUANTITY,
        value
    }
}

export function addItemToCart(item, quantity) {
    return {
        type: ADD_ITEM,
        cart: true,
        item,
        quantity
    }
}

export function changeCartItemQuantity(id, item) {
    return {
        type: CHANGE_ITEM,
        cart: true,
        itemId: id,
        change: item
    }
}

export function changeCurrentItem(item) {
    return {
        type: CHANGE_CURRENT_ITEM,
        change: item,
    }
}

export function removeFromCart(item) {
    return {
        type: REMOVE_ITEM,
        cart: true,
        item,
    }
}

export function fetchProducts() {
    return apiAction({
        url: "/api/products",
        onSuccess: setProducts,
        label: FETCH_PRODUCTS
    });
}

export function fetchProduct(id) {
    return apiAction({
        url: `/api/products/${id}`,
        onSuccess: setProductToItem,
        label: FETCH_PRODUCT
    });
}

export function setProductToItem(product, label) {
    return productToItem(product, 1)
}


export function setProducts(products) {
    return {
        type: SET_PRODUCTS,
        products
    }
}


