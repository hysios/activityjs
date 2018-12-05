import { COMPILE_FAILED } from '../actionTypes';
import { vet } from '../actions/compile';


export default ({dispatch}) => next => action =>  {
    next(action)

    if (action.type != COMPILE_FAILED) return

    dispatch(vet(action.code))
}