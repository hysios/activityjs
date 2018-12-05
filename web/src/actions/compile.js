import { COMPILE_FAILED, COMPILE_GOCODE, EVALUATED_CODE, EVALUATED_CODE_ERRORS, VET_GOCODE } from '../actionTypes';
import { apiAction } from './api';

export function compile(code, debug=false) {
    let q = ""
    if (debug) {
        q = "?debug=true"
    }

    return apiAction({
        url: `/api/compile${q}`,
        method: 'POST',
        onSuccess: fetchActivityScript,
        onFailure: compileErrorCode(code),
        label: COMPILE_GOCODE,
        data: code
    })
}

export function vet(code) {
    return apiAction({
        url: `/api/vet`,
        method: 'POST',
        onSuccess: vetErrors,
        label: VET_GOCODE,
        data: code
    })
}


export function fetchActivityScript({js, size, duration, errors}) {
    loadScript(js, () => { })
    return {
        type: EVALUATED_CODE,
        size,
        duration: duration,
        errors
    }
}

export const compileErrorCode = (code) => (err) => {
    return {
        type: COMPILE_FAILED,
        status: err.status,
        statusText: err.msg,
        error: err.json,
        code: code
    }
}


// export function compileError(err) {
//     return {
//         type: COMPILE_FAILED,
//         status: err.status,
//         statusText: err.msg,
//         error: err.json
//     }
// }


export function vetErrors({lineErrors}) {
    return {
        type: EVALUATED_CODE_ERRORS,
        errors: lineErrors
    }
}

function loadScript(url, callback){

    var script = document.createElement("script")
    script.type = "text/javascript";

    if (script.readyState){  //IE
        script.onreadystatechange = function(){
            if (script.readyState == "loaded" ||
                    script.readyState == "complete"){
                script.onreadystatechange = null;
                callback();
            }
        };
    } else {  //Others
        script.onload = function(){
            callback();
        };
    }

    script.src = url;
    document.getElementsByTagName("head")[0].appendChild(script);
}