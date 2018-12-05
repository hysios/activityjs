import React, { Component } from 'react'
import { connect } from 'react-redux';
import { Currency } from './utils'
import { sprintf } from 'sprintf-js';

const Subtotal = ({ subtotal }) => (
    <Currency>小计￥ {sprintf("%0.2f", subtotal)}</Currency>
)

const mapStateToProps = state => {
    let { item } = state

    return { subtotal: item.price * item.quantity }
}

export default connect(
    mapStateToProps
)(Subtotal)
