import { API_END, API_START } from '../actionTypes';

const initialState = {
    working: {}
}

export default function api(state = initialState, action) {
    switch (action.type) {
    case API_START: {
        let {payload} = action
        let {working} = state

        return {
            working: {
                [payload.label]: true,
                ...working
            }
        }
    }
    case API_END: {
        let {payload} = action
        let {working} = state

        delete working[payload.label]
        return {
            working: {
                ...working
            }
        }
    }
    default:
        return state;
    }
}