import { PRODUCT_TO_ITEM } from '../actionTypes';
import { genProductCode } from '../actions/code';


export default ({dispatch}) => next => action =>  {
    next(action)

    if (action.type != PRODUCT_TO_ITEM ) return

    dispatch(genProductCode(action.product))
}