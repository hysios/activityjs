import {
  ADD_ITEM_TO_CODE,
  ADD_PRODUCT_TO_CODE,
  GENERATE_CODE,
  MODIFY_CODE,
  SPLIT_CODE_WIDTH,
  UPDATE_CODE
} from '../actionTypes';
import { generateCode } from '../utils/compile'

const initialCode = {
    code: "",
    change: false,
    productCode: "",
    userCode: "",
    itemCode: "",
    codeEditorWidth: "600px",
}

export default function code(state = initialCode, action) {
    switch (action.type) {
        case GENERATE_CODE: {
            return {
                [`${action.model}Code`]: generateCode(action.model, action.value),
                ...state
            }
        }
        case ADD_PRODUCT_TO_CODE: {
            return {
                product: action.value
            }
        }
        case ADD_ITEM_TO_CODE: {
            return {
                item: action.value
            }
        }
        case UPDATE_CODE: {
            return {
                ...state,
                code: action.code
            }
        }
        case MODIFY_CODE: {
            return {
                change: true
            }
        }
        case SPLIT_CODE_WIDTH: {
            return {
                codeEditorWidth: cssWidth(action.value),
            }
        }
        default:
            return state;
    }
}

function cssWidth(val, def) {
    if (typeof val === 'string') {
        return val
    } else if (typeof val === 'number') {
        return `${val}px`;
    } else {
        return def;
    }
}
