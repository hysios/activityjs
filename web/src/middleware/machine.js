import { ITEM_ADD_QUANTITY } from '../actionTypes'
import { changeCurrentItem } from '../actions'
import { applyCurrentEffects } from '../actions/effect'

const activityMachineMiddleware = ({getState, dispatch}) => next => action => {
    next(action)

    if (!window.ActivityMachine) return
    const ActivityMachine = window.ActivityMachine
    let store = getState()
    console.log(getState())

    switch (action.type) {
    case ITEM_ADD_QUANTITY: {
        // store.item
        let {item, products: {product: {id, price}}, user, order} = store
        let {effects, ...rest} = item
        let oriItem = {
            ...rest,
            id: id,
            price,
        }

        let ctx = { item: oriItem, user, order }

        console.log(ctx)
        let state = ActivityMachine.evaluate(ctx)
        let newItem = {
            ...oriItem,
            ...state.item,
        }
        dispatch(changeCurrentItem(newItem))
        console.log(state)

        effects = newItem.effects || {}
        // if (!effects) {
        //     resetEffect
        // }
        dispatch(applyCurrentEffects(effects))
        // Object.keys(effects).forEach((field) => {
        //     let effect = effects[field]
        //     console.log(field, effect)
        //     dispatch(applyEffect(field, effect))
        // })
        break;
    }
    default:
        return
    }
}

export default activityMachineMiddleware;
