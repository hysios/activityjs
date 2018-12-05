// inspired by https://leanpub.com/redux-book
import { API } from "../actionTypes";
import { accessDenied, apiError, apiStart, apiEnd } from "../actions/api";
import path from 'path'
import { HttpError } from './errors'

const apiMiddleware = ({ dispatch }) => next => action => {
    next(action);

    if (action.type !== API) return;

    const {
        url,
        method,
        data,
        query,
        accessToken,
        onSuccess,
        onFailure,
        label,
        headers
    } = action.payload;
    const dataOrParams = ["GET", "DELETE"].includes(method) ? "qs" : "body";

    const defaults = {
        baseURL: process.env.REACT_APP_API_URL || process.env.REACT_APP_BASE_URL || "",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`
        }
    };

    const apiurl = defaults.baseURL + url;

    if (label) {
        dispatch(apiStart(label));
    }

    fetch(apiurl, {
        method,
        headers,
        qs: query,
        body: data
    })
        .then(res => {
            return catchHttpError(res)
        })
        .then(data => {
            dispatch(onSuccess(data, label));
        })
        .catch(error => {
            dispatch(apiError(error));
            console.log(error.json)
            let failedState = onFailure(error);
            if (failedState != undefined) {
                dispatch(failedState);
            }

            if (error.response && error.response.status === 403) {
                dispatch(accessDenied(window.location.pathname));
            }
        })
        .finally(() => {
            if (label) {
                dispatch(apiEnd(label));
            }
        });
};

function catchHttpError(res) {
    if (res.ok) {
        return res.json()
    }

    return res.json().then((jsonError) => {
        throw new HttpError(res.status, res.statusText, jsonError)
    })
}

export default apiMiddleware;
