import { applyMiddleware, createStore, compose } from 'redux'
import app from './reducers'
import logger from 'redux-logger'
import apiMiddleware from './middleware/api';
import codeMiddleware from './middleware/code';
import compileMiddleware from './middleware/compile'
import activityMachineMiddleware from './middleware/machine';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const store = createStore(app, /* preloadedState, */ composeEnhancers(
    applyMiddleware(logger, apiMiddleware, codeMiddleware, compileMiddleware, activityMachineMiddleware)
))

export default store