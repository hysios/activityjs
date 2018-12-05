import React, { Component } from 'react'
import { connect } from 'react-redux'
import styled, { css } from 'react-emotion'
import { Button, Input, InputGroup, InputGroupAddon, } from 'reactstrap'
import { itemAddQuantity, itemChangeQuantity } from '../../actions'

const ButtonGroup = styled(InputGroup)`
    font-size: 1.5rem;
    input {
        width: 10rem!important;
        flex: 0!important;
    }
`
const Minus = (props) => (
    <InputGroupAddon addonType="prepend"><Button onClick={props.onClick}>-</Button></InputGroupAddon>
)
const Plus = (props) => (
    <InputGroupAddon addonType="append"><Button onClick={props.onClick}>+</Button></InputGroupAddon>
)

class Quantity extends Component {
    state = {
        quantity: 1
    }

    handleMinus = (e) => {
        e.preventDefault()
        this.props.itemAddQuantity(-1)
    }

    handlePlus = (e) => {
        e.preventDefault()
        this.props.itemAddQuantity(1)
    }

    handleChange = (e) => {
        this.props.itemChangeQuantity(e.target.value)
    }

    handleEnter = (e) => {
        e.target.select()
    }

    render() {
        let { quantity } = this.props;

        return <ButtonGroup>
            数量 <Minus onClick={this.handleMinus} /><Input value={quantity} onFocus={this.handleEnter} onChange={this.handleChange} /><Plus onClick={this.handlePlus} />
        </ButtonGroup>
    }
}

const mapStateToProps = state => {
    let { item } = state

    return { quantity: item.quantity }
}

export default connect(
    mapStateToProps, {
        itemAddQuantity,
        itemChangeQuantity
    }
)(Quantity)