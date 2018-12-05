import React, {Component} from 'react'
import {sprintf} from 'sprintf-js'

const Currency = ({as = "p", money, className, prefix=null, flag="￥"}) => {
    return  React.createElement(as, {money, className}, React.Children.toArray([prefix, flag, sprintf("%0.2f", money)]))
}

export default Currency