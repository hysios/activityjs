import { APPLY_EFFECT, RESET_EFFECT, APPLY_CURRENT_EFFECTS } from '../actionTypes';
import { camelCase } from 'lodash';

const initialEffects = {
}

export default function effect(state = initialEffects, action) {
    switch (action.type) {
    case APPLY_EFFECT:{
        return {
            [action.field]: action.effect,
        }
    }
    case RESET_EFFECT: {
        return {
            [action.field]: [],
        }
    }
    case APPLY_CURRENT_EFFECTS: {
        let { effects } = action
        console.log(action)
        let newState = {}
        Object.keys(effects).forEach((key) => newState[camelCase(key)] = effects[key] )
        return {
            ...newState
        }
    }
    default:
        return state;
    }
}