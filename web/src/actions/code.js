import {
  ADD_ITEM_TO_CODE,
  ADD_PRODUCT_TO_CODE,
  FETCH_EXAMPLE,
  GENERATE_CODE,
  MODIFY_CODE,
  SPLIT_CODE_WIDTH,
  UPDATE_CODE,
} from '../actionTypes';
import { apiAction } from './api';

export function genProductCode(product) {
    return {
        type: GENERATE_CODE,
        model: "Product",
        value: product
    }
}

export function genUserCode(user) {
    return {
        type: GENERATE_CODE,
        model: "User",
        value: user
    }
}

export function genItemCode(item) {
    return {
        type: GENERATE_CODE,
        model: "OrderItem",
        value: item
    }
}

export function setCode(code) {
    return {
        type: UPDATE_CODE,
        code
    }
}

export function modifyCode() {
    return {
        type: MODIFY_CODE
    }
}

export function fetchExample(filename) {
    return apiAction({
        url: `/api/examples/${filename}`,
        onSuccess: loadExample,
        label: FETCH_EXAMPLE
    })
}

export function loadExample(data) {
    return setCode(data.file)
}

export function addProductToCode(product) {
    return {
        type: ADD_PRODUCT_TO_CODE,
        value: product
    }
}

export function addItemToCode(item) {
    return {
        type: ADD_ITEM_TO_CODE,
        value: item
    }
}

export function splitCode(width) {
    return {
        type: SPLIT_CODE_WIDTH,
        value: width,
    }
}
