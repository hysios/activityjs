import { COMPILE_FAILED, EVALUATED_CODE, EVALUATED_CODE_ERRORS } from '../actionTypes';

const initialCompile = {
    state: null
}

export default function compile(state = initialCompile, action) {
    switch (action.type) {
    case EVALUATED_CODE: {
        return {
            state: "success",
            size: action.size,
            errors: action.errors,
            duration: action.duration,
            annotations: []
        }
    }
    case COMPILE_FAILED: {
        return {
            state: "error",
        }
    }
    case EVALUATED_CODE_ERRORS: {
        return {
            state: "error",
            errors: action.errors,
            annotations: action.errors.map(({line, pos, msg}) => ({row: line-1, column: pos, type: 'error', text: msg}))
        }
    }
    default:
        return state;
    }
}