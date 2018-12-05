import { camelCase } from 'lodash'

import {
  APPLY_EFFECT,
  EFFECT_NEED_COMPILE,
  RESET_EFFECT,
  APPLY_CURRENT_EFFECTS,
} from '../actionTypes';

export function needCompileBefore(msg) {
    return {
        type: EFFECT_NEED_COMPILE,
        msg
    }
}

export function applyEffect(field, effect) {
    return {
        type: APPLY_EFFECT,
        field: camelCase(field),
        effect: effect,
    }
}

export function resetEffect(field) {
    return {
        type: RESET_EFFECT,
        field: camelCase(field),
    }
}

export function applyCurrentEffects(effects = {}) {
    return {
        type: APPLY_CURRENT_EFFECTS,
        effects,
    }
}